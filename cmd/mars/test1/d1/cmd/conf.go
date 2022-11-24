package main

import (
	"flag"
	"github.com/fengleng/mars/config"
	"github.com/fengleng/mars/config/file"
	"github.com/fengleng/mars/contrib/config/etcd"
	"github.com/fengleng/mars/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

var (
	conf config.Config


)

func init() {
	flag.StringVar(&flagconf, "conf", "./config.yaml", "config path, eg: -conf config.yaml")







}

func getLogLevel() log.Level {
	conf = config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	values, err := conf.Value("etcd").Slice()
	if err !=nil {
		log.Fatal("err: %s",err)
		return log.LevelDebug
	}

	var endPointList []string

	for _, value := range values {
		s, err := value.String()
		if err !=nil {
			log.Errorf("err: %s",err)
			return log.LevelDebug
		}
		endPointList = append(endPointList,s)
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endPointList,
		DialTimeout: 5 * time.Second,
	})
	if err !=nil {
		log.Errorf("err: %s",err)
		return log.LevelDebug
	}

	conf = config.New(
		config.WithSource(
			etcd.New()
		),
	)

	return log.LevelDebug
}



