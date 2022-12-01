package main

import (
	"github.com/fengleng/mars/log"
	"os"
)

var (
	// Version go build -ldflags "-X main.Version=x.y.z"
	Version string

	// ServiceName is the name of the compiled software.
	ServiceName = "t2"

	id, _ = os.Hostname()
)

func main() {
	app, cleanup := newApp()

	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
}
