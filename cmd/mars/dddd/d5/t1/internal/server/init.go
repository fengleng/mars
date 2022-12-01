package server

import (
	"github.com/fengleng/d5/t1/internal/conf"
	"github.com/fengleng/mars/log"
	"github.com/soheilhy/cmux"
	"net"
)

func init()  {
	value := conf.Conf.Value("port")
	if value==nil {
		panic("invalid port!")
	}
	addr, err := value.String()
	if err !=nil {
		log.Errorf("err: %s",err)
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