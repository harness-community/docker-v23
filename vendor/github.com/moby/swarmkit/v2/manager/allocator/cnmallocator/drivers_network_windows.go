package cnmallocator

import (
	"github.com/harness-community/docker-v23/libnetwork/drivers/overlay/ovmanager"
	"github.com/harness-community/docker-v23/libnetwork/drivers/remote"
	"github.com/moby/swarmkit/v2/manager/allocator/networkallocator"
)

var initializers = []initializer{
	{remote.Init, "remote"},
	{ovmanager.Init, "overlay"},
}

// PredefinedNetworks returns the list of predefined network structures
func PredefinedNetworks() []networkallocator.PredefinedNetworkData {
	return nil
}
