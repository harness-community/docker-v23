package daemon // import "github.com/harness-community/docker-v23/daemon"

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/harness-community/docker-v23/daemon/config"
	"github.com/sirupsen/logrus"
)

// Reload reads configuration changes and modifies the
// daemon according to those changes.
// These are the settings that Reload changes:
// - Platform runtime
// - Daemon debug log level
// - Daemon max concurrent downloads
// - Daemon max concurrent uploads
// - Daemon max download attempts
// - Daemon shutdown timeout (in seconds)
// - Cluster discovery (reconfigure and restart)
// - Daemon labels
// - Insecure registries
// - Registry mirrors
// - Daemon live restore
func (daemon *Daemon) Reload(conf *config.Config) (err error) {
	daemon.configStore.Lock()
	attributes := map[string]string{}

	defer func() {
		if err == nil {
			jsonString, _ := json.Marshal(&struct {
				*config.Config
				config.Proxies `json:"proxies"`
			}{
				Config: daemon.configStore,
				Proxies: config.Proxies{
					HTTPProxy:  config.MaskCredentials(daemon.configStore.HTTPProxy),
					HTTPSProxy: config.MaskCredentials(daemon.configStore.HTTPSProxy),
					NoProxy:    config.MaskCredentials(daemon.configStore.NoProxy),
				},
			})
			logrus.Infof("Reloaded configuration: %s", jsonString)
		}
		daemon.configStore.Unlock()
		if err == nil {
			daemon.LogDaemonEventWithAttributes("reload", attributes)
		}
	}()

	if err := daemon.reloadPlatform(conf, attributes); err != nil {
		return err
	}
	daemon.reloadDebug(conf, attributes)
	daemon.reloadMaxConcurrentDownloadsAndUploads(conf, attributes)
	daemon.reloadMaxDownloadAttempts(conf, attributes)
	daemon.reloadShutdownTimeout(conf, attributes)
	daemon.reloadFeatures(conf, attributes)

	if err := daemon.reloadLabels(conf, attributes); err != nil {
		return err
	}
	if err := daemon.reloadAllowNondistributableArtifacts(conf, attributes); err != nil {
		return err
	}
	if err := daemon.reloadInsecureRegistries(conf, attributes); err != nil {
		return err
	}
	if err := daemon.reloadRegistryMirrors(conf, attributes); err != nil {
		return err
	}
	if err := daemon.reloadLiveRestore(conf, attributes); err != nil {
		return err
	}
	return daemon.reloadNetworkDiagnosticPort(conf, attributes)
}

// reloadDebug updates configuration with Debug option
// and updates the passed attributes
func (daemon *Daemon) reloadDebug(conf *config.Config, attributes map[string]string) {
	// update corresponding configuration
	if conf.IsValueSet("debug") {
		daemon.configStore.Debug = conf.Debug
	}
	// prepare reload event attributes with updatable configurations
	attributes["debug"] = strconv.FormatBool(daemon.configStore.Debug)
}

// reloadMaxConcurrentDownloadsAndUploads updates configuration with max concurrent
// download and upload options and updates the passed attributes
func (daemon *Daemon) reloadMaxConcurrentDownloadsAndUploads(conf *config.Config, attributes map[string]string) {
	// We always "reset" as the cost is lightweight and easy to maintain.
	daemon.configStore.MaxConcurrentDownloads = config.DefaultMaxConcurrentDownloads
	daemon.configStore.MaxConcurrentUploads = config.DefaultMaxConcurrentUploads

	if conf.IsValueSet("max-concurrent-downloads") && conf.MaxConcurrentDownloads != 0 {
		daemon.configStore.MaxConcurrentDownloads = conf.MaxConcurrentDownloads
	}
	if conf.IsValueSet("max-concurrent-uploads") && conf.MaxConcurrentUploads != 0 {
		daemon.configStore.MaxConcurrentUploads = conf.MaxConcurrentUploads
	}
	if daemon.imageService != nil {
		daemon.imageService.UpdateConfig(
			daemon.configStore.MaxConcurrentDownloads,
			daemon.configStore.MaxConcurrentUploads,
		)
	}

	// prepare reload event attributes with updatable configurations
	attributes["max-concurrent-downloads"] = strconv.Itoa(daemon.configStore.MaxConcurrentDownloads)
	attributes["max-concurrent-uploads"] = strconv.Itoa(daemon.configStore.MaxConcurrentUploads)
	logrus.Debug("Reset Max Concurrent Downloads: ", attributes["max-concurrent-downloads"])
	logrus.Debug("Reset Max Concurrent Uploads: ", attributes["max-concurrent-uploads"])
}

// reloadMaxDownloadAttempts updates configuration with max concurrent
// download attempts when a connection is lost and updates the passed attributes
func (daemon *Daemon) reloadMaxDownloadAttempts(conf *config.Config, attributes map[string]string) {
	// We always "reset" as the cost is lightweight and easy to maintain.
	daemon.configStore.MaxDownloadAttempts = config.DefaultDownloadAttempts
	if conf.IsValueSet("max-download-attempts") && conf.MaxDownloadAttempts != 0 {
		daemon.configStore.MaxDownloadAttempts = conf.MaxDownloadAttempts
	}

	// prepare reload event attributes with updatable configurations
	attributes["max-download-attempts"] = strconv.Itoa(daemon.configStore.MaxDownloadAttempts)
	logrus.Debug("Reset Max Download Attempts: ", attributes["max-download-attempts"])
}

