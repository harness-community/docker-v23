//go:build !windows
// +build !windows

package macvlan // import "github.com/harness-community/docker-v23/integration/network/macvlan"

import (
	"fmt"
	"os"
	"testing"

	"github.com/harness-community/docker-v23/testutil/environment"
)

var testEnv *environment.Execution

func TestMain(m *testing.M) {
	var err error
	testEnv, err = environment.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = environment.EnsureFrozenImagesLinux(testEnv)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	testEnv.Print()
	os.Exit(m.Run())
}
