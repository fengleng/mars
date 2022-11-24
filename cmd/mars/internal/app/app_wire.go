package app

import (
	"github.com/fengleng/mars/log"
	"os"
	"path"
)

var appWire = `//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/fengleng/mars"
	"github.com/fengleng/mars-layout/internal/biz"
	"github.com/fengleng/mars-layout/internal/data"
	"github.com/fengleng/mars-layout/internal/server"
	"github.com/fengleng/mars-layout/internal/service"
	"github.com/fengleng/mars/config"

	"github.com/google/wire"
)

// wireApp init mars application.
func wireApp(conf config.Config) (*mars.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
`

func (a *App) initAppWire() {
	to := path.Join(a.AppDir, a.ServiceName, "cmd", "wire.go")
	_, err := os.Stat(to)
	if !os.IsNotExist(err) {
		return
	}
	file, err := os.OpenFile(to, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
	_, err = file.Write([]byte(appWire))
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
}
