package container // import "github.com/harness-community/docker-v23/api/server/router/container"

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/harness-community/docker-v23/api/server/httputils"
	"github.com/harness-community/docker-v23/api/types"
	"github.com/harness-community/docker-v23/api/types/container"
	"github.com/harness-community/docker-v23/api/types/versions"
	"github.com/harness-community/docker-v23/errdefs"
	"github.com/harness-community/docker-v23/pkg/stdcopy"
	"github.com/sirupsen/logrus"
)

func (s *containerRouter) getExecByID(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	eConfig, err := s.backend.ContainerExecInspect(vars["id"])
	if err != nil {
		return err
	}

	return httputils.WriteJSON(w, http.StatusOK, eConfig)
}

type execCommandError struct{}

func (execCommandError) Error() string {
	return "No exec command specified"
}

func (execCommandError) InvalidParameter() {}

func (s *containerRouter) postContainerExecCreate(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}

	execConfig := &types.ExecConfig{}
	if err := httputils.ReadJSON(r, execConfig); err != nil {
		return err
	}

	if len(execConfig.Cmd) == 0 {
		return execCommandError{}
	}

	version := httputils.VersionFromContext(ctx)
	if versions.LessThan(version, "1.42") {
		// Not supported by API versions before 1.42
		execConfig.ConsoleSize = nil
	}

	// Register an instance of Exec in container.
	id, err := s.backend.ContainerExecCreate(vars["name"], execConfig)
	if err != nil {
		logrus.Errorf("Error setting up exec command in container %s: %v", vars["name"], err)
		return err
	}

	return httputils.WriteJSON(w, http.StatusCreated, &types.IDResponse{
		ID: id,
	})
}

// TODO(vishh): Refactor the code to avoid having to specify stream config as part of both create and start.
func (s *containerRouter) postContainerExecStart(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}

	version := httputils.VersionFromContext(ctx)
	if versions.LessThan(version, "1.22") {
		// API versions before 1.22 did not enforce application/json content-type.
		// Allow older clients to work by patching the content-type.
		if r.Header.Get("Content-Type") != "application/json" {
			r.Header.Set("Content-Type", "application/json")
		}
	}

	var (
		execName                  = vars["name"]
		stdin, inStream           io.ReadCloser
		stdout, stderr, outStream io.Writer
	)

	execStartCheck := &types.ExecStartCheck{}
	if err := httputils.ReadJSON(r, execStartCheck); err != nil {
		return err
	}

	if exists, err := s.backend.ExecExists(execName); !exists {
		return err
	}

	if execStartCheck.ConsoleSize != nil {
		// Not supported before 1.42
		if versions.LessThan(version, "1.42") {
			execStartCheck.ConsoleSize = nil
		}

		// No console without tty
		if !execStartCheck.Tty {
			execStartCheck.ConsoleSize = nil
		}
	}

	if !execStartCheck.Detach {
		var err error
		// Setting up the streaming http interface.
		inStream, outStream, err = httputils.HijackConnection(w)
		if err != nil {
			return err
		}
		defer httputils.CloseStreams(inStream, outStream)

		if _, ok := r.Header["Upgrade"]; ok {
			contentType := types.MediaTypeRawStream
			if !execStartCheck.Tty && versions.GreaterThanOrEqualTo(httputils.VersionFromContext(ctx), "1.42") {
				contentType = types.MediaTypeMultiplexedStream
			}
			fmt.Fprint(outStream, "HTTP/1.1 101 UPGRADED\r\nContent-Type: "+contentType+"\r\nConnection: Upgrade\r\nUpgrade: tcp\r\n")
		} else {
			fmt.Fprint(outStream, "HTTP/1.1 200 OK\r\nContent-Type: application/vnd.docker.raw-stream\r\n")
		}

		// copy headers that were removed as part of hijack
		if err := w.Header().WriteSubset(outStream, nil); err != nil {
			return err
		}
		fmt.Fprint(outStream, "\r\n")

		stdin = inStream
		stdout = outStream
		if !execStartCheck.Tty {
			stderr = stdcopy.NewStdWriter(outStream, stdcopy.Stderr)
			stdout = stdcopy.NewStdWriter(outStream, stdcopy.Stdout)
		}
	}

	options := container.ExecStartOptions{
		Stdin:       stdin,
		Stdout:      stdout,
		Stderr:      stderr,
		ConsoleSize: execStartCheck.ConsoleSize,
	}

	// Now run the user process in container.
	// Maybe we should we pass ctx here if we're not detaching?
	if err := s.backend.ContainerExecStart(context.Background(), execName, options); err != nil {
		if execStartCheck.Detach {
			return err
		}
		stdout.Write([]byte(err.Error() + "\r\n"))
		logrus.Errorf("Error running exec %s in container: %v", execName, err)
	}
	return nil
}

func (s *containerRouter) postContainerExecResize(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}
	height, err := strconv.Atoi(r.Form.Get("h"))
	if err != nil {
		return errdefs.InvalidParameter(err)
	}
	width, err := strconv.Atoi(r.Form.Get("w"))
	if err != nil {
		return errdefs.InvalidParameter(err)
	}

	return s.backend.ContainerExecResize(vars["name"], height, width)
}
