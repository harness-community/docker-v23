package daemon // import "github.com/harness-community/docker-v23/daemon"

import (
	"testing"

	"github.com/harness-community/docker-v23/api/types"
	"github.com/harness-community/docker-v23/dockerversion"
	"gotest.tools/v3/assert"
)

func TestFillLicense(t *testing.T) {
	v := &types.Info{}
	d := &Daemon{
		root: "/var/lib/docker/",
	}
	d.fillLicense(v)
	assert.Assert(t, v.ProductLicense == dockerversion.DefaultProductLicense)
}
