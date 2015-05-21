package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	dl "github.com/openshift/dockerexec/pkg/libdocker"
)

var standardEnvironment = &cli.StringSlice{
	"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
	"TERM=xterm",
}

var execCommand = cli.Command{
	Name:   "exec",
	Usage:  "execute a new command inside a container",
	Action: execAction,
	Flags: []cli.Flag{
		cli.BoolFlag{Name: "tty,t", Usage: "allocate a TTY for the exec process"},
		cli.StringFlag{Name: "id", Value: "", Usage: "specify the ID of a running docker container"},
		cli.StringFlag{Name: "user,u", Value: "", Usage: "set the user, uid, and/or gid for the process"},
		cli.StringFlag{Name: "cwd", Value: "", Usage: "set the current working dir"},
		cli.StringSliceFlag{Name: "env", Value: standardEnvironment, Usage: "set environment variables for the process"},
	},
}

func getDockerPid() (int, error) {
	args := []string{"show", "-p", "MainPID", "docker.service"}
	out, err := exec.Command("systemctl", args...).Output()
	if err != nil {
		return -1, err
	}
	line := strings.TrimSpace(string(out))
	parts := strings.Split(line, "=")
	if len(parts) != 2 {
		return -1, fmt.Errorf("unexpected output trying to get docker pid: %v", line)
	}

	pid, err := strconv.Atoi(parts[1])
	if err != nil {
		return -1, fmt.Errorf("failed to parse docker pid: %v", err)
	}

	return pid, nil
}

func execAction(context *cli.Context) {
	runtime.LockOSThread()
	containerId := context.String("id")

	dockerPid, err := getDockerPid()
	if err != nil {
		fmt.Println("Error getting pid: %v", err)
	}
	fmt.Println("Docker PID: ", dockerPid)

	if containerId == "" {
		log.Fatal("Please specify a docker id")
	}

	if len(context.Args()) == 0 {
		log.Fatal("Please specify a command to run in the container")
	}

	execOptions := &dl.DockerExecOptions{
		Args:   context.Args(),
		Env:    context.StringSlice("env"),
		User:   context.String("user"),
		Cwd:    context.String("cwd"),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    context.Bool("tty"),
	}

	retCode, err := dl.RunInContainer(containerId, execOptions)
	if err != nil {
		fatal(err)
	}

	os.Exit(retCode)
}
