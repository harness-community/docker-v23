package daemon // import "github.com/DevanshMathur19/docker-v23/daemon"

import (
	"github.com/DevanshMathur19/docker-v23/container"
	"github.com/DevanshMathur19/docker-v23/pkg/archive"
)

func (daemon *Daemon) tarCopyOptions(container *container.Container, noOverwriteDirNonDir bool) (*archive.TarOptions, error) {
	return daemon.defaultTarCopyOptions(noOverwriteDirNonDir), nil
}
