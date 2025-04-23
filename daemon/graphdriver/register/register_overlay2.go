//go:build !exclude_graphdriver_overlay2 && linux
// +build !exclude_graphdriver_overlay2,linux

package register // import "github.com/DevanshMathur19/docker-v23/daemon/graphdriver/register"

import (
	// register the overlay2 graphdriver
	_ "github.com/DevanshMathur19/docker-v23/daemon/graphdriver/overlay2"
)
