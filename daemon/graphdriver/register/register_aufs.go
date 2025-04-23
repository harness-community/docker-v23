//go:build !exclude_graphdriver_aufs && linux
// +build !exclude_graphdriver_aufs,linux

package register // import "github.com/harness-community/docker-v23/daemon/graphdriver/register"

import (
	// register the aufs graphdriver
	_ "github.com/harness-community/docker-v23/daemon/graphdriver/aufs"
)
