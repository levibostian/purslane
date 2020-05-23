package ssh

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/levibostian/Purslane/cloud"
	"github.com/levibostian/Purslane/config"
	"github.com/levibostian/Purslane/ui"
	"golang.org/x/crypto/ssh"
)

type PurslaneSSHExecutor struct {
	coreConfig *config.CoreConfig
	server     *cloud.CreatedServer
	volume     *cloud.CreatedVolume
	sshClient  *ssh.Client
}

type SSHExecutor interface {
	DockerImagePull() bool
	DockerRegistryLogin() bool
	RunDockerImage() bool
	Close()
}

func GetSSHExecutor(coreConfig *config.CoreConfig, server *cloud.CreatedServer, volume *cloud.CreatedVolume) SSHExecutor {
	signer, err := ssh.ParsePrivateKey([]byte(coreConfig.PrivateSSHKey))
	ui.HandleError(err)

	config := &ssh.ClientConfig{
		User: server.OSUsername,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", server.IPAddress, server.SSHPort), config)
	ui.HandleError(err)

	var sshExecutor SSHExecutor = PurslaneSSHExecutor{coreConfig, server, volume, client}

	return sshExecutor
}

func (executor PurslaneSSHExecutor) DockerImagePull() bool {
	return runSSHCommand(fmt.Sprintf("docker pull %s", executor.coreConfig.DockerImageName), executor.sshClient)
}

func (executor PurslaneSSHExecutor) DockerRegistryLogin() bool {
	return runSSHCommand(fmt.Sprintf("echo \"%s\" | docker login %s --username %s --password-stdin", executor.coreConfig.DockerRegistry.Password, executor.coreConfig.DockerRegistry.RegistryName, executor.coreConfig.DockerRegistry.Username), executor.sshClient)
}

func (executor PurslaneSSHExecutor) RunDockerImage() bool {
	// Everything before options being added
	dockerRunCommand := "docker run --name purslane_executed"

	if executor.volume != nil {
		dockerRunCommand = fmt.Sprintf("%s -v %s:%s", dockerRunCommand, executor.volume.FileSystemMount, executor.coreConfig.DockerRunConfig.VolumeMountPoint)
	}

	if extraArgs := executor.coreConfig.DockerRunConfig.ExtraArgs; extraArgs != nil {
		dockerRunCommand = fmt.Sprintf("%s %s", dockerRunCommand, *extraArgs)
	}

	// image name goes last in command.
	dockerRunCommand = fmt.Sprintf("%s %s", dockerRunCommand, executor.coreConfig.DockerImageName)

	return runSSHCommand(dockerRunCommand, executor.sshClient)
}

func (executor PurslaneSSHExecutor) Close() {
	executor.sshClient.Close()
}

func runSSHCommand(command string, client *ssh.Client) bool {
	session, err := client.NewSession()
	if err != nil {
		client.Close()
		ui.HandleError(err)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	ui.Debug("[COMMAND] %s", command) // Running in Debug only mode because commands may contain private information.
	err = session.Run(command)

	successful := err == nil
	if !successful {
		fmt.Println(err)
	}

	return successful
}
