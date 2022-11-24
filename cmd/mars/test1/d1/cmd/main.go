package main

import (
	"flag"
	"os"

	"github.com/fengleng/mars"
	"github.com/fengleng/mars/config"
	"github.com/fengleng/mars/config/file"
	"github.com/fengleng/mars/log"
	"github.com/fengleng/mars/middleware/tracing"
	"github.com/fengleng/mars/transport/grpc"
	"github.com/fengleng/mars/transport/http"
	//_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string

	// flagconf is the config flag.
	flagconf string

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
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc svcConf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
