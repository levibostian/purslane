package run

import (
	"github.com/levibostian/purslane/cloud"
	"github.com/levibostian/purslane/config"
	"github.com/levibostian/purslane/ui"
)

// Execute The run command. Runs the Docker container in a cloud vm. Delegates out work to other parts of the code.
func Execute() {
	/**
	Perform all checks first. To make sure the program is setup correctly as much as possible.

	Here, we are getting configurations to make sure we have the information that we need to run.
	*/
	coreConfig := config.GetCoreConfig()
	if coreConfig == nil {
		ui.Abort("You did not configure the core configurations for Purslane. Purslane cannot run without it!")
	}

	cloudConfig := config.GetCloudConfig()
	if cloudConfig == nil {
		ui.Abort("You did not configure a cloud provider. Purslane cannot run without it!")
	}

	createVolumeConfig := config.CreateVolume()
	var createdVolume *cloud.CreatedVolume = nil
	if createVolumeConfig != nil {
		ui.Message("Creating volume")
		createdVolume = cloud.CreateVolume(cloudConfig, createVolumeConfig)
	}

	createServerConfig := config.CreateServer()
	ui.Message("Creating server")
	_ = cloud.CreateServer(coreConfig, cloudConfig, createServerConfig, createdVolume)
}
