//go:build !windows
// +build !windows

package dockerfile // import "github.com/harness-community/docker-v23/builder/dockerfile"

func defaultShellForOS(os string) []string {
	return []string{"/bin/sh", "-c"}
}
