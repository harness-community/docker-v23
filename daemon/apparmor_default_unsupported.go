//go:build !linux
// +build !linux

package daemon // import "github.com/DevanshMathur19/docker-v23/daemon"

func ensureDefaultAppArmorProfile() error {
	return nil
}

// DefaultApparmorProfile returns an empty string.
func DefaultApparmorProfile() string {
	return ""
}
