//go:build !linux && !windows
// +build !linux,!windows

package sysinfo // import "github.com/harness-community/docker-v23/pkg/sysinfo"

import (
	"runtime"
)

// NumCPU returns the number of CPUs
func NumCPU() int {
	return runtime.NumCPU()
}
