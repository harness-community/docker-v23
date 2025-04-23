//go:build !linux && !darwin && !freebsd && !windows
// +build !linux,!darwin,!freebsd,!windows

package daemon // import "github.com/harness-community/docker-v23/daemon"

func (daemon *Daemon) setupDumpStackTrap(_ string) {
	return
}
