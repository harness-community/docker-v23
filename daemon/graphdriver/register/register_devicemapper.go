//go:build !exclude_graphdriver_devicemapper && !static_build && linux
// +build !exclude_graphdriver_devicemapper,!static_build,linux

package register // import "github.com/harness-community/docker-v23/daemon/graphdriver/register"

import (
	// register the devmapper graphdriver
	_ "github.com/harness-community/docker-v23/daemon/graphdriver/devmapper"
)
