package libnetwork

import (
	"github.com/harness-community/docker-v23/libnetwork/drvregistry"
	"github.com/harness-community/docker-v23/libnetwork/ipamapi"
	builtinIpam "github.com/harness-community/docker-v23/libnetwork/ipams/builtin"
	nullIpam "github.com/harness-community/docker-v23/libnetwork/ipams/null"
	remoteIpam "github.com/harness-community/docker-v23/libnetwork/ipams/remote"
	"github.com/harness-community/docker-v23/libnetwork/ipamutils"
)

func initIPAMDrivers(r *drvregistry.DrvRegistry, lDs, gDs interface{}, addressPool []*ipamutils.NetworkToSplit) error {
	builtinIpam.SetDefaultIPAddressPool(addressPool)
	for _, fn := range [](func(ipamapi.Callback, interface{}, interface{}) error){
		builtinIpam.Init,
		remoteIpam.Init,
		nullIpam.Init,
	} {
		if err := fn(r, lDs, gDs); err != nil {
			return err
		}
	}

	return nil
}
