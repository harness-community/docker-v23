//go:build windows
// +build windows

package daemon // import "github.com/harness-community/docker-v23/daemon"

import "github.com/harness-community/docker-v23/pkg/plugingetter"

func registerMetricsPluginCallback(getter plugingetter.PluginGetter, sockPath string) {
}

func (daemon *Daemon) listenMetricsSock() (string, error) {
	return "", nil
}
