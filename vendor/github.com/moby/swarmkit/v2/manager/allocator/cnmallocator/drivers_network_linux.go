package cnmallocator

import (
	"github.com/harness-community/docker-v23/libnetwork/drivers/bridge/brmanager"
	"github.com/harness-community/docker-v23/libnetwork/drivers/host"
	"github.com/harness-community/docker-v23/libnetwork/drivers/ipvlan/ivmanager"
	"github.com/harness-community/docker-v23/libnetwork/drivers/macvlan/mvmanager"
	"github.com/harness-community/docker-v23/libnetwork/drivers/overlay/ovmanager"
	"github.com/harness-community/docker-v23/libnetwork/drivers/remote"
	"github.com/moby/swarmkit/v2/manager/allocator/networkallocator"
)

var initializers = []initializer{
	{remote.Init, "remote"},
	{ovmanager.Init, "overlay"},
	{mvmanager.Init, "macvlan"},
	{brmanager.Init, "bridge"},
	{ivmanager.Init, "ipvlan"},
	{host.Init, "host"},
}

// PredefinedNetworks returns the list of predefined network structures
func PredefinedNetworks() []networkallocator.PredefinedNetworkData {
	return []networkallocator.PredefinedNetworkData{
		{Name: "bridge", Driver: "bridge"},
		{Name: "host", Driver: "host"},
	}
}
