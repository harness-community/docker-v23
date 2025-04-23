package mounts // import "github.com/DevanshMathur19/docker-v23/volume/mounts"

func (p *linuxParser) HasResource(m *MountPoint, absolutePath string) bool {
	return false
}
