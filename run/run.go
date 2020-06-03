package run

import (
	"github.com/levibostian/Purslane/cloud"
	"github.com/levibostian/Purslane/config"
	"github.com/levibostian/Purslane/ssh"
	"github.com/levibostian/Purslane/ui"
)

// Execute The run command. Runs the Docker container in a cloud vm. Delegates out work to other parts of the code.
func Execute() {
	/**
	Perform all checks first. To make sure the program is setup correctly as much as possible.

	Here, we are getting configurations to make sure we have the information that we need to run.
	*/
	ui.Debug("Getting core config")
	coreConfig := config.GetCoreConfig()
	if coreConfig == nil {
		ui.Abort("You did not configure the core configurations for Purslane. Purslane cannot run without it!")
	}
	ui.Debug("Core config: %+v", *coreConfig)

	ui.Debug("Getting cloud provider config")
	cloudConfig := config.GetCloudConfig()
	if cloudConfig == nil {
		ui.Abort("You did not configure a cloud provider. Purslane cannot run without it!")
	}
	ui.Debug("Cloud provider config: %+v", *cloudConfig)

	cloudProvider := cloud.GetCloudProvider(coreConfig, cloudConfig)

	ui.Debug("Getting volume config")
	createVolumeConfig := config.CreateVolume()
	var createdVolume *cloud.CreatedVolume = nil
	if createVolumeConfig != nil {
		ui.Debug("Volume config exists: %+v", *createVolumeConfig)

		ui.Message("Creating volume")
		createdVolume = cloudProvider.CreateVolume(createVolumeConfig)

		defer func() {
			ui.Message("Deleting created volume.")
			cloudProvider.DeleteVolume(createdVolume)
		}()
	}

	createServerConfig := config.CreateServer()
	ui.Debug("Create server config: %+v", *createServerConfig)

	ui.Message("Creating server")
	serverReference := cloudProvider.CreateServer(createServerConfig, createdVolume)
	defer func() {
		ui.Message("Deleting created server")
		cloudProvider.DeleteServer(serverReference)
	}()
	createdServer := cloudProvider.WaitForServerToBeReady(serverReference)
	ui.Debug("Server ready to connect %+v", *createdServer)

	// We want to only create 1 SSH session and run all commands against it.
	sshExecutor := ssh.GetSSHExecutor(coreConfig, createdServer, createdVolume)
	defer sshExecutor.Close()

	if coreConfig.DockerRegistry != nil {
		ui.Message("Logging Docker into Docker registry.")

		handleSSHSessionResult(sshExecutor.DockerRegistryLogin())
	}

	ui.Message("Pulling Docker image in new server")
	handleSSHSessionResult(sshExecutor.DockerImagePull())

	ui.Message("Running Docker container in new server")
	handleSSHSessionResult(sshExecutor.RunDockerImage())

	ui.Message("Success! Cleaning up now...")
	return // all defer statements will run in a stack, resulting in clean up functions being called.
}

func handleSSHSessionResult(successful bool) {
	if successful {
		return
	}

	ui.Abort("Command failed. Exiting...")
}
