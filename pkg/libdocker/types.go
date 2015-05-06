package libdocker

import (
	"io"
)

type DockerExecOptions struct {
	// The command to be run followed by any arguments.
	Args []string

	// Env specifies the environment variables for the process.
	Env []string

	// User will set the uid and gid of the executing process running inside the container
	// local to the container's user and group configuration.
	User string

	// Cwd will change the processes current working directory inside the container's rootfs.
	Cwd string

	// Stdin is a pointer to a reader which provides the standard input stream.
	Stdin io.Reader

	// Stdout is a pointer to a writer which receives the standard output stream.
	Stdout io.Writer

	// Stderr is a pointer to a writer which receives the standard error stream.
	Stderr io.Writer

	// Capabilities specify the capabilities to keep when executing the process inside the container.
	Capabilities []string

	// Tty specifies whether a tty should be allocated or not.
	Tty bool
}

// The structs below this line are used for reading docker state
type State struct {
	Running    bool
	Paused     bool
	Restarting bool
	OOMKilled  bool
	Dead       bool
	Pid        int
	ExitCode   int
	Error      string // contains last known error when starting the container
}

type Config struct {
	Hostname        string
	Domainname      string
	User            string
	AttachStdin     bool
	AttachStdout    bool
	AttachStderr    bool
	PortSpecs       []string // Deprecated - Can be in the format of 8080/tcp
	Tty             bool     // Attach standard streams to a tty, including stdin if it is not closed.
	OpenStdin       bool     // Open stdin
	StdinOnce       bool     // If true, close stdin after the 1 attached client disconnects.
	Env             []string
	Image           string // Name of the image as it was passed by the operator (eg. could be symbolic)
	Volumes         map[string]struct{}
	WorkingDir      string
	NetworkDisabled bool
	MacAddress      string
	OnBuild         []string
	Labels          map[string]string
}

type ContainerConfig struct {
	State  `json:"State"`
	ID     string
	Path   string
	Args   []string
	Config Config `json:"Config"`
}
