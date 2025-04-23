//go:build linux && cgo && !static_build && journald && !journald_compat
// +build linux,cgo,!static_build,journald,!journald_compat

package journald // import "github.com/DevanshMathur19/docker-v23/daemon/logger/journald"

// #cgo pkg-config: libsystemd
import "C"
