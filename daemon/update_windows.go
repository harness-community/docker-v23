package daemon // import "github.com/harness-community/docker-v23/daemon"

import (
	"github.com/harness-community/docker-v23/api/types/container"
	libcontainerdtypes "github.com/harness-community/docker-v23/libcontainerd/types"
)

func toContainerdResources(resources container.Resources) *libcontainerdtypes.Resources {
	// We don't support update, so do nothing
	return nil
}
