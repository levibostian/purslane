package cloud

import (
	"github.com/levibostian/purslane/config"
)

// Cloud provider
type Cloud interface {
	createVolume(*config.CreateVolumeConfig) *CreatedVolume
	createServer(*config.CoreConfig, *config.CreateServerConfig, *CreatedVolume) *CreatedServerReference
	waitForServerToBeReady(*CreatedServerReference) *CreatedServer
	deleteVolume(*CreatedVolume)
	deleteServer(*CreatedServerReference)
}

// CreatedVolume - info about created volume
type CreatedVolume struct {
	ID              string
	Name            string
	FileSystemMount string
}

// CreatedServer - info about created server. Available after server is ready.
type CreatedServer struct {
	Reference CreatedServerReference

	IPAddress  string
	SSHPort    int
	OSUsername string
}

type CreatedServerReference struct {
	DO *DigitalOceanCreatedServerExtras
}
