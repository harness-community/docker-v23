//go:build !windows
// +build !windows

package containerfs // import "github.com/DevanshMathur19/docker-v23/pkg/containerfs"

import "path/filepath"

// cleanScopedPath preappends a to combine with a mnt path.
func cleanScopedPath(path string) string {
	return filepath.Join(string(filepath.Separator), path)
}
