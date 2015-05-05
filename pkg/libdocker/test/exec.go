package main

import (
	"fmt"
	"os"

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
	Usage:  "execute a new command inside a docker container",
	Action: execAction,
	Flags: []cli.Flag{
		cli.BoolFlag{Name: "tty,t", Usage: "allocate a TTY to the container"},
		cli.StringFlag{Name: "id", Value: "nsinit", Usage: "specify the ID for a container"},
		cli.StringFlag{Name: "user,u", Value: "root", Usage: "set the user, uid, and/or gid for the process"},
		cli.StringFlag{Name: "cwd", Value: "", Usage: "set the current working dir"},
		cli.StringSliceFlag{Name: "env", Value: standardEnvironment, Usage: "set environment variables for the process"},
	},
}

func execAction(context *cli.Context) {
	containerId := context.String("id")
	execOptions := &dl.DockerExecOptions{
		Args:   context.Args(),
		Env:    context.StringSlice("env"),
		User:   context.String("user"),
		Cwd:    context.String("cwd"),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	retCode, err := dl.RunInContainer(containerId, execOptions)
	if err != nil {
		fatal(err)
	}

	os.Exit(retCode)
}

func fatal(err error) {
	if lerr, ok := err.(libcontainer.Error); ok {
		lerr.Detail(os.Stderr)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
