package server

import (
	"github.com/fengleng/dddd/t1/internal/conf"
	"github.com/fengleng/dddd/t1/internal/service"
	"github.com/fengleng/mars/log"
	"github.com/soheilhy/cmux"
	"net"

	"github.com/fengleng/dddd/client/go/t1"


)

func init() {
	value := conf.SvcConf.Value("port")
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
	go func() {
		m.Serve()
	}()
}

func init() {
	t1.RegisterT1Server(gRPCServer, service.NewT1Service())
	t1.RegisterT1HTTPServer(httpServer, service.NewT1Service())
}