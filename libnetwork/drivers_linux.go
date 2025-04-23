package libnetwork

import (
	"github.com/harness-community/docker-v23/libnetwork/drivers/bridge"
	"github.com/harness-community/docker-v23/libnetwork/drivers/host"
	"github.com/harness-community/docker-v23/libnetwork/drivers/ipvlan"
	"github.com/harness-community/docker-v23/libnetwork/drivers/macvlan"
	"github.com/harness-community/docker-v23/libnetwork/drivers/null"
	"github.com/harness-community/docker-v23/libnetwork/drivers/overlay"
	"github.com/harness-community/docker-v23/libnetwork/drivers/remote"
)

func getInitializers() []initializer {
	in := []initializer{
		{bridge.Init, "bridge"},
		{host.Init, "host"},
		{ipvlan.Init, "ipvlan"},
		{macvlan.Init, "macvlan"},
		{null.Init, "null"},
		{overlay.Init, "overlay"},
		{remote.Init, "remote"},
	}
	return in
}
