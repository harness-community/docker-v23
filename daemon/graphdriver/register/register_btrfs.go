//go:build !exclude_graphdriver_btrfs && linux
// +build !exclude_graphdriver_btrfs,linux

package register // import "github.com/harness-community/docker-v23/daemon/graphdriver/register"

import (
	// register the btrfs graphdriver
	_ "github.com/harness-community/docker-v23/daemon/graphdriver/btrfs"
)
