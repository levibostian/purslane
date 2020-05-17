package cloud

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/digitalocean/godo"

	"github.com/levibostian/purslane/config"
	"github.com/levibostian/purslane/ui"
)

type digitalocean struct {
	client godo.Client
	Region string
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

func (cloud digitalocean) createServer(coreConfig *config.CoreConfig, createServerConfig *config.CreateServerConfig, createdVolume *CreatedVolume) *CreatedServer {
	cloud.assertAuth()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	sshKeyID := cloud.getSSHKeyReferenceID(coreConfig)

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

	ui.Message("Server created. Waiting until server is ready for connection...")
	cloud.waitUntilServerAvailable(dropletID)

	var ipAddress *string
	for _, network := range droplet.Networks.V4 {
		if network.Type == "public" {
			ipAddress = &network.IPAddress
		}
	}
	if ipAddress == nil {
		ui.ShouldNotHappen(fmt.Errorf("Droplet created, but it does not have a public IP address. Droplet: %+v", droplet))
	}

	return &CreatedServer{strconv.FormatInt(int64(dropletID), 10), createRequest.Name, *ipAddress, 22, "root"}
}

func (cloud digitalocean) getSSHKeyReferenceID(coreConfig *config.CoreConfig) (sshKeyID int) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	ui.Message("Checking if the public SSH key given in config is found in your DigitalOcean account...")
	key, _, err := cloud.client.Keys.GetByFingerprint(ctx, coreConfig.PublicSSHKeyFingerprint)
	ui.HandleError(err)
	if key == nil {
		createdSSHKeyName := "Purslane Public SSH Key"

		ui.Message("Public SSH key given in config not found in DigitalOcean account. Adding it now. (Note: Name of created SSH key will be named %s in your account)", createdSSHKeyName)

		createRequest := &godo.KeyCreateRequest{
			Name:      createdSSHKeyName,
			PublicKey: coreConfig.PublicSSHKey,
		}

		transfer, _, err := cloud.client.Keys.Create(ctx, createRequest)
		ui.HandleError(err)

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

func (cloud digitalocean) waitUntilServerAvailable(dropletID int) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	droplet, _, err := cloud.client.Droplets.Get(ctx, dropletID)
	ui.HandleError(err)

	dropletStatus := droplet.Status
	if dropletStatus == "new" {
		var secondsToWait time.Duration = 2

		ui.Message(fmt.Sprintf("Server not ready. Retrying in %d seconds", secondsToWait))
		time.Sleep(secondsToWait * time.Second)
		cloud.waitUntilServerAvailable(dropletID)
	} else if dropletStatus == "active" {
		ui.Message("Server ready to connect!")
		return
	} else {
		ui.Abort(fmt.Sprintf("Created server is in an invalid state, %s. Purslane cannot use a server that is in this unknown state.", dropletStatus))
	}
}

// Wants format: "s-1vcpu-1gb"
// (s|m|c)-(\dvcpu)-(\dgb)
func getDropletSize(createServerConfig *config.CreateServerConfig) string {
	return fmt.Sprintf("%s-%dvcpu-%dgb", createServerConfig.SizeType, createServerConfig.SizeCPU, createServerConfig.SizeMemory)
}
