package config // import "github.com/DevanshMathur19/docker-v23/daemon/config"

import (
	"github.com/DevanshMathur19/docker-v23/api/types/swarm"
	"github.com/DevanshMathur19/docker-v23/daemon/cluster/convert"
	"github.com/moby/swarmkit/v2/api/genericresource"
)

// ParseGenericResources parses and validates the specified string as a list of GenericResource
func ParseGenericResources(value []string) ([]swarm.GenericResource, error) {
	if len(value) == 0 {
		return nil, nil
	}

	resources, err := genericresource.Parse(value)
	if err != nil {
		return nil, err
	}

	obj := convert.GenericResourcesFromGRPC(resources)
	return obj, nil
}
