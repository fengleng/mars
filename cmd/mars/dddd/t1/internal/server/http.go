package server

import (
	"github.com/fengleng/dddd/t1/internal/conf"
	"github.com/fengleng/mars/config"
	"github.com/fengleng/mars/middleware/recovery"
	"github.com/fengleng/mars/transport/http"
	"github.com/soheilhy/cmux"
	"net"
)

var (
	httpServer *http.Server
	httpL net.Listener
)

func initHTTPServer(m cmux.CMux) {
	httpL = m.Match(cmux.HTTP1Fast())
	httpServer = newHTTPServer(conf.Conf)
}

// newHTTPServer new a HTTP server.
func newHTTPServer(c config.Config) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
		http.Listener(httpL),
	}
	srv := http.NewServer(opts...)
	return srv
}

func NewHTTPServer(c config.Config) *http.Server {
	return httpServer
}
