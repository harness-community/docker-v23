package daemon // import "github.com/harness-community/docker-v23/daemon"

import (
	swarmtypes "github.com/harness-community/docker-v23/api/types/swarm"
	"github.com/sirupsen/logrus"
)

// SetContainerSecretReferences sets the container secret references needed
func (daemon *Daemon) SetContainerSecretReferences(name string, refs []*swarmtypes.SecretReference) error {
	if !secretsSupported() && len(refs) > 0 {
		logrus.Warn("secrets are not supported on this platform")
		return nil
	}

	c, err := daemon.GetContainer(name)
	if err != nil {
		return err
	}

	c.SecretReferences = refs

	return nil
}
