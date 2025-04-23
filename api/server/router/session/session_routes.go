package session // import "github.com/DevanshMathur19/docker-v23/api/server/router/session"

import (
	"context"
	"net/http"

	"github.com/DevanshMathur19/docker-v23/errdefs"
)

func (sr *sessionRouter) startSession(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	err := sr.backend.HandleHTTPRequest(ctx, w, r)
	if err != nil {
		return errdefs.InvalidParameter(err)
	}
	return nil
}
