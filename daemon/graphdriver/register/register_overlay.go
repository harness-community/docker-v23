//go:build !exclude_graphdriver_overlay && linux
// +build !exclude_graphdriver_overlay,linux

package register // import "github.com/harness-community/docker-v23/daemon/graphdriver/register"

import (
	// register the overlay graphdriver
	_ "github.com/harness-community/docker-v23/daemon/graphdriver/overlay"
)
