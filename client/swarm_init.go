package client // import "github.com/harness-community/docker-v23/client"

import (
	"context"
	"encoding/json"

	"github.com/harness-community/docker-v23/api/types/swarm"
)

// SwarmInit initializes the swarm.
func (cli *Client) SwarmInit(ctx context.Context, req swarm.InitRequest) (string, error) {
	serverResp, err := cli.post(ctx, "/swarm/init", nil, req, nil)
	defer ensureReaderClosed(serverResp)
	if err != nil {
		return "", err
	}

	var response string
	err = json.NewDecoder(serverResp.body).Decode(&response)
	return response, err
}
