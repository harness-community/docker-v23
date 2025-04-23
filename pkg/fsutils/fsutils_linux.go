package fsutils // import "github.com/harness-community/docker-v23/pkg/fsutils"

import "github.com/containerd/continuity/fs"

// SupportsDType returns whether the filesystem mounted on path supports d_type.
//
// Deprecated: use github.com/containerd/continuity/fs.SupportsDType
var SupportsDType = fs.SupportsDType
