package container // import "github.com/harness-community/docker-v23/container"

// Mount contains information for a mount operation.
type Mount struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Writable    bool   `json:"writable"`
}
