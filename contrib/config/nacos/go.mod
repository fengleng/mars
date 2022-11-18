module github.com/fengleng/mars/contrib/config/nacos

go 1.16

require (
	github.com/fengleng/mars v0.0.0-00010101000000-000000000000
	github.com/nacos-group/nacos-sdk-go v1.0.9
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/fengleng/mars => ../../../

replace github.com/buger/jsonparser => github.com/buger/jsonparser v1.1.1
