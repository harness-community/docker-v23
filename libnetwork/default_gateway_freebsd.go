package libnetwork

import "github.com/harness-community/docker-v23/libnetwork/types"

const libnGWNetwork = "docker_gwbridge"

func getPlatformOption() EndpointOption {
	return nil
}

func (c *controller) createGWNetwork() (Network, error) {
	return nil, types.NotImplementedErrorf("default gateway functionality is not implemented in freebsd")
}
