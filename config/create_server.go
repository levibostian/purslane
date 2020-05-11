package config

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/levibostian/purslane/ui"
)

// CreateServerConfig -
type CreateServerConfig struct {
	Size       string // size string that is universally used for purslane. it will be converted to the cloud's own type later. This is parsed into smaller values below.
	SizeType   string
	SizeCPU    int
	SizeMemory int
}

// CreateServer -
func CreateServer() *CreateServerConfig {
	serverSize := GetEnv("SERVER_SIZE", "server.size")
	serverSizeString := "s-1cpu-1gb"
	if serverSize != nil {
		// regexr.com/53ust
		matched, err := regexp.MatchString(`(s|m|c)-(\dcpu)-(\dgb)`, *serverSize)
		ui.HandleError(err)
		if !matched {
			ui.Abort("Server size given not valid format. Not using default value because it may not be the size you need for the job. Provide a valid value.")
		}

		serverSizeString = *serverSize
	}

	serverSizeSplit := strings.Split(serverSizeString, "-")
	serverSizeType := serverSizeSplit[0]
	serverSizeCpus, _ := strconv.Atoi(strings.ReplaceAll(serverSizeSplit[1], "cpu", ""))
	serverSizeMemory, _ := strconv.Atoi(strings.ReplaceAll(serverSizeSplit[1], "gb", ""))

	return &CreateServerConfig{serverSizeString, serverSizeType, serverSizeCpus, serverSizeMemory}
}
