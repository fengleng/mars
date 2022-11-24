package app

import (
	"github.com/fengleng/mars/log"
	"os"
	"path"
)

var appMain = `package main

import (
	"flag"
	"github.com/fengleng/test1/d1/internal/conf"
	"os"

	"github.com/fengleng/mars"
	"github.com/fengleng/mars/log"
	"github.com/fengleng/mars/transport/grpc"
	"github.com/fengleng/mars/transport/http"
	//_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string



	id, _ = os.Hostname()
)

//var app *mars.App

func init() {
	flag.StringVar(&flagconf, "svcConf", "./config.yaml", "config path, eg: -svcConf config.yaml")
}

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

	app, cleanup, err := wireApp(conf.Conf)
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
	_, err = file.Write([]byte(appMain))
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
}
