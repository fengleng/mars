package conf

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

	configPath string = "config/t1.json"
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
	if err != nil {
		log.Errorf("err: %s", err)
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

	source, err := etcd.New(client, etcd.WithPath(configPath))
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
	Conf = config.New(config.WithSource(source))
	err = Conf.Load()
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
}
