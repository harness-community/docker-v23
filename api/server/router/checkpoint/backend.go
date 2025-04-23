package checkpoint // import "github.com/DevanshMathur19/docker-v23/api/server/router/checkpoint"

import "github.com/DevanshMathur19/docker-v23/api/types"

// Backend for Checkpoint
type Backend interface {
	CheckpointCreate(container string, config types.CheckpointCreateOptions) error
	CheckpointDelete(container string, config types.CheckpointDeleteOptions) error
	CheckpointList(container string, config types.CheckpointListOptions) ([]types.Checkpoint, error)
}
