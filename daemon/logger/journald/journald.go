//go:build linux
// +build linux

// Package journald provides the log driver for forwarding server logs
// to endpoints that receive the systemd format.
package journald // import "github.com/harness-community/docker-v23/daemon/logger/journald"

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/coreos/go-systemd/v22/journal"
	"github.com/harness-community/docker-v23/daemon/logger"
	"github.com/harness-community/docker-v23/daemon/logger/loggerutils"
)

const name = "journald"

type journald struct {
	vars map[string]string // additional variables and values to send to the journal along with the log message

	closed chan struct{}
}

func init() {
	if err := logger.RegisterLogDriver(name, New); err != nil {
		panic(err)
	}
	if err := logger.RegisterLogOptValidator(name, validateLogOpt); err != nil {
		panic(err)
	}
}

// sanitizeKeyMode returns the sanitized string so that it could be used in journald.
// In journald log, there are special requirements for fields.
// Fields must be composed of uppercase letters, numbers, and underscores, but must
// not start with an underscore.
func sanitizeKeyMod(s string) string {
	n := ""
	for _, v := range s {
		if 'a' <= v && v <= 'z' {
			v = unicode.ToUpper(v)
		} else if ('Z' < v || v < 'A') && ('9' < v || v < '0') {
			v = '_'
		}
		// If (n == "" && v == '_'), then we will skip as this is the beginning with '_'
		if !(n == "" && v == '_') {
			n += string(v)
		}
	}
	return n
}

// New creates a journald logger using the configuration passed in on
// the context.
func New(info logger.Info) (logger.Logger, error) {
	if !journal.Enabled() {
		return nil, fmt.Errorf("journald is not enabled on this host")
	}

	// parse log tag
	tag, err := loggerutils.ParseLogTag(info, loggerutils.DefaultTemplate)
	if err != nil {
		return nil, err
	}

	vars := map[string]string{
		"CONTAINER_ID":      info.ContainerID[:12],
		"CONTAINER_ID_FULL": info.ContainerID,
		"CONTAINER_NAME":    info.Name(),
		"CONTAINER_TAG":     tag,
		"IMAGE_NAME":        info.ImageName(),
		"SYSLOG_IDENTIFIER": tag,
	}
	extraAttrs, err := info.ExtraAttributes(sanitizeKeyMod)
	if err != nil {
		return nil, err
	}
	for k, v := range extraAttrs {
		vars[k] = v
	}
	return &journald{vars: vars, closed: make(chan struct{})}, nil
}

// We don't actually accept any options, but we have to supply a callback for
// the factory to pass the (probably empty) configuration map to.
func validateLogOpt(cfg map[string]string) error {
	for key := range cfg {
		switch key {
		case "labels":
		case "labels-regex":
		case "env":
		case "env-regex":
		case "tag":
		default:
			return fmt.Errorf("unknown log opt '%s' for journald log driver", key)
		}
	}
	return nil
}

func (s *journald) Log(msg *logger.Message) error {
	vars := map[string]string{}
	for k, v := range s.vars {
		vars[k] = v
	}
	if msg.PLogMetaData != nil {
		vars["CONTAINER_PARTIAL_ID"] = msg.PLogMetaData.ID
		vars["CONTAINER_PARTIAL_ORDINAL"] = strconv.Itoa(msg.PLogMetaData.Ordinal)
		vars["CONTAINER_PARTIAL_LAST"] = strconv.FormatBool(msg.PLogMetaData.Last)
		if !msg.PLogMetaData.Last {
			vars["CONTAINER_PARTIAL_MESSAGE"] = "true"
		}
	}

	line := string(msg.Line)
	source := msg.Source
	logger.PutMessage(msg)

	if source == "stderr" {
		return journal.Send(line, journal.PriErr, vars)
	}
	return journal.Send(line, journal.PriInfo, vars)
}

func (s *journald) Name() string {
	return name
}

func (s *journald) Close() error {
	close(s.closed)
	return nil
}
