package main

import (
	"github.com/fengleng/mars/log"
	"os"
)

var (
	// Version go build -ldflags "-X main.Version=x.y.z"
	Version string

	// Name is the name of the compiled software.
	Name string

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
