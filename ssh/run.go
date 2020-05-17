package ssh

import (
	"fmt"
	"os"

	"github.com/levibostian/purslane/cloud"
	"github.com/levibostian/purslane/config"
	"github.com/levibostian/purslane/ui"
	"golang.org/x/crypto/ssh"
)

type PurslaneSSHExecutor struct {
	coreConfig *config.CoreConfig
	server     *cloud.CreatedServer
	volume     *cloud.CreatedVolume
	session    *ssh.Session
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
	}

	conn, _ := ssh.Dial("tcp", fmt.Sprintf("%s:%d", server.IPAddress, server.SSHPort), config)
	session, err := conn.NewSession()
	ui.HandleError(err)

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	var sshExecutor SSHExecutor = PurslaneSSHExecutor{coreConfig, server, volume, session}

	return sshExecutor
}

func (executor PurslaneSSHExecutor) DockerImagePull() bool {
	return runSSHCommand(fmt.Sprintf("docker pull %s", executor.coreConfig.DockerImageName), executor.session)
}

func (executor PurslaneSSHExecutor) DockerRegistryLogin() bool {
	return runSSHCommand(fmt.Sprintf("echo \"%s\" | docker login %s --username %s --password-stdin", executor.coreConfig.DockerRegistry.Password, executor.coreConfig.DockerRegistry.RegistryName, executor.coreConfig.DockerRegistry.Username), executor.session)
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

	return runSSHCommand(dockerRunCommand, executor.session)
}

func (executor PurslaneSSHExecutor) Close() {
	executor.session.Close()
}

func runSSHCommand(command string, session *ssh.Session) bool {
	ui.Debug("[COMMAND] %s", command) // Running in Debug only mode because commands may contain private information.
	err := session.Run(command)

	return err != nil
}
