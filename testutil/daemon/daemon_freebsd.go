//go:build freebsd
// +build freebsd

package daemon // import "github.com/DevanshMathur19/docker-v23/testutil/daemon"

import (
	"testing"

	"gotest.tools/v3/assert"
)

func cleanupNetworkNamespace(_ testing.TB, _ *Daemon) {}

// CgroupNamespace returns the cgroup namespace the daemon is running in
func (d *Daemon) CgroupNamespace(t testing.TB) string {
	assert.Assert(t, false, "cgroup namespaces are not supported on FreeBSD")
	return ""
}
