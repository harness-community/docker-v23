//go:build !windows
// +build !windows

package container // import "github.com/harness-community/docker-v23/daemon/cluster/executor/container"

const (
	testAbsPath        = "/foo"
	testAbsNonExistent = "/some-non-existing-host-path/"
)
