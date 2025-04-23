package mounts // import "github.com/harness-community/docker-v23/volume/mounts"

func (p *linuxParser) HasResource(m *MountPoint, absolutePath string) bool {
	return false
}
