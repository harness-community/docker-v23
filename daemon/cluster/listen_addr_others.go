//go:build !linux
// +build !linux

package cluster // import "github.com/DevanshMathur19/docker-v23/daemon/cluster"

import "net"

func (c *Cluster) resolveSystemAddr() (net.IP, error) {
	return c.resolveSystemAddrViaSubnetCheck()
}
