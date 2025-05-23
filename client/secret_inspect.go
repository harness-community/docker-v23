package client // import "github.com/harness-community/docker-v23/client"

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/harness-community/docker-v23/api/types/swarm"
)

// SecretInspectWithRaw returns the secret information with raw data
func (cli *Client) SecretInspectWithRaw(ctx context.Context, id string) (swarm.Secret, []byte, error) {
	if err := cli.NewVersionError("1.25", "secret inspect"); err != nil {
		return swarm.Secret{}, nil, err
	}
	if id == "" {
		return swarm.Secret{}, nil, objectNotFoundError{object: "secret", id: id}
	}
	resp, err := cli.get(ctx, "/secrets/"+id, nil, nil)
	defer ensureReaderClosed(resp)
	if err != nil {
		return swarm.Secret{}, nil, err
	}

	body, err := io.ReadAll(resp.body)
	if err != nil {
		return swarm.Secret{}, nil, err
	}

	var secret swarm.Secret
	rdr := bytes.NewReader(body)
	err = json.NewDecoder(rdr).Decode(&secret)

	return secret, body, err
}
