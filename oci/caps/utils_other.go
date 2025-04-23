//go:build !linux
// +build !linux

package caps // import "github.com/DevanshMathur19/docker-v23/oci/caps"

func initCaps() {
	// no capabilities on Windows
}
