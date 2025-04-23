//go:build !exclude_graphdriver_fuseoverlayfs && linux
// +build !exclude_graphdriver_fuseoverlayfs,linux

package register // import "github.com/harness-community/docker-v23/daemon/graphdriver/register"

import (
	// register the fuse-overlayfs graphdriver
	_ "github.com/harness-community/docker-v23/daemon/graphdriver/fuse-overlayfs"
)
