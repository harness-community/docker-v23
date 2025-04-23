//go:build !linux && !windows
// +build !linux,!windows

package daemon // import "github.com/harness-community/docker-v23/daemon"

func secretsSupported() bool {
	return false
}
