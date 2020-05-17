package config

import (
	"errors"

	"github.com/levibostian/purslane/ui"
	"github.com/levibostian/purslane/util"
)

type DockerRegistryConfig struct {
	RegistryName string
	Username     string
	Password     string
}

type DockerRunConfig struct {
	ExtraArgs        *string
	VolumeMountPoint string
}

type CoreConfig struct {
	PublicSSHKeyFingerprint string
	PublicSSHKey            string
	PrivateSSHKey           string
	DockerImageName         string
	DockerRegistry          *DockerRegistryConfig
	DockerRunConfig         DockerRunConfig
}

func GetCoreConfig() *CoreConfig {
	sshKeyFingerprint, publicSSHKey, privateSSHKey, err := getSSHKeyInfo()
	ui.HandleError(err)

	dockerImageName, dockerRegistryInfo, err := getDockerInformation()
	ui.HandleError(err)

	dockerRunInfo := getDockerRunInfo()

	return &CoreConfig{sshKeyFingerprint, publicSSHKey, privateSSHKey, dockerImageName, dockerRegistryInfo, dockerRunInfo}
}

func getDockerRunInfo() DockerRunConfig {
	mountPoint := "/home"
	if value := GetEnv("VOLUME_CONTAINER_MOUNT_PATH", "volume.container_mount_path"); value != nil {
		mountPoint = *value
	}

	var extraArgs *string
	if value := GetEnv("DOCKER_RUN_EXTRA_ARGS", "docker_run.extra_args"); value != nil {
		extraArgs = value
	}

	return DockerRunConfig{extraArgs, mountPoint}
}

func getDockerInformation() (imageName string, registryInfo *DockerRegistryConfig, err error) {
	dockerImageName := GetEnv("DOCKER_IMAGE", "docker.image")
	if dockerImageName == nil {
		return "", nil, errors.New("No docker image specified in configuration file")
	}

	registryName := GetEnv("DOCKER_REGISTRY_NAME", "docker.registry.name")
	registryUsername := GetEnv("DOCKER_REGISTRY_USERNAME", "docker.registry.username")
	registryPassword := GetEnv("DOCKER_REGISTRY_PASSWORD", "docker.registry.password")

	if registryName != nil && registryUsername != nil && registryPassword != nil {
		return *dockerImageName, &DockerRegistryConfig{*registryName, *registryUsername, *registryPassword}, nil
	} else {
		return *dockerImageName, nil, nil
	}
}

func getSSHKeyInfo() (fingerprint string, publicKey string, privateKey string, err error) {
	publicSSHKeyEnv := GetEnv("PUBLIC_SSH_KEY_PATH", "public_ssh_key_path")
	if publicSSHKeyEnv == nil {
		return "", "", "", errors.New("No public SSH key information set in the configuration file")
	}
	privateSSHKeyEnv := GetEnv("PRIVATE_SSH_KEY_PATH", "private_ssh_key_path")
	if privateSSHKeyEnv == nil {
		return "", "", "", errors.New("No private SSH key information set in the configuration file")
	}

	// ssh key
	publicSSHKeyPath := util.GetFullFilePath(*publicSSHKeyEnv)
	publicSSHKey := string(util.GetFileContents(publicSSHKeyPath, "Public SSH key"))
	publicSSHKeyFingerprint := util.GetSSHFootprint(publicSSHKey)

	privateSSHKeyPath := util.GetFullFilePath(*privateSSHKeyEnv)
	privateSSHKey := string(util.GetFileContents(privateSSHKeyPath, "Private SSH key"))

	return publicSSHKeyFingerprint, publicSSHKey, privateSSHKey, nil
}
