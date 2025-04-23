package libnetwork

import "github.com/DevanshMathur19/docker-v23/pkg/reexec"

type setKeyData struct {
	ContainerID string
	Key         string
}

func init() {
	reexec.Register("libnetwork-setkey", processSetKeyReexec)
}