// reloadShutdownTimeout updates configuration with daemon shutdown timeout option
// and updates the passed attributes
func (daemon *Daemon) reloadShutdownTimeout(conf *config.Config, attributes map[string]string) {
	// update corresponding configuration
	if conf.IsValueSet("shutdown-timeout") {
		daemon.configStore.ShutdownTimeout = conf.ShutdownTimeout
		logrus.Debugf("Reset Shutdown Timeout: %d", daemon.configStore.ShutdownTimeout)
	}

	// prepare reload event attributes with updatable configurations
	attributes["shutdown-timeout"] = strconv.Itoa(daemon.configStore.ShutdownTimeout)
}

// reloadLabels updates configuration with engine labels
// and updates the passed attributes
func (daemon *Daemon) reloadLabels(conf *config.Config, attributes map[string]string) error {
	// update corresponding configuration
	if conf.IsValueSet("labels") {
		daemon.configStore.Labels = conf.Labels
	}

	// prepare reload event attributes with updatable configurations
	if daemon.configStore.Labels != nil {
		labels, err := json.Marshal(daemon.configStore.Labels)
		if err != nil {
			return err
		}
		attributes["labels"] = string(labels)
	} else {
		attributes["labels"] = "[]"
	}

	return nil
}

// reloadAllowNondistributableArtifacts updates the configuration with allow-nondistributable-artifacts options
// and updates the passed attributes.
func (daemon *Daemon) reloadAllowNondistributableArtifacts(conf *config.Config, attributes map[string]string) error {
	// Update corresponding configuration.
	if conf.IsValueSet("allow-nondistributable-artifacts") {
		daemon.configStore.AllowNondistributableArtifacts = conf.AllowNondistributableArtifacts
		if err := daemon.registryService.LoadAllowNondistributableArtifacts(conf.AllowNondistributableArtifacts); err != nil {
			return err
		}
	}

	// Prepare reload event attributes with updatable configurations.
	if daemon.configStore.AllowNondistributableArtifacts != nil {
		v, err := json.Marshal(daemon.configStore.AllowNondistributableArtifacts)
		if err != nil {
			return err
		}
		attributes["allow-nondistributable-artifacts"] = string(v)
	} else {
		attributes["allow-nondistributable-artifacts"] = "[]"
	}

	return nil
}

// reloadInsecureRegistries updates configuration with insecure registry option
// and updates the passed attributes
func (daemon *Daemon) reloadInsecureRegistries(conf *config.Config, attributes map[string]string) error {
	// update corresponding configuration
	if conf.IsValueSet("insecure-registries") {
		daemon.configStore.InsecureRegistries = conf.InsecureRegistries
		if err := daemon.registryService.LoadInsecureRegistries(conf.InsecureRegistries); err != nil {
			return err
		}
	}

	// prepare reload event attributes with updatable configurations
	if daemon.configStore.InsecureRegistries != nil {
		insecureRegistries, err := json.Marshal(daemon.configStore.InsecureRegistries)
		if err != nil {
			return err
		}
		attributes["insecure-registries"] = string(insecureRegistries)
	} else {
		attributes["insecure-registries"] = "[]"
	}

	return nil
}

// reloadRegistryMirrors updates configuration with registry mirror options
// and updates the passed attributes
func (daemon *Daemon) reloadRegistryMirrors(conf *config.Config, attributes map[string]string) error {
	// update corresponding configuration
	if conf.IsValueSet("registry-mirrors") {
		daemon.configStore.Mirrors = conf.Mirrors
		if err := daemon.registryService.LoadMirrors(conf.Mirrors); err != nil {
			return err
		}
	}

	// prepare reload event attributes with updatable configurations
	if daemon.configStore.Mirrors != nil {
		mirrors, err := json.Marshal(daemon.configStore.Mirrors)
		if err != nil {
			return err
		}
		attributes["registry-mirrors"] = string(mirrors)
	} else {
		attributes["registry-mirrors"] = "[]"
	}

	return nil
}

// reloadLiveRestore updates configuration with live restore option
// and updates the passed attributes
func (daemon *Daemon) reloadLiveRestore(conf *config.Config, attributes map[string]string) error {
	// update corresponding configuration
	if conf.IsValueSet("live-restore") {
		daemon.configStore.LiveRestoreEnabled = conf.LiveRestoreEnabled
	}

	// prepare reload event attributes with updatable configurations
	attributes["live-restore"] = strconv.FormatBool(daemon.configStore.LiveRestoreEnabled)
	return nil
}

// reloadNetworkDiagnosticPort updates the network controller starting the diagnostic if the config is valid
func (daemon *Daemon) reloadNetworkDiagnosticPort(conf *config.Config, attributes map[string]string) error {
	if conf == nil || daemon.netController == nil || !conf.IsValueSet("network-diagnostic-port") ||
		conf.NetworkDiagnosticPort < 1 || conf.NetworkDiagnosticPort > 65535 {
		// If there is no config make sure that the diagnostic is off
		if daemon.netController != nil {
			daemon.netController.StopDiagnostic()
		}
		return nil
	}
	// Enable the network diagnostic if the flag is set with a valid port within the range
	logrus.WithFields(logrus.Fields{"port": conf.NetworkDiagnosticPort, "ip": "127.0.0.1"}).Warn("Starting network diagnostic server")
	daemon.netController.StartDiagnostic(conf.NetworkDiagnosticPort)

	return nil
}

// reloadFeatures updates configuration with enabled/disabled features
func (daemon *Daemon) reloadFeatures(conf *config.Config, attributes map[string]string) {
	// update corresponding configuration
	// note that we allow features option to be entirely unset
	daemon.configStore.Features = conf.Features

	// prepare reload event attributes with updatable configurations
	attributes["features"] = fmt.Sprintf("%v", daemon.configStore.Features)
}
