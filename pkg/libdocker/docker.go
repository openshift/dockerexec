package libdocker

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/docker/libcontainer"
	"github.com/docker/libcontainer/configs"
	"github.com/docker/libcontainer/utils"
)

// RunInContainer runs a process in a running docker container using the options
// specified. It returns the exit code and/or error.
func RunInContainer(id string, options *DockerExecOptions) (int, error) {
	var factory libcontainer.Factory
	var config *configs.Config
	var containerConfig *ContainerConfig
	var err error

	containerId, err := expandContainerId(id)
	if err != nil {
		return -1, err
	}
	factory, err = loadDockerFactory()
	if err != nil {
		return -1, err
	}
	config, err = loadDockerConfig(containerId)
	if err != nil {
		return -1, err
	}

	container, err := factory.Load(containerId)
	if err != nil {
		return -1, err
	}

	containerConfig, err = loadContainerConfig(containerId)
	if err != nil {
		return -1, err
	}

	user := options.User
	if user == "" {
		user = containerConfig.Config.User
	}

	env := containerConfig.Config.Env
	env = append(env, options.Env...)

	hostname, err := getContainerHostName(containerId)
	if err == nil {
		hostnameEnv := fmt.Sprintf("HOSTNAME=%s", hostname)
		env = append(env, hostnameEnv)
	}

	process := &libcontainer.Process{
		Args:   options.Args,
		Env:    env,
		User:   user,
		Cwd:    options.Cwd,
		Stdin:  options.Stdin,
		Stdout: options.Stdout,
		Stderr: options.Stderr,
	}
	rootuid, err := config.HostUID()
	if err != nil {
		return -1, err
	}
	tty, err := newTty(options.Tty, process, rootuid)
	if err != nil {
		return -1, err
	}
	if err := tty.attach(process); err != nil {
		return -1, err
	}
	go handleSignals(process, tty)
	err = container.Start(process)
	if err != nil {
		tty.Close()
		return -1, err
	}

	status, err := process.Wait()
	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok {
			status = exitError.ProcessState
		} else {
			tty.Close()
			return -1, err
		}
	}
	tty.Close()
	return utils.ExitStatus(status.Sys().(syscall.WaitStatus)), nil
}

func handleSignals(container *libcontainer.Process, tty *tty) {
	sigc := make(chan os.Signal, 10)
	signal.Notify(sigc)
	tty.resize()
	for sig := range sigc {
		switch sig {
		case syscall.SIGWINCH:
			tty.resize()
		default:
			container.Signal(sig)
		}
	}
}
