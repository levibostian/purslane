package cloud

import (
	"context"
	"fmt"
	"time"

	"github.com/digitalocean/godo"

	"github.com/levibostian/Purslane/config"
	"github.com/levibostian/Purslane/ui"
)

type digitalocean struct {
	client godo.Client
	Region string
}

type DigitalOceanCreatedServerExtras struct {
	ID              int
	CreatedSSHKeyId *int // We only want to delete the ssh key from your account if Purslane created it. That's why it's optional.
}

// GetDigitalOceanCloud -
func GetDigitalOceanCloud(config *config.CloudConfig) Cloud {
	return digitalocean{*godo.NewFromToken(config.APIKey), "nyc1"}
}

func (cloud digitalocean) createVolume(config *config.CreateVolumeConfig) *CreatedVolume {
	cloud.assertAuth()

	createRequest := &godo.VolumeCreateRequest{
		Region:          cloud.Region,
		Name:            "purslane-volume",
		Description:     "Storage for Droplet created by Purslane CLI",
		FilesystemLabel: "purslane",
		SizeGigaBytes:   int64(config.Gigs),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	volume, _, err := cloud.client.Storage.CreateVolume(ctx, createRequest)
	ui.HandleError(err)

	ui.Debug("DigitalOcean volume created: %+v", *volume)

	return &CreatedVolume{volume.ID, createRequest.Name, createRequest.FilesystemLabel}
}

func (cloud digitalocean) deleteVolume(createdVolume *CreatedVolume) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := cloud.client.Storage.DeleteVolume(ctx, createdVolume.ID)
	ui.HandleError(err)
}

func (cloud digitalocean) deleteServer(serverReference *CreatedServerReference) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if sshKeyID := serverReference.DO.CreatedSSHKeyId; sshKeyID != nil {
		_, err := cloud.client.Keys.DeleteByID(ctx, *sshKeyID)
		ui.HandleError(err)
	}

	_, err := cloud.client.Droplets.Delete(ctx, serverReference.DO.ID)
	ui.HandleError(err)
}

func (cloud digitalocean) createServer(coreConfig *config.CoreConfig, createServerConfig *config.CreateServerConfig, createdVolume *CreatedVolume) *CreatedServerReference {
	cloud.assertAuth()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	sshKeyID, createdSSHKey := cloud.getSSHKeyReferenceID(coreConfig)

	createRequest := &godo.DropletCreateRequest{
		Name:   "purslane",
		Region: cloud.Region,
		Size:   getDropletSize(createServerConfig),
		Image: godo.DropletCreateImage{
			Slug: "docker-18-04", // https://marketplace.digitalocean.com/apps/docker
		},
		SSHKeys: []godo.DropletCreateSSHKey{
			{ID: sshKeyID},
		},
	}

	droplet, _, err := cloud.client.Droplets.Create(ctx, createRequest)
	ui.HandleError(err)
	dropletID := droplet.ID
	ui.Debug("DigitalOcean server created: %+v", *droplet)

	var extrasSSHKeyID *int
	if createdSSHKey {
		extrasSSHKeyID = &sshKeyID
	}

	return &CreatedServerReference{&DigitalOceanCreatedServerExtras{dropletID, extrasSSHKeyID}}
}

func (cloud digitalocean) waitForServerToBeReady(serverReference *CreatedServerReference) *CreatedServer {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	ui.Message("Waiting until server is ready for connection...")

	ipAddress := loopUntilServerReady(ctx, cloud, serverReference)

	return &CreatedServer{*serverReference, ipAddress, 22, "root"}
}

func loopUntilServerReady(ctx context.Context, cloud digitalocean, serverReference *CreatedServerReference) (ipAddress string) {
	droplet, _, err := cloud.client.Droplets.Get(ctx, serverReference.DO.ID)
	ui.HandleError(err)

	for _, network := range droplet.Networks.V4 {
		if network.Type == "public" {
			ipAddress = network.IPAddress

			ui.Message("Server ready to connect!")

			return
		}
	}

	var secondsToWait time.Duration = 2

	ui.Message(fmt.Sprintf("Server not ready. Retrying in %d seconds", secondsToWait))
	time.Sleep(secondsToWait * time.Second)

	return loopUntilServerReady(ctx, cloud, serverReference)
}

func (cloud digitalocean) getSSHKeyReferenceID(coreConfig *config.CoreConfig) (sshKeyID int, createdSSHKey bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	ui.Message("Checking if the public SSH key given in config is found in your DigitalOcean account...")
	key, response, err := cloud.client.Keys.GetByFingerprint(ctx, coreConfig.PublicSSHKeyFingerprint)
	if key == nil && response.StatusCode != 404 { // only handle error if we are not already handing it.
		ui.HandleError(err)
	}
	sshKeyNeedsAddedToAccount := key == nil || response.StatusCode == 404

	createdSSHKey = false

	if sshKeyNeedsAddedToAccount {
		createdSSHKeyName := "Purslane Public SSH Key"

		ui.Message("Public SSH key given in config not found in DigitalOcean account. Adding it now. (Note: Name of created SSH key will be named %s in your account)", createdSSHKeyName)

		createRequest := &godo.KeyCreateRequest{
			Name:      createdSSHKeyName,
			PublicKey: coreConfig.PublicSSHKey,
		}

		transfer, _, err := cloud.client.Keys.Create(ctx, createRequest)
		ui.HandleError(err)

		createdSSHKey = true

		sshKeyID = transfer.ID
	} else {
		sshKeyID = key.ID
	}

	return
}

func (cloud digitalocean) assertAuth() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	apiKeyInvalidMessage := "Sorry! The API key you gave for DigitalOcean is not associated with a DigitalOcean account."

	account, response, err := cloud.client.Account.Get(ctx)
	if response.StatusCode == 401 { // If a response came back, that's good. This check here is to see if the api key does not belong to an account.
		ui.Abort(apiKeyInvalidMessage)
	}
	ui.HandleError(err) // Maybe there was not a response or the account is invalid. Handle it.

	if account == nil {
		ui.Abort(apiKeyInvalidMessage)
	} else {
		if account.Status != "active" {
			ui.Abort("Sorry! The API key you gave for DigitalOcean belongs to an account, but the account is not active. Login to your DigitalOcean account to fix it then try again.")
		}
	}

	return
}

// Wants format: "s-1vcpu-1gb"
// (s|m|c)-(\dvcpu)-(\dgb)
func getDropletSize(createServerConfig *config.CreateServerConfig) string {
	return fmt.Sprintf("%s-%dvcpu-%dgb", createServerConfig.SizeType, createServerConfig.SizeCPU, createServerConfig.SizeMemory)
}
