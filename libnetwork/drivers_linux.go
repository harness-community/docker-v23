package libnetwork

import (
	"github.com/DevanshMathur19/docker-v23/libnetwork/drivers/bridge"
	"github.com/DevanshMathur19/docker-v23/libnetwork/drivers/host"
	"github.com/DevanshMathur19/docker-v23/libnetwork/drivers/ipvlan"
	"github.com/DevanshMathur19/docker-v23/libnetwork/drivers/macvlan"
	"github.com/DevanshMathur19/docker-v23/libnetwork/drivers/null"
	"github.com/DevanshMathur19/docker-v23/libnetwork/drivers/overlay"
	"github.com/DevanshMathur19/docker-v23/libnetwork/drivers/remote"
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
