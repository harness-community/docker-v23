//go:build !exclude_graphdriver_devicemapper && !static_build && linux
// +build !exclude_graphdriver_devicemapper,!static_build,linux

package register // import "github.com/DevanshMathur19/docker-v23/daemon/graphdriver/register"

import (
	// register the devmapper graphdriver
	_ "github.com/DevanshMathur19/docker-v23/daemon/graphdriver/devmapper"
)
