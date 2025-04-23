package libnetwork

import (
	"github.com/harness-community/docker-v23/libnetwork/drivers/null"
	"github.com/harness-community/docker-v23/libnetwork/drivers/remote"
)

func getInitializers() []initializer {
	return []initializer{
		{null.Init, "null"},
		{remote.Init, "remote"},
	}
}
