package app

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"os"
)

// CmdApp represents the new command.
var CmdApp = &cobra.Command{
	Use:   "app",
	Short: "Create a app",
	Long:  "Create a app project using the repository template. Example: mars new helloworld",
	//Run:   app.new,
}

func init() {
	CmdApp.AddCommand(cmdAppAdd)
	CmdApp.AddCommand(cmdAppNew)
}

type App struct {
	Dir     string `json:"dir" toml:"dir"`
	AppName string `json:"app_name" toml:"app_name"`
	AppDir  string `toml:"app_dir" `

	//Front   string `json:"Front"`
	//FrontPath   string
	//Backend string `json:"Backend"`
	//BackendPath string
	Proto string `json:"Proto" toml:"proto"`
	//ProtoPath   string

	GitUrl string `json:"git_url" toml:"git_url"`

	done chan error `toml:"done"`
}

var app *App

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	app = &App{
		Dir:   wd,
		Proto: "../proto",
		done:  make(chan error),
	}

	if app.GitUrl = os.Getenv("MARS_GIT"); app.GitUrl == "" {
		app.GitUrl = "https://github.com/fengleng"
	}

	CmdApp.Flags().StringVarP(&app.GitUrl, "git-url", "g", app.GitUrl, "git url")
	CmdApp.Flags().StringVarP(&app.Proto, "proto", "p", app.Proto, "proto Dir")
	//CmdApp.Flags().StringVarP(&app.Backend, "backend", "b", app.Backend, "Backend Dir")
	CmdApp.Flags().StringVar(&app.AppName, "app", "", "app name")
}

func (a *App) tryNewDir(dir string) {
	to := dir
	if _, err := os.Stat(to); !os.IsNotExist(err) {
		fmt.Printf("🚫 %s already exists\n", a.AppName)
		override := false
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("📂 Do you want to override the folder:%s ?", to),
			Help:    fmt.Sprintf("Delete the existing folder %s.", to),
		}
		e := survey.AskOne(prompt, &override)
		if e != nil {
			a.done <- e
			return
		}
		if !override {
			err := errors.New(fmt.Sprintf("app dir is existed:%s", to))
			a.done <- err
			return
		}
		e = os.RemoveAll(to)
		if e != nil {
			a.done <- e
			return
		}
	}
	err := os.MkdirAll(to, os.ModePerm)
	if err != nil {
		a.done <- err
		return
	}
}
