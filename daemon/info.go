package daemon // import "github.com/harness-community/docker-v23/daemon"

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/harness-community/docker-v23/api"
	"github.com/harness-community/docker-v23/api/types"
	"github.com/harness-community/docker-v23/cli/debug"
	"github.com/harness-community/docker-v23/daemon/config"
	"github.com/harness-community/docker-v23/daemon/logger"
	"github.com/harness-community/docker-v23/dockerversion"
	"github.com/harness-community/docker-v23/pkg/fileutils"
	"github.com/harness-community/docker-v23/pkg/parsers/kernel"
	"github.com/harness-community/docker-v23/pkg/parsers/operatingsystem"
	"github.com/harness-community/docker-v23/pkg/platform"
	"github.com/harness-community/docker-v23/pkg/sysinfo"
	"github.com/harness-community/docker-v23/pkg/system"
	"github.com/harness-community/docker-v23/registry"
	metrics "github.com/docker/go-metrics"
	"github.com/opencontainers/selinux/go-selinux"
	"github.com/sirupsen/logrus"
)

// SystemInfo returns information about the host server the daemon is running on.
func (daemon *Daemon) SystemInfo() *types.Info {
	defer metrics.StartTimer(hostInfoFunctions.WithValues("system_info"))()

	sysInfo := daemon.RawSysInfo()

	v := &types.Info{
		ID:                 daemon.id,
		Images:             daemon.imageService.CountImages(),
		IPv4Forwarding:     !sysInfo.IPv4ForwardingDisabled,
		BridgeNfIptables:   !sysInfo.BridgeNFCallIPTablesDisabled,
		BridgeNfIP6tables:  !sysInfo.BridgeNFCallIP6TablesDisabled,
		Name:               hostName(),
		SystemTime:         time.Now().Format(time.RFC3339Nano),
		LoggingDriver:      daemon.defaultLogConfig.Type,
		KernelVersion:      kernelVersion(),
		OperatingSystem:    operatingSystem(),
		OSVersion:          osVersion(),
		IndexServerAddress: registry.IndexServer,
		OSType:             platform.OSType,
		Architecture:       platform.Architecture,
		RegistryConfig:     daemon.registryService.ServiceConfig(),
		NCPU:               sysinfo.NumCPU(),
		MemTotal:           memInfo().MemTotal,
		GenericResources:   daemon.genericResources,
		DockerRootDir:      daemon.configStore.Root,
		Labels:             daemon.configStore.Labels,
		ExperimentalBuild:  daemon.configStore.Experimental,
		ServerVersion:      dockerversion.Version,
		HTTPProxy:          config.MaskCredentials(getConfigOrEnv(daemon.configStore.HTTPProxy, "HTTP_PROXY", "http_proxy")),
		HTTPSProxy:         config.MaskCredentials(getConfigOrEnv(daemon.configStore.HTTPSProxy, "HTTPS_PROXY", "https_proxy")),
		NoProxy:            getConfigOrEnv(daemon.configStore.NoProxy, "NO_PROXY", "no_proxy"),
		LiveRestoreEnabled: daemon.configStore.LiveRestoreEnabled,
		Isolation:          daemon.defaultIsolation,
	}

	daemon.fillContainerStates(v)
	daemon.fillDebugInfo(v)
	daemon.fillAPIInfo(v)
	// Retrieve platform specific info
	daemon.fillPlatformInfo(v, sysInfo)
	daemon.fillDriverInfo(v)
	daemon.fillPluginsInfo(v)
	daemon.fillSecurityOptions(v, sysInfo)
	daemon.fillLicense(v)
	daemon.fillDefaultAddressPools(v)

	return v
}

