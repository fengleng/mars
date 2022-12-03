package main

import (
	"github.com/fengleng/dddd/t3/internal/conf"
	"github.com/fengleng/mars/log"
	"github.com/fengleng/mars/plugin/registry/etcd"
	"github.com/fengleng/mars/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

var (
	Register registry.Registrar
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
	Register = etcd.New(client)
}
