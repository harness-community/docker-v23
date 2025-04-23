//go:build !linux
// +build !linux

package daemon // import "github.com/DevanshMathur19/docker-v23/daemon"

import (
	"context"

	"github.com/containerd/containerd/containers"
	coci "github.com/containerd/containerd/oci"
	"github.com/DevanshMathur19/docker-v23/container"
)

const supportsSeccomp = false

// WithSeccomp sets the seccomp profile
func WithSeccomp(daemon *Daemon, c *container.Container) coci.SpecOpts {
	return func(ctx context.Context, _ coci.Client, _ *containers.Container, s *coci.Spec) error {
		return nil
	}
}
