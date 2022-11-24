package mars_log

import (
	"github.com/fengleng/mars/log"
	"github.com/fengleng/mars/middleware/tracing"
	"os"
)

const AppName = ""

func init() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Errorf("err: %s", err)
		return
	}
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"hostname", hostname,
		"app.name", "",
		"app.service", "",
		"service.version", "",
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	log.SetLogger(logger)
}
