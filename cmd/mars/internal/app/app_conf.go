package app

import (
	"bytes"
	"github.com/fengleng/mars/log"
	"html/template"
	"os"
	"path"
)

var appConf = `package main

import (
	"github.com/fengleng/mars/log"
	"github.com/fengleng/mars/middleware/tracing"
	"os"
)

var (
	Version string
)

func init() {
	hostname, err := os.Hostname()
	if err !=nil {
		log.Errorf("err: %s",err)
		return 
	}
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"hostname", hostname,
		"app.name", "test1",
		"app.service", "d2",
		"app.service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	log.SetLogger(logger)
}
`

func (a *App) initAppConf() {
	to := path.Join(a.AppDir, a.ServiceName, "cmd", "conf.go")
	_, err := os.Stat(to)
	if !os.IsNotExist(err) {
		return
	}
	tmpl, err := template.New("mars-conf").Parse(appConf)
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
