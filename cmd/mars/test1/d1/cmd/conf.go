package main

import (
	"flag"
	"github.com/fengleng/mars/config"
	"github.com/fengleng/mars/config/file"
	"github.com/fengleng/mars/log"
	"github.com/fengleng/mars/plugin/config/etcd"
	"github.com/fengleng/test1/d1/internal/conf"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func init() {
	flag.StringVar(&flagconf, "svcConf", "./config.yaml", "config path, eg: -svcConf config.yaml")
	conf.SvcConf = config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	values, err := conf.SvcConf.Value("etcd").Slice()
	if err != nil {
		log.Fatal("err: %s", err)
		panic(err)
	}

	var endPointList []string
	for _, value := range values {
		s, err := value.String()
		if err != nil {
			log.Fatal("err: %s", err)
			panic(err)
		}
		endPointList = append(endPointList, s)
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endPointList,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal("err: %s", err)
		panic(err)
	}
	source, err := etcd.New(cli, etcd.WithPrefix(true), etcd.WithPath("config"))
	if err != nil {
		log.Fatal("err: %s", err)
		panic(err)
	}
	conf.Conf = config.New(
		config.WithSource(source),
	)
}