// SystemVersion returns version information about the daemon.
func (daemon *Daemon) SystemVersion() types.Version {
	defer metrics.StartTimer(hostInfoFunctions.WithValues("system_version"))()

	kernelVersion := kernelVersion()

	v := types.Version{
		Components: []types.ComponentVersion{
			{
				Name:    "Engine",
				Version: dockerversion.Version,
				Details: map[string]string{
					"GitCommit":     dockerversion.GitCommit,
					"ApiVersion":    api.DefaultVersion,
					"MinAPIVersion": api.MinVersion,
					"GoVersion":     runtime.Version(),
					"Os":            runtime.GOOS,
					"Arch":          runtime.GOARCH,
					"BuildTime":     dockerversion.BuildTime,
					"KernelVersion": kernelVersion,
					"Experimental":  fmt.Sprintf("%t", daemon.configStore.Experimental),
				},
			},
		},

		// Populate deprecated fields for older clients
		Version:       dockerversion.Version,
		GitCommit:     dockerversion.GitCommit,
		APIVersion:    api.DefaultVersion,
		MinAPIVersion: api.MinVersion,
		GoVersion:     runtime.Version(),
		Os:            runtime.GOOS,
		Arch:          runtime.GOARCH,
		BuildTime:     dockerversion.BuildTime,
		KernelVersion: kernelVersion,
		Experimental:  daemon.configStore.Experimental,
	}

	v.Platform.Name = dockerversion.PlatformName

	daemon.fillPlatformVersion(&v)
	return v
}

func (daemon *Daemon) fillDriverInfo(v *types.Info) {
	const warnMsg = `
WARNING: The %s storage-driver is deprecated, and will be removed in a future release.
         Refer to the documentation for more information: https://docs.docker.com/go/storage-driver/`

	switch daemon.graphDriver {
	case "aufs", "devicemapper", "overlay":
		v.Warnings = append(v.Warnings, fmt.Sprintf(warnMsg, daemon.graphDriver))
	}

	v.Driver = daemon.graphDriver
	v.DriverStatus = daemon.imageService.LayerStoreStatus()

	fillDriverWarnings(v)
}

func (daemon *Daemon) fillPluginsInfo(v *types.Info) {
	v.Plugins = types.PluginsInfo{
		Volume:  daemon.volumes.GetDriverList(),
		Network: daemon.GetNetworkDriverList(),

		// The authorization plugins are returned in the order they are
		// used as they constitute a request/response modification chain.
		Authorization: daemon.configStore.AuthorizationPlugins,
		Log:           logger.ListDrivers(),
	}
}

func (daemon *Daemon) fillSecurityOptions(v *types.Info, sysInfo *sysinfo.SysInfo) {
	var securityOptions []string
	if sysInfo.AppArmor {
		securityOptions = append(securityOptions, "name=apparmor")
	}
	if sysInfo.Seccomp && supportsSeccomp {
		if daemon.seccompProfilePath != config.SeccompProfileDefault {
			v.Warnings = append(v.Warnings, "WARNING: daemon is not using the default seccomp profile")
		}
		securityOptions = append(securityOptions, "name=seccomp,profile="+daemon.seccompProfilePath)
	}
	if selinux.GetEnabled() {
		securityOptions = append(securityOptions, "name=selinux")
	}
	if rootIDs := daemon.idMapping.RootPair(); rootIDs.UID != 0 || rootIDs.GID != 0 {
		securityOptions = append(securityOptions, "name=userns")
	}
	if daemon.Rootless() {
		securityOptions = append(securityOptions, "name=rootless")
	}
	if daemon.cgroupNamespacesEnabled(sysInfo) {
		securityOptions = append(securityOptions, "name=cgroupns")
	}

	v.SecurityOptions = securityOptions
}

func (daemon *Daemon) fillContainerStates(v *types.Info) {
	cRunning, cPaused, cStopped := stateCtr.get()
	v.Containers = cRunning + cPaused + cStopped
	v.ContainersPaused = cPaused
	v.ContainersRunning = cRunning
	v.ContainersStopped = cStopped
}

