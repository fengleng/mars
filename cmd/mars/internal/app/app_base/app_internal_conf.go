package app_base

import (
	"bytes"
	"github.com/fengleng/mars/log"
	"go/format"
	"html/template"
	"os"
	"path"
)

var appInternalConf = `package conf

import (
	"flag"
	"github.com/fengleng/mars/config"
	"github.com/fengleng/mars/config/file"
	"github.com/fengleng/mars/log"
	"github.com/fengleng/mars/plugin/config/etcd"
	clientV3 "go.etcd.io/etcd/client/v3"
	"time"
)

var (
	Conf config.Config

	SvcConf config.Config

	flagConf string
	
	configPath string = "config/{{.ServiceName}}.json"
)

func init() {
	flag.StringVar(&flagConf, "svcConf", "./config.yaml", "config path, eg: -svcConf config.yaml")
}

func init() {
	SvcConf = config.New(
		config.WithSource(
			file.NewSource(flagConf),
		),
	)

	err := SvcConf.Load()
	if err !=nil {
		log.Errorf("err: %s",err)
		panic(err)
	}

	values, err := SvcConf.Value("etcd").Slice()
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}

	var endPointList []string
	for _, value := range values {
		s, err := value.String()
		if err != nil {
			log.Errorf("err: %s", err)
			panic(err)
		}
		endPointList = append(endPointList, s)
	}

	client, err := clientV3.New(clientV3.Config{
		Endpoints:            endPointList,
		DialTimeout:          3 * time.Second,
		DialKeepAliveTimeout: 3 * time.Second,
	})
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}

	source, err := etcd.New(client,etcd.WithPath(configPath),etcd.WithPrefix(true))
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
	Conf = config.New(config.WithSource(source))
	err = Conf.Load()
	if err !=nil {
		log.Errorf("err: %s",err)
		panic(err)
	}
}
`

func (a *App) InitAppInternalConf() {
	to := path.Join(a.AppDir, a.ServiceName, "internal", "conf", "conf.go")
	_, err := os.Stat(to)
	if !os.IsNotExist(err) {
		return
	}
	file, err := os.OpenFile(to, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)

	}

	tmpl, err := template.New("mars_internal_conf").Parse(appInternalConf)
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
