//go:build !windows
// +build !windows

package daemon // import "github.com/harness-community/docker-v23/daemon"

import (
	"github.com/harness-community/docker-v23/container"
	"github.com/harness-community/docker-v23/errdefs"
)

func (daemon *Daemon) saveAppArmorConfig(container *container.Container) error {
	container.AppArmorProfile = "" // we don't care about the previous value.

	if !daemon.RawSysInfo().AppArmor {
		return nil // if apparmor is disabled there is nothing to do here.
	}

	if err := parseSecurityOpt(container, container.HostConfig); err != nil {
		return errdefs.InvalidParameter(err)
	}

	if container.HostConfig.Privileged {
		container.AppArmorProfile = unconfinedAppArmorProfile
	} else if container.AppArmorProfile == "" {
		container.AppArmorProfile = defaultAppArmorProfile
	}
	return nil
}
