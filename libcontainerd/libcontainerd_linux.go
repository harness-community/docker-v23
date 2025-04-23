package libcontainerd // import "github.com/DevanshMathur19/docker-v23/libcontainerd"

import (
	"context"

	"github.com/containerd/containerd"
	"github.com/DevanshMathur19/docker-v23/libcontainerd/remote"
	libcontainerdtypes "github.com/DevanshMathur19/docker-v23/libcontainerd/types"
)

// NewClient creates a new libcontainerd client from a containerd client
func NewClient(ctx context.Context, cli *containerd.Client, stateDir, ns string, b libcontainerdtypes.Backend) (libcontainerdtypes.Client, error) {
	return remote.NewClient(ctx, cli, stateDir, ns, b)
}
