package checkpoint // import "github.com/harness-community/docker-v23/api/server/router/checkpoint"

import (
	"context"
	"net/http"

	"github.com/harness-community/docker-v23/api/server/httputils"
	"github.com/harness-community/docker-v23/api/types"
)

func (s *checkpointRouter) postContainerCheckpoint(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}

	var options types.CheckpointCreateOptions
	if err := httputils.ReadJSON(r, &options); err != nil {
		return err
	}

	err := s.backend.CheckpointCreate(vars["name"], options)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}

func (s *checkpointRouter) getContainerCheckpoints(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}

	checkpoints, err := s.backend.CheckpointList(vars["name"], types.CheckpointListOptions{
		CheckpointDir: r.Form.Get("dir"),
	})

	if err != nil {
		return err
	}

	return httputils.WriteJSON(w, http.StatusOK, checkpoints)
}

func (s *checkpointRouter) deleteContainerCheckpoint(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}

	err := s.backend.CheckpointDelete(vars["name"], types.CheckpointDeleteOptions{
		CheckpointDir: r.Form.Get("dir"),
		CheckpointID:  vars["checkpoint"],
	})

	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
