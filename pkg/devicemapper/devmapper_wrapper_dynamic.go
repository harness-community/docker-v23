//go:build linux && cgo && !static_build
// +build linux,cgo,!static_build

package devicemapper // import "github.com/DevanshMathur19/docker-v23/pkg/devicemapper"

// #cgo pkg-config: devmapper
import "C"
