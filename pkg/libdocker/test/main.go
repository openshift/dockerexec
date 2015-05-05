package main

import (
	dl "github.com/openshift/dockerexec/pkg/libdocker"
)

func main() {
	options := dl.DockerExecOptions{}
	dl.RunInContainer("hello", options)
}
