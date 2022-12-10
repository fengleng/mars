package app

import (
	"github.com/fengleng/mars/cmd/mars/internal/app/app_base"
	"github.com/fengleng/mars/cmd/mars/internal/app/new"
	"github.com/fengleng/mars/cmd/mars/internal/app/service/GenAll"
	"github.com/fengleng/mars/cmd/mars/internal/app/service/add"
	"github.com/fengleng/mars/cmd/mars/internal/app/service/client"
	"github.com/fengleng/mars/cmd/mars/internal/app/service/server"
	"github.com/spf13/cobra"
	"os"
)

// CmdApp represents the new command.
var CmdApp = &cobra.Command{
	Use:   "app",
	Short: "Create a app",
	Long:  "Create a app project using the repository template. Example: mars new helloworld",
	//Run:   Instance.new,
}

func init() {
	CmdApp.AddCommand(add.CmdAppServiceAdd)
	CmdApp.AddCommand(new.CmdAppNew)
	CmdApp.AddCommand(client.CmdClient)
	CmdApp.AddCommand(server.CmdServer)
	CmdApp.AddCommand(GenAll.CmdGenAll)
}

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	app_base.Instance = &app_base.App{
		Dir:   wd,
		Proto: "./proto",
		Done:  make(chan error),
	}

	if app_base.Instance.GitUrl = os.Getenv("MARS_GIT"); app_base.Instance.GitUrl == "" {
		app_base.Instance.GitUrl = "https://github.com/fengleng"
	}

	CmdApp.Flags().StringVarP(&app_base.Instance.GitUrl, "git-url", "g", app_base.Instance.GitUrl, "git url")
	CmdApp.Flags().StringVarP(&app_base.Instance.Proto, "proto", "p", app_base.Instance.Proto, "proto dir")
	//CmdApp.Flags().StringVarP(&Instance.Backend, "backend", "b", Instance.Backend, "Backend dir")
	CmdApp.Flags().StringVar(&app_base.Instance.AppName, "app", "", "app name")
}
