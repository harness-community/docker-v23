package daemon // import "github.com/DevanshMathur19/docker-v23/daemon"

import (
	// Importing packages here only to make sure their init gets called and
	// therefore they register themselves to the logdriver factory.
	_ "github.com/DevanshMathur19/docker-v23/daemon/logger/awslogs"
	_ "github.com/DevanshMathur19/docker-v23/daemon/logger/fluentd"
	_ "github.com/DevanshMathur19/docker-v23/daemon/logger/gcplogs"
	_ "github.com/DevanshMathur19/docker-v23/daemon/logger/gelf"
	_ "github.com/DevanshMathur19/docker-v23/daemon/logger/journald"
	_ "github.com/DevanshMathur19/docker-v23/daemon/logger/jsonfilelog"
	_ "github.com/DevanshMathur19/docker-v23/daemon/logger/local"
	_ "github.com/DevanshMathur19/docker-v23/daemon/logger/logentries"
	_ "github.com/DevanshMathur19/docker-v23/daemon/logger/loggerutils/cache"
	_ "github.com/DevanshMathur19/docker-v23/daemon/logger/splunk"
	_ "github.com/DevanshMathur19/docker-v23/daemon/logger/syslog"
)
