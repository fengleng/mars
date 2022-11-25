package app

import (
	"bytes"
	"github.com/fengleng/mars/log"
	"html/template"
	"os"
	"path"
)

var appConf = `package main

import (
	"flag"
	"{{.GoMod}}/{{.ServiceName}}/internal/conf"
	"{{.GoMod}}/{{.ServiceName}}/internal/server"
	"github.com/fengleng/mars/config"
	"github.com/fengleng/mars/config/file"
	"github.com/fengleng/mars/log"
	"github.com/fengleng/mars/pkg/env"
	"github.com/fengleng/mars/plugin/config/etcd"
	clientV3 "go.etcd.io/etcd/client/v3"
	"time"
)

var (
	Version string
	flagConf string
)


func init() {
	flag.StringVar(&flagConf, "svcConf", "./config.yaml", "config path, eg: -svcConf config.yaml")
}

func init() {
	initConf()

	server.Init()
}


func initConf()  {
	conf.SvcConf = config.New(
		config.WithSource(
			file.NewSource(flagConf),
		),
	)

	values, err := conf.SvcConf.Value("etcd").Slice()
	if err !=nil {
		log.Errorf("err: %s",err)
		panic(err)
	}

	var endPointList []string
	for _, value := range values {
		s, err := value.String()
		if err !=nil {
			log.Errorf("err: %s",err)
			panic(err)
		}
		endPointList = append(endPointList,s)
	}

	client, err := clientV3.New(clientV3.Config{
		Endpoints:            endPointList,
		DialTimeout:          3 * time.Second,
		DialKeepAliveTimeout: 3 * time.Second,
	})
	if err !=nil {
		log.Errorf("err: %s",err)
		panic(err)
	}

	source, err := etcd.New(client)
	if err !=nil {
		log.Errorf("err: %s",err)
		panic(err)
	}
	conf.Conf = config.New(config.WithSource(source))
}

func getLogLevel() log.Level {
	value := conf.SvcConf.Value("env")
	s, err := value.String()
	if err !=nil {
		log.Errorf("err: %s",err)
		panic(err)
	}
	switch s {
	case env.Beta,env.Staging:
		return log.LevelDebug
	default:
		return log.LevelInfo
	}
}
`

func (a *App) initAppConf() {
	to := path.Join(a.AppDir, a.ServiceName, "cmd", "conf.go")
	_, err := os.Stat(to)
	if !os.IsNotExist(err) {
		return
	}
	tmpl, err := template.New("mars-conf").Parse(appConf)
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
