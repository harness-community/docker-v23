//go:build linux
// +build linux

package vfs // import "github.com/harness-community/docker-v23/daemon/graphdriver/vfs"

import (
	"testing"

	"github.com/harness-community/docker-v23/daemon/graphdriver/graphtest"

	"github.com/harness-community/docker-v23/pkg/reexec"
)

func init() {
	reexec.Init()
}

// This avoids creating a new driver for each test if all tests are run
// Make sure to put new tests between TestVfsSetup and TestVfsTeardown
func TestVfsSetup(t *testing.T) {
	graphtest.GetDriver(t, "vfs")
}

func TestVfsCreateEmpty(t *testing.T) {
	graphtest.DriverTestCreateEmpty(t, "vfs")
}

func TestVfsCreateBase(t *testing.T) {
	graphtest.DriverTestCreateBase(t, "vfs")
}

func TestVfsCreateSnap(t *testing.T) {
	graphtest.DriverTestCreateSnap(t, "vfs")
}

func TestVfsSetQuota(t *testing.T) {
	graphtest.DriverTestSetQuota(t, "vfs", false)
}

func TestVfsTeardown(t *testing.T) {
	graphtest.PutDriver(t)
}
