package daemon // import "github.com/harness-community/docker-v23/daemon"

import (
	"testing"

	containertypes "github.com/harness-community/docker-v23/api/types/container"
)

func TestMergeAndVerifyLogConfigNilConfig(t *testing.T) {
	d := &Daemon{defaultLogConfig: containertypes.LogConfig{Type: "json-file", Config: map[string]string{"max-file": "1"}}}
	cfg := containertypes.LogConfig{Type: d.defaultLogConfig.Type}
	if err := d.mergeAndVerifyLogConfig(&cfg); err != nil {
		t.Fatal(err)
	}
}
