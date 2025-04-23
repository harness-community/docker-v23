//go:build windows
// +build windows

package runconfig // import "github.com/harness-community/docker-v23/runconfig"

import (
	"testing"

	"github.com/harness-community/docker-v23/api/types/container"
)

func TestValidatePrivileged(t *testing.T) {
	expected := "Windows does not support privileged mode"
	err := validatePrivileged(&container.HostConfig{Privileged: true})
	if err == nil || err.Error() != expected {
		t.Fatalf("Expected %s", expected)
	}
}
