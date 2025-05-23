package system // import "github.com/harness-community/docker-v23/pkg/system"
import (
	"runtime"
	"strings"
)

// IsOSSupported determines if an operating system is supported by the host.
func IsOSSupported(os string) bool {
	return strings.EqualFold(runtime.GOOS, os)
}
