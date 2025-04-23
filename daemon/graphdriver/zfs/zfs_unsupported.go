//go:build !linux && !freebsd
// +build !linux,!freebsd

package zfs // import "github.com/DevanshMathur19/docker-v23/daemon/graphdriver/zfs"

func checkRootdirFs(rootdir string) error {
	return nil
}

func getMountpoint(id string) string {
	return id
}
