package app_base

import (
	"bytes"
	"github.com/fengleng/mars/log"
	"go/format"
	"html/template"
	"os"
	"path"
)

var appMain = `package main

import (
	"github.com/fengleng/mars/log"
	"os"
)

var (
	// Version go build -ldflags "-X main.Version=x.y.z"
	Version string
	
	// ServiceName is the name of the compiled software.
	ServiceName = "{{.ServiceName}}"

	id, _ = os.Hostname()
)

func main() {
	app, cleanup := newApp()

	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		log.Errorf("err: %s",err)
		panic(err)
	}
}`

func (a *App) InitAppMain() {
	to := path.Join(a.AppDir, a.ServiceName, "cmd", "main.go")
	_, err := os.Stat(to)
	if !os.IsNotExist(err) {
		return
	}
	file, err := os.OpenFile(to, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)

	}

	tmpl, err := template.New("mars_main").Parse(appMain)
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, a)
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
	bytes := buf.Bytes()
	bytes, err = format.Source(bytes)
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
	_, err = file.Write(bytes)
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
}
