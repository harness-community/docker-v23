package daemon // import "github.com/harness-community/docker-v23/daemon"

import (
	// Importing packages here only to make sure their init gets called and
	// therefore they register themselves to the logdriver factory.
	_ "github.com/harness-community/docker-v23/daemon/logger/awslogs"
	_ "github.com/harness-community/docker-v23/daemon/logger/etwlogs"
	_ "github.com/harness-community/docker-v23/daemon/logger/fluentd"
	_ "github.com/harness-community/docker-v23/daemon/logger/gcplogs"
	_ "github.com/harness-community/docker-v23/daemon/logger/gelf"
	_ "github.com/harness-community/docker-v23/daemon/logger/jsonfilelog"
	_ "github.com/harness-community/docker-v23/daemon/logger/logentries"
	_ "github.com/harness-community/docker-v23/daemon/logger/loggerutils/cache"
	_ "github.com/harness-community/docker-v23/daemon/logger/splunk"
	_ "github.com/harness-community/docker-v23/daemon/logger/syslog"
)
