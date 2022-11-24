module github.com/fengleng/mars/contrib/log/aliyun

go 1.16

require (
	github.com/aliyun/aliyun-log-go-sdk v0.1.37

	google.golang.org/protobuf v1.28.1
)

replace (
	github.com/fengleng/mars => ../../../
	github.com/gogo/protobuf v1.3.1 => github.com/gogo/protobuf v1.3.2
)
