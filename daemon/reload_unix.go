//go:build linux || freebsd
// +build linux freebsd

package daemon // import "github.com/harness-community/docker-v23/daemon"

import (
	"bytes"
	"strconv"

	"github.com/harness-community/docker-v23/api/types"
	"github.com/harness-community/docker-v23/daemon/config"
)

// reloadPlatform updates configuration with platform specific options
// and updates the passed attributes
func (daemon *Daemon) reloadPlatform(conf *config.Config, attributes map[string]string) error {
	if err := conf.ValidatePlatformConfig(); err != nil {
		return err
	}

	if conf.IsValueSet("runtimes") {
		// Always set the default one
		conf.Runtimes[config.StockRuntimeName] = types.Runtime{Path: config.DefaultRuntimeBinary}
		if err := daemon.initRuntimes(conf.Runtimes); err != nil {
			return err
		}
		daemon.configStore.Runtimes = conf.Runtimes
	}

	if conf.DefaultRuntime != "" {
		daemon.configStore.DefaultRuntime = conf.DefaultRuntime
	}

	if conf.IsValueSet("default-shm-size") {
		daemon.configStore.ShmSize = conf.ShmSize
	}

	if conf.CgroupNamespaceMode != "" {
		daemon.configStore.CgroupNamespaceMode = conf.CgroupNamespaceMode
	}

	if conf.IpcMode != "" {
		daemon.configStore.IpcMode = conf.IpcMode
	}

	// Update attributes
	var runtimeList bytes.Buffer
	for name, rt := range daemon.configStore.Runtimes {
		if runtimeList.Len() > 0 {
			runtimeList.WriteRune(' ')
		}
		runtimeList.WriteString(name + ":" + rt.Path)
	}

	attributes["runtimes"] = runtimeList.String()
	attributes["default-runtime"] = daemon.configStore.DefaultRuntime
	attributes["default-shm-size"] = strconv.FormatInt(int64(daemon.configStore.ShmSize), 10)
	attributes["default-ipc-mode"] = daemon.configStore.IpcMode
	attributes["default-cgroupns-mode"] = daemon.configStore.CgroupNamespaceMode

	return nil
}