// fillDebugInfo sets the current debugging state of the daemon, and additional
// debugging information, such as the number of Go-routines, and file descriptors.
//
// Note that this currently always collects the information, but the CLI only
// prints it if the daemon has debug enabled. We should consider to either make
// this information optional (cli to request "with debugging information"), or
// only collect it if the daemon has debug enabled. For the CLI code, see
// https://github.com/docker/cli/blob/v20.10.12/cli/command/system/info.go#L239-L244
func (daemon *Daemon) fillDebugInfo(v *types.Info) {
	v.Debug = debug.IsEnabled()
	v.NFd = fileutils.GetTotalUsedFds()
	v.NGoroutines = runtime.NumGoroutine()
	v.NEventsListener = daemon.EventsService.SubscribersCount()
}

func (daemon *Daemon) fillAPIInfo(v *types.Info) {
	const warn string = `
         Access to the remote API is equivalent to root access on the host. Refer
         to the 'Docker daemon attack surface' section in the documentation for
         more information: https://docs.docker.com/go/attack-surface/`

	cfg := daemon.configStore
	for _, host := range cfg.Hosts {
		// cnf.Hosts is normalized during startup, so should always have a scheme/proto
		h := strings.SplitN(host, "://", 2)
		proto := h[0]
		addr := h[1]
		if proto != "tcp" {
			continue
		}
		if cfg.TLS == nil || !*cfg.TLS {
			v.Warnings = append(v.Warnings, fmt.Sprintf("WARNING: API is accessible on http://%s without encryption.%s", addr, warn))
			continue
		}
		if cfg.TLSVerify == nil || !*cfg.TLSVerify {
			v.Warnings = append(v.Warnings, fmt.Sprintf("WARNING: API is accessible on https://%s without TLS client verification.%s", addr, warn))
			continue
		}
	}
}

func (daemon *Daemon) fillDefaultAddressPools(v *types.Info) {
	for _, pool := range daemon.configStore.DefaultAddressPools.Value() {
		v.DefaultAddressPools = append(v.DefaultAddressPools, types.NetworkAddressPool{
			Base: pool.Base,
			Size: pool.Size,
		})
	}
}

func hostName() string {
	hostname := ""
	if hn, err := os.Hostname(); err != nil {
		logrus.Warnf("Could not get hostname: %v", err)
	} else {
		hostname = hn
	}
	return hostname
}

func kernelVersion() string {
	var kernelVersion string
	if kv, err := kernel.GetKernelVersion(); err != nil {
		logrus.Warnf("Could not get kernel version: %v", err)
	} else {
		kernelVersion = kv.String()
	}
	return kernelVersion
}

func memInfo() *system.MemInfo {
	memInfo, err := system.ReadMemInfo()
	if err != nil {
		logrus.Errorf("Could not read system memory info: %v", err)
		memInfo = &system.MemInfo{}
	}
	return memInfo
}

func operatingSystem() (operatingSystem string) {
	defer metrics.StartTimer(hostInfoFunctions.WithValues("operating_system"))()

	if s, err := operatingsystem.GetOperatingSystem(); err != nil {
		logrus.Warnf("Could not get operating system name: %v", err)
	} else {
		operatingSystem = s
	}
	if inContainer, err := operatingsystem.IsContainerized(); err != nil {
		logrus.Errorf("Could not determine if daemon is containerized: %v", err)
		operatingSystem += " (error determining if containerized)"
	} else if inContainer {
		operatingSystem += " (containerized)"
	}

	return operatingSystem
}

func osVersion() (version string) {
	defer metrics.StartTimer(hostInfoFunctions.WithValues("os_version"))()

	version, err := operatingsystem.GetOperatingSystemVersion()
	if err != nil {
		logrus.Warnf("Could not get operating system version: %v", err)
	}

	return version
}

func getEnvAny(names ...string) string {
	for _, n := range names {
		if val := os.Getenv(n); val != "" {
			return val
		}
	}
	return ""
}

func getConfigOrEnv(config string, env ...string) string {
	if config != "" {
		return config
	}
	return getEnvAny(env...)
}
