//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/fengleng/mars-layout/internal/biz"
	"github.com/fengleng/mars-layout/internal/data"
	"github.com/fengleng/mars-layout/internal/server"
	"github.com/fengleng/mars-layout/internal/service"

	"github.com/fengleng/mars/log"
	"github.com/google/wire"
)

// wireApp init mars application.
func wireApp(*svcConf.Server, *svcConf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
