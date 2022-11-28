package app_base

import (
	"bytes"
	"github.com/fengleng/mars/log"
	"go/format"
	"os"
	"path"
	"text/template"
)

var appLog = `package main

import (
	"{{.GoMod}}/{{.ServiceName}}/internal/conf"
	"github.com/fengleng/mars/log"
	"github.com/fengleng/mars/middleware/tracing"
	"github.com/fengleng/mars/pkg/env"
	"os"
)

func initLog() (log.Logger,func()) {
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"hostname", id,
		"Instance.name", "test1",
		"Instance.service", "d2",
		"Instance.service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	logger = log.NewFilter(logger, log.FilterLevel(getLogLevel()))
	log.SetLogger(logger)
	return logger,func() {}
}


func getLogLevel() log.Level {
	value := conf.SvcConf.Value("env")
	s, err := value.String()
	if err !=nil {
		log.Errorf("err: %s",err)
		panic(err)
	}
	switch s {
	case env.Beta,env.Staging:
		return log.LevelDebug
	default:
		return log.LevelInfo
	}
}`

func (a *App) MarsLog() {
	to := path.Join(a.AppDir, a.ServiceName, "cmd", "mars_log.go")
	_, err := os.Stat(to)
	if !os.IsNotExist(err) {
		return
	}
	tmpl, err := template.New("mars-log").Parse(appLog)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, a)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
	bytes := buf.Bytes()
	bytes, err = format.Source(bytes)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}

	file, err := os.OpenFile(to, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
	_, err = file.Write(bytes)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
}
