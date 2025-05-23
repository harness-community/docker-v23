package plugin // import "github.com/harness-community/docker-v23/plugin"

import (
	"fmt"

	v2 "github.com/harness-community/docker-v23/plugin/v2"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func (pm *Manager) enable(p *v2.Plugin, c *controller, force bool) error {
	return fmt.Errorf("Not implemented")
}

func (pm *Manager) initSpec(p *v2.Plugin) (*specs.Spec, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (pm *Manager) disable(p *v2.Plugin, c *controller) error {
	return fmt.Errorf("Not implemented")
}

func (pm *Manager) restore(p *v2.Plugin, c *controller) error {
	return fmt.Errorf("Not implemented")
}

// Shutdown plugins
func (pm *Manager) Shutdown() {
}

func recursiveUnmount(_ string) error {
	return nil
}
