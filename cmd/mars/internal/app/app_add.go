package app

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/BurntSushi/toml"
	"github.com/fengleng/mars/log"
	"github.com/spf13/cobra"
	"os"
	"path"
)

var cmdAppAdd = &cobra.Command{
	Use:   "add",
	Short: "add a app service",
	Long:  "add a app service using the repository template. Example: mars app add helloworld",
	Run:   add,
}

func add(cmd *cobra.Command, args []string) {
	a := app
	_, err := toml.DecodeFile("./.env/env.toml", a)
	if err != nil {
		log.Errorf("err: %s", err)
		return
	}

	if len(args) > 0 {
		a.ServiceName = args[0]
	}
	if a.ServiceName == "" {
		prompt := &survey.Input{
			Message: "what is the app service name ?",
			Help:    "what is the app service name ?",
		}
		err := survey.AskOne(prompt, &a.ServiceName)
		if err != nil || app.AppName == "" {
			return
		}
	}
	a.initServiceDir()
	a.initAppInternalConf()
	a.initAppMain()
	a.initAppConf()
	a.initAppConfFile()
	a.initAppWire()
	a.marsLog()
}

func (a *App) initServiceDir() {
	err := os.MkdirAll(path.Join(a.AppDir, a.ServiceName), os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}

	err = os.MkdirAll(path.Join(a.AppDir, a.ServiceName, "cmd"), os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
	err = os.MkdirAll(path.Join(a.AppDir, a.ServiceName, "internal"), os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
	err = os.MkdirAll(path.Join(a.AppDir, a.ServiceName, "internal", "data"), os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
	err = os.MkdirAll(path.Join(a.AppDir, a.ServiceName, "internal", "biz"), os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
	err = os.MkdirAll(path.Join(a.AppDir, a.ServiceName, "internal", "service"), os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
	err = os.MkdirAll(path.Join(a.AppDir, a.ServiceName, "internal", "conf"), os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
}
