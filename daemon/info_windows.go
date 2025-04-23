package daemon // import "github.com/harness-community/docker-v23/daemon"

import (
	"github.com/harness-community/docker-v23/api/types"
	"github.com/harness-community/docker-v23/pkg/sysinfo"
)

// fillPlatformInfo fills the platform related info.
func (daemon *Daemon) fillPlatformInfo(v *types.Info, sysInfo *sysinfo.SysInfo) {
}

func (daemon *Daemon) fillPlatformVersion(v *types.Version) {}

func fillDriverWarnings(v *types.Info) {
}

func (daemon *Daemon) cgroupNamespacesEnabled(sysInfo *sysinfo.SysInfo) bool {
	return false
}

// Rootless returns true if daemon is running in rootless mode
func (daemon *Daemon) Rootless() bool {
	return false
}
