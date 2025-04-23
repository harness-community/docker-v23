//go:build !exclude_graphdriver_btrfs && linux
// +build !exclude_graphdriver_btrfs,linux

package register // import "github.com/DevanshMathur19/docker-v23/daemon/graphdriver/register"

import (
	// register the btrfs graphdriver
	_ "github.com/DevanshMathur19/docker-v23/daemon/graphdriver/btrfs"
)
