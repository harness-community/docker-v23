package daemon // import "github.com/harness-community/docker-v23/daemon"

import (
	"context"

	"github.com/harness-community/docker-v23/api/types"
	"github.com/harness-community/docker-v23/dockerversion"
)

// AuthenticateToRegistry checks the validity of credentials in authConfig
func (daemon *Daemon) AuthenticateToRegistry(ctx context.Context, authConfig *types.AuthConfig) (string, string, error) {
	return daemon.registryService.Auth(ctx, authConfig, dockerversion.DockerUserAgent(ctx))
}
