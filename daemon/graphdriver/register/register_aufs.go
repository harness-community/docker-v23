//go:build !exclude_graphdriver_aufs && linux
// +build !exclude_graphdriver_aufs,linux

package register // import "github.com/DevanshMathur19/docker-v23/daemon/graphdriver/register"

import (
	// register the aufs graphdriver
	_ "github.com/DevanshMathur19/docker-v23/daemon/graphdriver/aufs"
)
