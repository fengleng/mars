package main

import (
	"github.com/fengleng/mars/log"
	"os"
)

var (
	// Version go build -ldflags "-X main.Version=x.y.z"
	Version string

	// ServiceName is the name of the compiled software.
	ServiceName = "t1"

	id, _ = os.Hostname()
)

func main() {
	Instance, cleanup := newApp()

	defer cleanup()

	// start and wait for stop signal
	if err := Instance.Run(); err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
}
