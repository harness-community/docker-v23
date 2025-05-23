//go:build !windows
// +build !windows

package system // import "github.com/harness-community/docker-v23/pkg/system"

import (
	"time"
)

// setCTime will set the create time on a file. On Unix, the create
// time is updated as a side effect of setting the modified time, so
// no action is required.
func setCTime(path string, ctime time.Time) error {
	return nil
}
