package client // import "github.com/harness-community/docker-v23/client"

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/harness-community/docker-v23/api/types/swarm"
	"github.com/harness-community/docker-v23/errdefs"
	"github.com/pkg/errors"
)

func TestTaskInspectError(t *testing.T) {
	client := &Client{
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	}

	_, _, err := client.TaskInspectWithRaw(context.Background(), "nothing")
	if !errdefs.IsSystem(err) {
		t.Fatalf("expected a Server Error, got %[1]T: %[1]v", err)
	}
}

func TestTaskInspectWithEmptyID(t *testing.T) {
	client := &Client{
		client: newMockClient(func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("should not make request")
		}),
	}
	_, _, err := client.TaskInspectWithRaw(context.Background(), "")
	if !IsErrNotFound(err) {
		t.Fatalf("Expected NotFoundError, got %v", err)
	}
}

func TestTaskInspect(t *testing.T) {
	expectedURL := "/tasks/task_id"
	client := &Client{
		client: newMockClient(func(req *http.Request) (*http.Response, error) {
			if !strings.HasPrefix(req.URL.Path, expectedURL) {
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			}
			content, err := json.Marshal(swarm.Task{
				ID: "task_id",
			})
			if err != nil {
				return nil, err
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(content)),
			}, nil
		}),
	}

	taskInspect, _, err := client.TaskInspectWithRaw(context.Background(), "task_id")
	if err != nil {
		t.Fatal(err)
	}
	if taskInspect.ID != "task_id" {
		t.Fatalf("expected `task_id`, got %s", taskInspect.ID)
	}
}
