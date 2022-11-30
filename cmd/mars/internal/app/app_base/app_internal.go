package app_base

import (
	"bytes"
	"github.com/fengleng/mars/log"
	"go/format"
	"html/template"
	"os"
	"path"
)

//var biz = `package biz
//import "{{.GoMod}}/{{.ServiceName}}/internal"
//// ProviderSet is biz providers.
//func init() {
//	internal.RegisterProvider()
//}`
//
//var data = `package data
//import "{{.GoMod}}/{{.ServiceName}}/internal"
//// ProviderSet is data providers.
//func init() {
//	internal.RegisterProvider()
//}`
//
//var server = `package server
//import "{{.GoMod}}/{{.ServiceName}}/internal"
//// ProviderSet is server providers.
//func init() {
//	internal.RegisterProvider(NewGRPCServer,NewHTTPServer)
//}`
//
//var service = `package service
//import "{{.GoMod}}/{{.ServiceName}}/internal"
//// ProviderSet is service providers.
//func init() {
//	internal.RegisterProvider()
//}`
//
//var myWire = `package wire
//import (
//	"{{.GoMod}}/{{.ServiceName}}/internal/server"
//	"github.com/google/wire"
//)
//var ProviderSet = wire.NewSet(server.NewGRPCServer,server.NewHTTPServer)`

var serverInit = `package server

import (
	"{{.GoMod}}/{{.ServiceName}}/internal/conf"
	"github.com/fengleng/mars/log"
	"github.com/soheilhy/cmux"
	"net"
)

func init()  {
	value := conf.SvcConf.Value("port")
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
}`

var httpServer = `package server

import (
	"{{.GoMod}}/{{.ServiceName}}/internal/conf"
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
`

var grpcServer = `package server

import (
	"{{.GoMod}}/{{.ServiceName}}/internal/conf"
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
}`

var internalInitGo = `package internal
import "github.com/google/wire"
var ProviderSet wire.ProviderSet
var providerList []interface{}
func RegisterProvider(providers ...interface{})  {
	providerList = append(providerList,providers...)
}
func init() {
	ProviderSet = wire.NewSet(providerList...)
}
`

func (a *App) InitInternal() {
	//a.wireInit("biz","wire.go",biz)
	//a.wireInit("data","wire.go",data)
	//a.wireInit("service","wire.go",service)
	//a.wireInit("wire", "wire.go", myWire)

	//a.internalWireInit()

	a.serverInit("init.go", serverInit)
	a.serverInit("http.go", httpServer)
	a.serverInit("grpc.go", grpcServer)
}

func (a *App) internalWireInit() {
	to := path.Join(a.AppDir, a.ServiceName, "internal", "init.go")
	_, err := os.Stat(to)
	if !os.IsNotExist(err) {
		return
	}
	tmpl, err := template.New("server").Parse(internalInitGo)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, a)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
	bytes := buf.Bytes()

	bytes, err = format.Source(bytes)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}

	file, err := os.OpenFile(to, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
	_, err = file.Write(bytes)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
}

func (a *App) wireInit(d, f, v string) {
	to := path.Join(a.AppDir, a.ServiceName, "internal", d, f)
	_, err := os.Stat(to)
	if !os.IsNotExist(err) {
		return
	}
	tmpl, err := template.New("internal").Parse(v)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, a)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
	bytes := buf.Bytes()

	file, err := os.OpenFile(to, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
	_, err = file.Write(bytes)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
}

func (a *App) serverInit(f, s string) {
	to := path.Join(a.AppDir, a.ServiceName, "internal", "server", f)
	_, err := os.Stat(to)
	if !os.IsNotExist(err) {
		return
	}
	tmpl, err := template.New("server").Parse(s)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, a)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
	bytes := buf.Bytes()

	file, err := os.OpenFile(to, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
	_, err = file.Write(bytes)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
}
