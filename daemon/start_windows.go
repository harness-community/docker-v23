package daemon // import "github.com/DevanshMathur19/docker-v23/daemon"

import (
	"github.com/Microsoft/hcsshim/cmd/containerd-shim-runhcs-v1/options"
	"github.com/DevanshMathur19/docker-v23/container"
	"github.com/DevanshMathur19/docker-v23/pkg/system"
)

func (daemon *Daemon) getLibcontainerdCreateOptions(_ *container.Container) (string, interface{}, error) {
	if system.ContainerdRuntimeSupported() {
		opts := &options.Options{}
		return "io.containerd.runhcs.v1", opts, nil
	}
	return "", nil, nil
}
