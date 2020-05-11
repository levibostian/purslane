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

	return &CreatedVolume{volume.ID, createRequest.Name, createRequest.FilesystemLabel}
}

func (cloud digitalocean) createServer(coreConfig *config.CoreConfig, createServerConfig *config.CreateServerConfig, createdVolume *CreatedVolume) *CreatedServer {
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
			godo.DropletCreateSSHKey{ID: sshKeyID},
		},
	}

	droplet, _, err := cloud.client.Droplets.Create(ctx, createRequest)
	ui.HandleError(err)
	dropletID := droplet.ID

	cloud.waitUntilServerAvailable(dropletID)

	return &CreatedServer{strconv.FormatInt(int64(dropletID), 10), createRequest.Name}
}

func (cloud digitalocean) getSSHKeyReferenceID(coreConfig *config.CoreConfig) (sshKeyID int) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	key, _, err := cloud.client.Keys.GetByFingerprint(ctx, coreConfig.PublicSSHKeyFingerprint)
	ui.HandleError(err)
	if key == nil {
		createRequest := &godo.KeyCreateRequest{
			Name:      "Purslane Public SSH Key",
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

func (cloud digitalocean) waitUntilServerAvailable(dropletID int) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	droplet, _, err := cloud.client.Droplets.Get(ctx, dropletID)
	ui.HandleError(err)

	dropletStatus := droplet.Status
	if dropletStatus == "new" {
		time.Sleep(2 * time.Second)
		cloud.waitUntilServerAvailable(dropletID)
	} else if dropletStatus == "active" {
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
