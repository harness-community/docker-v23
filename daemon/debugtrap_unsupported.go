//go:build !linux && !darwin && !freebsd && !windows
// +build !linux,!darwin,!freebsd,!windows

package daemon // import "github.com/DevanshMathur19/docker-v23/daemon"

func (daemon *Daemon) setupDumpStackTrap(_ string) {
	return
}
