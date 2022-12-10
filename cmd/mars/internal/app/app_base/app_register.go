package app_base

import (
	"bytes"
	"github.com/fengleng/mars/log"
	"go/format"
	"os"
	"path"
	"text/template"
)

var register = `package main

import (
	"{{.GoMod}}/{{.ServiceName}}/internal/conf"
	"{{.GoMod}}/pkg/consts"
	"github.com/fengleng/mars/log"
	"github.com/fengleng/mars/plugin/registry/etcd"
	"github.com/fengleng/mars/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"path"
	"time"
)

var (
	register registry.Registrar
)

func init() {

	endpointList := conf.GetEtcdEndpointList()

	client, err := clientv3.New(clientv3.Config{
		Endpoints:            endpointList,
		DialTimeout:          3 * time.Second,
		DialKeepAliveTimeout: 3 * time.Second,
	})
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
	register = etcd.New(client,etcd.Namespace(path.Join(etcd.DefaultNameSpace,consts.AppName,ServiceName)))
}`

func (a *App) InitAppRegister() {
	to := path.Join(a.AppDir, a.ServiceName, "cmd", "register.go")
	_, err := os.Stat(to)
	if !os.IsNotExist(err) {
		return
	}
	tmpl, err := template.New("register").Parse(register)
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
