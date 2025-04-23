package libnetwork

import (
	"github.com/DevanshMathur19/docker-v23/libnetwork/drivers/null"
	"github.com/DevanshMathur19/docker-v23/libnetwork/drivers/remote"
)

func getInitializers() []initializer {
	return []initializer{
		{null.Init, "null"},
		{remote.Init, "remote"},
	}
}
