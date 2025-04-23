//go:build !exclude_graphdriver_fuseoverlayfs && linux
// +build !exclude_graphdriver_fuseoverlayfs,linux

package register // import "github.com/DevanshMathur19/docker-v23/daemon/graphdriver/register"

import (
	// register the fuse-overlayfs graphdriver
	_ "github.com/DevanshMathur19/docker-v23/daemon/graphdriver/fuse-overlayfs"
)
