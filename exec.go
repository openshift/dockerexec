package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/libcontainer"
	dl "github.com/openshift/dockerexec/pkg/libdocker"
)

var standardEnvironment = &cli.StringSlice{
	"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
	"HOSTNAME=nsinit",
	"TERM=xterm",
}

var execCommand = cli.Command{
	Name:   "exec",
	Usage:  "execute a new command inside a container",
	Action: execAction,
	Flags: append([]cli.Flag{
		cli.BoolFlag{Name: "tty,t", Usage: "allocate a TTY to the container"},
		cli.StringFlag{Name: "id", Value: "", Usage: "specify the ID for a container"},
		cli.StringFlag{Name: "user,u", Value: "", Usage: "set the user, uid, and/or gid for the process"},
		cli.StringFlag{Name: "cwd", Value: "", Usage: "set the current working dir"},
		cli.StringSliceFlag{Name: "env", Value: standardEnvironment, Usage: "set environment variables for the process"},
	}, createFlags...),
}

func execAction(context *cli.Context) {
	containerId := context.String("id")

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
