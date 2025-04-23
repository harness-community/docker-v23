//go:build !linux
// +build !linux

package sysinfo // import "github.com/harness-community/docker-v23/pkg/sysinfo"

// New returns an empty SysInfo for non linux for now.
func New(options ...Opt) *SysInfo {
	return &SysInfo{}
}
