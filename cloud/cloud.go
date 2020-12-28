package cloud

import (
	"github.com/levibostian/purslane/config"
)

// Cloud provider
type CloudProvider interface {
	CreateVolume(*config.CreateVolumeConfig) *CreatedVolume
	CreateServer(*config.CreateServerConfig, *CreatedVolume) *CreatedServerReference
	WaitForServerToBeReady(*CreatedServerReference) *CreatedServer
	DeleteVolume(*CreatedVolume)
	DeleteServer(*CreatedServerReference)
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
