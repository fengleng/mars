module github.com/fengleng/mars/cmd/mars

go 1.17

replace .github.com/fengleng/mars v0.0.0-20221124152600-0a38125a6e99 => ../../../mars

require (
	github.com/AlecAivazis/survey/v2 v2.3.6
	github.com/emicklei/proto v1.10.0
	github.com/fatih/color v1.13.0
	github.com/fengleng/mars v0.0.0-20221124152600-0a38125a6e99
	github.com/pelletier/go-toml/v2 v2.0.6
	github.com/spf13/cobra v1.4.0
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4
	golang.org/x/text v0.3.5
	google.golang.org/grpc v1.46.2
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/go-playground/form/v4 v4.2.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/sys v0.2.0 // indirect
	golang.org/x/term v0.2.0 // indirect
	google.golang.org/genproto v0.0.0-20220519153652-3a47de7e79bd // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
