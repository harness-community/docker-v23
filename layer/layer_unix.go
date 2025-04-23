//go:build linux || freebsd || darwin || openbsd
// +build linux freebsd darwin openbsd

package layer // import "github.com/harness-community/docker-v23/layer"

import "github.com/harness-community/docker-v23/pkg/stringid"

func (ls *layerStore) mountID(name string) string {
	return stringid.GenerateRandomID()
}
