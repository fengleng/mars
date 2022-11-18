module github.com/fengleng/mars/contrib/config/etcd

go 1.16

require (
	github.com/fengleng/mars v0.0.0-00010101000000-000000000000
	go.etcd.io/etcd/client/v3 v3.5.4
	google.golang.org/grpc v1.46.2
)

replace github.com/fengleng/mars => ../../../
