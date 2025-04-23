//go:build !linux
// +build !linux

package kernel // import "github.com/DevanshMathur19/docker-v23/pkg/parsers/kernel"

import (
	"errors"
)

// Utsname represents the system name structure.
// It is defined here to make it portable as it is available on linux but not
// on windows.
type Utsname struct {
	Release [65]byte
}

func uname() (*Utsname, error) {
	return nil, errors.New("Kernel version detection is available only on linux")
}
