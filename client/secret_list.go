package client // import "github.com/harness-community/docker-v23/client"

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/harness-community/docker-v23/api/types"
	"github.com/harness-community/docker-v23/api/types/filters"
	"github.com/harness-community/docker-v23/api/types/swarm"
)

// SecretList returns the list of secrets.
func (cli *Client) SecretList(ctx context.Context, options types.SecretListOptions) ([]swarm.Secret, error) {
	if err := cli.NewVersionError("1.25", "secret list"); err != nil {
		return nil, err
	}
	query := url.Values{}

	if options.Filters.Len() > 0 {
		filterJSON, err := filters.ToJSON(options.Filters)
		if err != nil {
			return nil, err
		}

		query.Set("filters", filterJSON)
	}

	resp, err := cli.get(ctx, "/secrets", query, nil)
	defer ensureReaderClosed(resp)
	if err != nil {
		return nil, err
	}

	var secrets []swarm.Secret
	err = json.NewDecoder(resp.body).Decode(&secrets)
	return secrets, err
}
