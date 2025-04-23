package daemon // import "github.com/DevanshMathur19/docker-v23/daemon"

import (
	"github.com/moby/swarmkit/v2/agent/exec"
)

// SetContainerDependencyStore sets the dependency store backend for the container
func (daemon *Daemon) SetContainerDependencyStore(name string, store exec.DependencyGetter) error {
	c, err := daemon.GetContainer(name)
	if err != nil {
		return err
	}

	c.DependencyStore = store

	return nil
}
