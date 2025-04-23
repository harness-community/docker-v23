//go:build !linux
// +build !linux

package vfs // import "github.com/DevanshMathur19/docker-v23/daemon/graphdriver/vfs"

import (
	"github.com/DevanshMathur19/docker-v23/pkg/chrootarchive"
	"github.com/DevanshMathur19/docker-v23/pkg/idtools"
)

func dirCopy(srcDir, dstDir string) error {
	return chrootarchive.NewArchiver(idtools.IdentityMapping{}).CopyWithTar(srcDir, dstDir)
}
