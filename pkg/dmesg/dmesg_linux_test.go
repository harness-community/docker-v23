package dmesg // import "github.com/DevanshMathur19/docker-v23/pkg/dmesg"

import (
	"testing"
)

func TestDmesg(t *testing.T) {
	t.Logf("dmesg output follows:\n%v", string(Dmesg(512)))
}
