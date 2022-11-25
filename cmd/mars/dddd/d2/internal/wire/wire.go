package wire

import (
	"github.com/fengleng/dddd/d2/internal/server"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(server.NewGRPCServer, server.NewHTTPServer)
