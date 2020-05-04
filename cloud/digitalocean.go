package cloud

import (
	"context"
	"fmt"
	"time"

	"github.com/digitalocean/godo"

	"github.com/levibostian/purslane/config"
	"github.com/levibostian/purslane/ui"
)

type digitalocean struct {
	client godo.Client
}

func (cloud digitalocean) createVolume(config config.CreateVolumeConfig) {
	createRequest := &godo.VolumeCreateRequest{
		Region:          "nyc1",
		Name:            "purslane-volume",
		Description:     "Storage for Droplet created by Purslane CLI",
		FilesystemLabel: "purslane",
		SizeGigaBytes:   int64(config.GIGS),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	volume, _, err := cloud.client.Storage.CreateVolume(ctx, createRequest)
	ui.HandleError(err)

	fmt.Printf("%+v\n", volume)
}

var digitaloceanCloud = digitalocean{*godo.NewFromToken("1233")}
