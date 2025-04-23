package libnetwork

import (
	"github.com/DevanshMathur19/docker-v23/libnetwork/drivers/null"
	"github.com/DevanshMathur19/docker-v23/libnetwork/drivers/remote"
	"github.com/DevanshMathur19/docker-v23/libnetwork/drivers/windows"
	"github.com/DevanshMathur19/docker-v23/libnetwork/drivers/windows/overlay"
)

func getInitializers() []initializer {
	return []initializer{
		{null.Init, "null"},
		{overlay.Init, "overlay"},
		{remote.Init, "remote"},
		{windows.GetInit("transparent"), "transparent"},
		{windows.GetInit("l2bridge"), "l2bridge"},
		{windows.GetInit("l2tunnel"), "l2tunnel"},
		{windows.GetInit("nat"), "nat"},
		{windows.GetInit("internal"), "internal"},
		{windows.GetInit("private"), "private"},
		{windows.GetInit("ics"), "ics"},
	}
}
