package libdocker

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/docker/libcontainer"
	"github.com/docker/libcontainer/cgroups/systemd"
	"github.com/docker/libcontainer/configs"
)

const (
	dockerRoot    = "/var/run/docker/execdriver/native"
	stateFilename = "state.json"
)

func loadDockerFactory() (libcontainer.Factory, error) {
	cgm := libcontainer.Cgroupfs
	if systemd.UseSystemd() {
		cgm = libcontainer.SystemdCgroups
	}
	return libcontainer.New(dockerRoot, cgm)
}

func loadState(root string) (*libcontainer.State, error) {
	f, err := os.Open(filepath.Join(root, stateFilename))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		return nil, err
	}
	defer f.Close()
	var state *libcontainer.State
	if err := json.NewDecoder(f).Decode(&state); err != nil {
		return nil, err
	}
	return state, nil
}

func loadDockerConfig(containerId string) (*configs.Config, error) {
	containerRoot := filepath.Join(dockerRoot, containerId)
	state, err := loadState(containerRoot)
	if err != nil {
		return nil, err
	}
	config := &state.Config
	//modify(config, context)
	return config, nil
}
