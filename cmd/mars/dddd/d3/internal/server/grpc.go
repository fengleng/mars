package server

import (
	"github.com/fengleng/dddd/d3/internal/conf"
	"github.com/fengleng/mars/config"
	"github.com/fengleng/mars/middleware/recovery"
	"github.com/fengleng/mars/transport/grpc"
	"github.com/soheilhy/cmux"
	"net"
)

var (
	gRPCServer *grpc.Server
	grpcL net.Listener
)

func initGRPCServer(m cmux.CMux) {
	grpcL = m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	gRPCServer = newGRPCServer(conf.Conf)
}

// newGRPCServer new a gRPC server.
func newGRPCServer(c config.Config) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
		grpc.Listener(grpcL),
	}
	srv := grpc.NewServer(opts...)
	return srv
}

func NewGRPCServer(c config.Config) *grpc.Server {
	return gRPCServer
}