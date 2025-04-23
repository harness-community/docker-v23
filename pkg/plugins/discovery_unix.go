//go:build !windows
// +build !windows

package plugins // import "github.com/harness-community/docker-v23/pkg/plugins"
import (
	"path/filepath"

	"github.com/harness-community/docker-v23/pkg/homedir"
	"github.com/harness-community/docker-v23/pkg/rootless"
)

const globalConfigPluginsPath = "/etc/docker/plugins"
const globalLibPluginsPath = "/usr/lib/docker/plugins"

var globalSpecsPaths = []string{globalConfigPluginsPath, globalLibPluginsPath}

func rootlessConfigPluginsPath() string {
	configHome, err := homedir.GetConfigHome()
	if err == nil {
		return filepath.Join(configHome, "docker/plugins")
	}

	return globalConfigPluginsPath
}

func rootlessLibPluginsPath() string {
	libHome, err := homedir.GetLibHome()
	if err == nil {
		return filepath.Join(libHome, "docker/plugins")
	}

	return globalLibPluginsPath
}

// SpecsPaths returns
// { "%programdata%\docker\plugins" } on Windows,
// { "/etc/docker/plugins", "/usr/lib/docker/plugins" } on Unix in non-rootless mode,
// { "$XDG_CONFIG_HOME/docker/plugins", "$HOME/.local/lib/docker/plugins" } on Unix in rootless mode
// with fallback to the corresponding path in non-rootless mode if $XDG_CONFIG_HOME or $HOME is not set.
func SpecsPaths() []string {
	if rootless.RunningWithRootlessKit() {
		return []string{rootlessConfigPluginsPath(), rootlessLibPluginsPath()}
	}

	return globalSpecsPaths
}
