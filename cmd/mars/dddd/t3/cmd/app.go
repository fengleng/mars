package main

import (
	"github.com/fengleng/dddd/t3/internal/conf"
	"github.com/fengleng/dddd/t3/internal/server"
	"github.com/fengleng/mars"
)

func newApp() (*mars.App, func()) {

	logger, cleanLog := initLog()
	grpcServer := server.NewGRPCServer(conf.SvcConf)
	httpServer := server.NewHTTPServer(conf.SvcConf)

	return mars.New(
			mars.ID(id),
			mars.Name(ServiceName),
			mars.Version(Version),
			mars.Metadata(map[string]string{}),
			mars.Logger(logger),
			mars.Registrar(Register),
			mars.Server(
				grpcServer,
				httpServer,
			),
		), func() {
			cleanLog()
		}
}
