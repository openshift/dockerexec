package libdocker

import (
	"fmt"
	"path/filepath"
)

// expandContainerId takes in a prefix of a container id
// and attempts to translate it to a unique container id
func expandContainerId(id string) (string, error) {
	idPath := filepath.Join(containerRoot, id)
	globPattern := fmt.Sprintf("%s*", idPath)
	matches, err := filepath.Glob(globPattern)
	if err != nil {
		return "", err
	}
	if len(matches) > 1 {
		return "", fmt.Errorf("More than one match found for the id, please provide a longer prefix.")
	}
	return filepath.Base(matches[0]), nil
}
