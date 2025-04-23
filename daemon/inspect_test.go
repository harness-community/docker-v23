package daemon // import "github.com/harness-community/docker-v23/daemon"

import (
	"testing"

	containertypes "github.com/harness-community/docker-v23/api/types/container"
	"github.com/harness-community/docker-v23/container"
	"github.com/harness-community/docker-v23/daemon/config"
	"github.com/harness-community/docker-v23/daemon/exec"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestGetInspectData(t *testing.T) {
	c := &container.Container{
		ID:           "inspect-me",
		HostConfig:   &containertypes.HostConfig{},
		State:        container.NewState(),
		ExecCommands: exec.NewStore(),
	}

	d := &Daemon{
		linkIndex:   newLinkIndex(),
		configStore: &config.Config{},
	}

	_, err := d.getInspectData(c)
	assert.Check(t, is.ErrorContains(err, ""))

	c.Dead = true
	_, err = d.getInspectData(c)
	assert.Check(t, err)
}
