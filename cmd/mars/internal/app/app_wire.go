package app

import (
	"bytes"
	"github.com/fengleng/mars/log"
	"html/template"
	"os"
	"path"
)

var appWire = `//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/fengleng/mars"
	"github.com/fengleng/mars/config"
	"github.com/fengleng/mars/log"

	myWire "{{.GoMod}}/{{.ServiceName}}/internal/wire"
	"github.com/google/wire"
)

// wireApp init mars application.
func wireApp(conf config.Config,logger log.Logger) (*mars.App, func(), error) {
	panic(wire.Build(myWire.ProviderSet, newApp))
}
`

func (a *App) initAppWire() {
	to := path.Join(a.AppDir, a.ServiceName, "cmd", "wire.go")
	_, err := os.Stat(to)
	if !os.IsNotExist(err) {
		return
	}
	tmpl, err := template.New("mars-log").Parse(appWire)
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
