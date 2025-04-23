//go:build !linux && !windows
// +build !linux,!windows

package service // import "github.com/DevanshMathur19/docker-v23/volume/service"

import (
	"github.com/DevanshMathur19/docker-v23/pkg/idtools"
	"github.com/DevanshMathur19/docker-v23/volume/drivers"
)

func setupDefaultDriver(_ *drivers.Store, _ string, _ idtools.Identity) error { return nil }
