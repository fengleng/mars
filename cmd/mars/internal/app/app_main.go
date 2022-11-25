package app

import (
	"bytes"
	"github.com/fengleng/mars/log"
	"html/template"
	"os"
	"path"
)

var appMain = `package main

import (
	"github.com/fengleng/mars"
	"github.com/fengleng/mars/log"
	"github.com/fengleng/mars/transport/grpc"
	"github.com/fengleng/mars/transport/http"
	"{{.GoMod}}/{{.ServiceName}}/internal/conf"
	"os"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	id, _ = os.Hostname()
)

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *mars.App {
	return mars.New(
		mars.ID(id),
		mars.Name(Name),
		mars.Version(Version),
		mars.Metadata(map[string]string{}),
		mars.Logger(logger),
		mars.Server(
			gs,
			hs,
		),
	)
}

func main() {
	app, cleanup, err := wireApp(conf.Conf,logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
`

func (a *App) initAppMain() {
	to := path.Join(a.AppDir, a.ServiceName, "cmd", "main.go")
	_, err := os.Stat(to)
	if !os.IsNotExist(err) {
		return
	}
	file, err := os.OpenFile(to, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}

	tmpl, err := template.New("mars_main").Parse(appMain)
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
	_, err = file.Write(bytes)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
}
