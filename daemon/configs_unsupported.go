//go:build !linux && !windows
// +build !linux,!windows

package daemon // import "github.com/DevanshMathur19/docker-v23/daemon"

func configsSupported() bool {
	return false
}
