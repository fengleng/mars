module github.com/fengleng/mars/contrib/config/apollo

go 1.16

require github.com/apolloconfig/agollo/v4 v4.2.1

require (
	github.com/fengleng/mars v0.0.0-00010101000000-000000000000
	github.com/spf13/viper v1.11.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/fengleng/mars => ../../../
