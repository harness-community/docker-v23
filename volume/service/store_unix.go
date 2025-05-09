//go:build linux || freebsd || darwin
// +build linux freebsd darwin

package service // import "github.com/harness-community/docker-v23/volume/service"

// normalizeVolumeName is a platform specific function to normalize the name
// of a volume. This is a no-op on Unix-like platforms
func normalizeVolumeName(name string) string {
	return name
}
