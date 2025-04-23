package vfs // import "github.com/harness-community/docker-v23/daemon/graphdriver/vfs"

import "github.com/harness-community/docker-v23/daemon/graphdriver/copy"

func dirCopy(srcDir, dstDir string) error {
	return copy.DirCopy(srcDir, dstDir, copy.Content, false)
}
