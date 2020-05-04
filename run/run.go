package run

import (
	"github.com/levibostian/purslane/cloud"
	"github.com/levibostian/purslane/config"
	"github.com/levibostian/purslane/ui"
)

// Execute The run command. Runs the Docker container in a cloud vm. Delegates out work to other parts of the code.
func Execute() {
	cloudConfig := config.Cloud()
	if cloudConfig == nil {
		ui.Abort("You did not configure a cloud provider. Purslane cannot run without it!")
	}

	createVolumeConfig := config.CreateVolume()
	if createVolumeConfig != nil {
		ui.Message("Creating volume")
		cloud.CreateVolume(*createVolumeConfig)
	}
}
