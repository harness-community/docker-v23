package client // import "github.com/harness-community/docker-v23/client"

import (
	"fmt"

	"github.com/harness-community/docker-v23/api/types/versions"
	"github.com/harness-community/docker-v23/errdefs"
	"github.com/pkg/errors"
)

// errConnectionFailed implements an error returned when connection failed.
type errConnectionFailed struct {
	host string
}

// Error returns a string representation of an errConnectionFailed
func (err errConnectionFailed) Error() string {
	if err.host == "" {
		return "Cannot connect to the Docker daemon. Is the docker daemon running on this host?"
	}
	return fmt.Sprintf("Cannot connect to the Docker daemon at %s. Is the docker daemon running?", err.host)
}

// IsErrConnectionFailed returns true if the error is caused by connection failed.
func IsErrConnectionFailed(err error) bool {
	return errors.As(err, &errConnectionFailed{})
}

// ErrorConnectionFailed returns an error with host in the error message when connection to docker daemon failed.
func ErrorConnectionFailed(host string) error {
	return errConnectionFailed{host: host}
}

// Deprecated: use the errdefs.NotFound() interface instead. Kept for backward compatibility
type notFound interface {
	error
	NotFound() bool
}

// IsErrNotFound returns true if the error is a NotFound error, which is returned
// by the API when some object is not found.
func IsErrNotFound(err error) bool {
	if errdefs.IsNotFound(err) {
		return true
	}
	var e notFound
	return errors.As(err, &e)
}

type objectNotFoundError struct {
	object string
	id     string
}

func (e objectNotFoundError) NotFound() {}

func (e objectNotFoundError) Error() string {
	return fmt.Sprintf("Error: No such %s: %s", e.object, e.id)
}

// IsErrUnauthorized returns true if the error is caused
// when a remote registry authentication fails
//
// Deprecated: use errdefs.IsUnauthorized
func IsErrUnauthorized(err error) bool {
	return errdefs.IsUnauthorized(err)
}

type pluginPermissionDenied struct {
	name string
}

func (e pluginPermissionDenied) Error() string {
	return "Permission denied while installing plugin " + e.name
}

// IsErrNotImplemented returns true if the error is a NotImplemented error.
// This is returned by the API when a requested feature has not been
// implemented.
//
// Deprecated: use errdefs.IsNotImplemented
func IsErrNotImplemented(err error) bool {
	return errdefs.IsNotImplemented(err)
}

// NewVersionError returns an error if the APIVersion required
// if less than the current supported version
func (cli *Client) NewVersionError(APIrequired, feature string) error {
	if cli.version != "" && versions.LessThan(cli.version, APIrequired) {
		return fmt.Errorf("%q requires API version %s, but the Docker daemon API version is %s", feature, APIrequired, cli.version)
	}
	return nil
}
