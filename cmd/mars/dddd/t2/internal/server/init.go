package server

import (
	"github.com/fengleng/dddd/t2/internal/conf"
	"github.com/fengleng/mars/log"
	"github.com/soheilhy/cmux"
	"net"

	"github.com/fengleng/dddd/client/go/t2"
	"github.com/fengleng/dddd/t2/internal/service"
)

func init() {
	value := conf.Conf.Value("port")
	if value == nil {
		panic("invalid port!")
	}
	addr, err := value.String()
	if err != nil {
		log.Errorf("err: %s", err)
		panic("invalid port!")
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	m := cmux.New(l)
	initGRPCServer(m)
	initHTTPServer(m)
}

func init() {
	t2.RegisterT2Server(gRPCServer, service.NewT2Service())
	t2.RegisterT2HTTPServer(httpServer, service.NewT2Service())
}
