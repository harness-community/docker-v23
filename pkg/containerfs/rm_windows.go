package containerfs // import "github.com/harness-community/docker-v23/pkg/containerfs"

import "os"

// EnsureRemoveAll is an alias to os.RemoveAll on Windows
var EnsureRemoveAll = os.RemoveAll
