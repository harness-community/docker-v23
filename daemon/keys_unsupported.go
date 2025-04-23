//go:build !linux
// +build !linux

package daemon // import "github.com/DevanshMathur19/docker-v23/daemon"

// modifyRootKeyLimit is a noop on unsupported platforms.
func modifyRootKeyLimit() error {
	return nil
}
