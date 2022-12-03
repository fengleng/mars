package main

import (
	"github.com/fengleng/dddd/t3/internal/conf"
	"github.com/fengleng/mars/log"
	"github.com/fengleng/mars/middleware/tracing"
	"github.com/fengleng/mars/pkg/env"
	"os"
)

func initLog() (log.Logger, func()) {
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.Caller(5),
		"hostname", id,
		"app.name", "test1",
		"app.service", "d2",
		"app.service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	logger = log.NewFilter(logger, log.FilterLevel(getLogLevel()))
	log.SetLogger(logger)
	return logger, func() {}
}

func getLogLevel() log.Level {
	value := conf.SvcConf.Value("env")
	s, err := value.String()
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
	switch s {
	case env.Beta, env.Staging:
		return log.LevelDebug
	default:
		return log.LevelInfo
	}
}
