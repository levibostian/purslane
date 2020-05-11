package cloud

import (
	"github.com/levibostian/purslane/config"
)

// Cloud provider
type Cloud interface {
	createVolume(*config.CreateVolumeConfig) *CreatedVolume
	createServer(*config.CoreConfig, *config.CreateServerConfig, *CreatedVolume) *CreatedServer
}

// CreatedVolume - info about created volume
type CreatedVolume struct {
	ID              string
	Name            string
	FileSystemMount string
}

// CreatedServer - info about created server
type CreatedServer struct {
	ID   string
	Name string
}
