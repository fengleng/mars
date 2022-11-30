package app_base

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/fengleng/mars/log"
	"github.com/pelletier/go-toml/v2"
	"os"
	"path"
	"path/filepath"
)

type App struct {
	Dir     string `toml:"-"`
	AppName string `json:"app_name" toml:"app_name"`

	ServiceName string `toml:"-"`

	AppDir string `toml:"app_dir" `

	Proto string `json:"Proto" toml:"proto"`

	GitUrl string `json:"git_url" toml:"git_url"`
	GoMod  string `json:"go_mod"`

	Done chan error `toml:"-"`
}

var Instance *App

func (a *App) ProtoCol() string {
	return path.Join(a.Proto, "protocol")
}

func (a *App) ServerDir2(s string) string {
	return path.Join(s, "internal", "server")
}

func (a *App) ServiceDir2(s string) string {
	return path.Join(s, "internal", "service")
}

func (a *App) ConfDir2(s string) string {
	return path.Join(s, "internal", "conf")
}

func (a *App) DataDir2(s string) string {
	return path.Join(s, "internal", "service")
}

func (a *App) ServerDir(s string) string {
	return path.Join(a.AppDir, s, "internal", "server")
}

func (a *App) ServiceDir(s string) string {
	return path.Join(a.AppDir, s, "internal", "service")
}

func (a *App) ConfDir(s string) string {
	return path.Join(a.AppDir, s, "internal", "conf")
}

func (a *App) DataDir(s string) string {
	return path.Join(a.AppDir, s, "internal", "service")
}

func GetApp() *App {
	Instance.ReadFromToml()
	return Instance
}

func (a *App) ProtoClientGo() string {
	return path.Join(a.AppDir, a.ProtoClientGoSuf())
}

func (a *App) ProtoClientGoSuf() string {
	return path.Join("client", "go")
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
			a.Done <- e
			return
		}
		if !override {
			err := errors.New(fmt.Sprintf("Instance Dir is existed:%s", to))
			a.Done <- err
			return
		}
		e = os.RemoveAll(to)
		if e != nil {
			a.Done <- e
			return
		}
	}
	err := os.MkdirAll(to, os.ModePerm)
	if err != nil {
		a.Done <- err
		return
	}
}

func (a *App) WriteToml() {
	ew, err := os.OpenFile(path.Join(a.AppDir, ".env", "env.toml"), os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	err = toml.NewEncoder(ew).Encode(a)
	if err != nil {
		log.Errorf("err: %s", err)
		a.Done <- err
		return
	}
}

func (a *App) ReadFromToml() {
	fmt.Println(a)
	ew, err := os.OpenFile(path.Join(a.AppDir, ".env", "env.toml"), os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}

	err = toml.NewDecoder(ew).Decode(a)
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
}

func (a *App) InitAppDir() {
	a.AppDir = filepath.Join(Instance.Dir, Instance.AppName)
	a.tryNewDir(a.AppDir)
}

func (a *App) WriteMakeFile() {
	file, err := os.OpenFile(path.Join(a.AppDir, "makefile"), os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		a.Done <- err
		return
	}

	_, err = file.WriteString(makeFile)
	if err != nil {
		log.Errorf("err: %s", err)
		a.Done <- err
		return
	}
}
