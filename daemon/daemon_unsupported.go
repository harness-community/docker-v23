//go:build !linux && !freebsd && !windows
// +build !linux,!freebsd,!windows

package daemon // import "github.com/harness-community/docker-v23/daemon"

import (
	"github.com/harness-community/docker-v23/daemon/config"
	"github.com/harness-community/docker-v23/pkg/sysinfo"
)

const platformSupported = false

func setupResolvConf(config *config.Config) {
}

func getSysInfo(daemon *Daemon) *sysinfo.SysInfo {
	return sysinfo.New()
}
