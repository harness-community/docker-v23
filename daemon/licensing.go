package daemon // import "github.com/harness-community/docker-v23/daemon"

import (
	"github.com/harness-community/docker-v23/api/types"
	"github.com/harness-community/docker-v23/dockerversion"
)

func (daemon *Daemon) fillLicense(v *types.Info) {
	v.ProductLicense = dockerversion.DefaultProductLicense
}
