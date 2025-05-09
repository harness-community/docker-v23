package system // import "github.com/harness-community/docker-v23/api/server/router/system"

import (
	"context"
	"time"

	"github.com/harness-community/docker-v23/api/types"
	"github.com/harness-community/docker-v23/api/types/events"
	"github.com/harness-community/docker-v23/api/types/filters"
	"github.com/harness-community/docker-v23/api/types/swarm"
)

// DiskUsageOptions holds parameters for system disk usage query.
type DiskUsageOptions struct {
	// Containers controls whether container disk usage should be computed.
	Containers bool

	// Images controls whether image disk usage should be computed.
	Images bool

	// Volumes controls whether volume disk usage should be computed.
	Volumes bool
}

// Backend is the methods that need to be implemented to provide
// system specific functionality.
type Backend interface {
	SystemInfo() *types.Info
	SystemVersion() types.Version
	SystemDiskUsage(ctx context.Context, opts DiskUsageOptions) (*types.DiskUsage, error)
	SubscribeToEvents(since, until time.Time, ef filters.Args) ([]events.Message, chan interface{})
	UnsubscribeFromEvents(chan interface{})
	AuthenticateToRegistry(ctx context.Context, authConfig *types.AuthConfig) (string, string, error)
}

// ClusterBackend is all the methods that need to be implemented
// to provide cluster system specific functionality.
type ClusterBackend interface {
	Info() swarm.Info
}

// StatusProvider provides methods to get the swarm status of the current node.
type StatusProvider interface {
	Status() string
}
