package app_base

import (
	"github.com/fengleng/mars/log"
	"os"
	"path"
)

func (a *App) InitServiceDir() {
	err := os.MkdirAll(path.Join(a.AppDir, a.ServiceName), os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}

	err = os.MkdirAll(path.Join(a.AppDir, a.ServiceName, "cmd"), os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
	err = os.MkdirAll(path.Join(a.AppDir, a.ServiceName, "internal"), os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
	err = os.MkdirAll(path.Join(a.AppDir, a.ServiceName, "internal", "data"), os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
	err = os.MkdirAll(path.Join(a.AppDir, a.ServiceName, "internal", "service"), os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
	err = os.MkdirAll(path.Join(a.AppDir, a.ServiceName, "internal", "conf"), os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
	err = os.MkdirAll(path.Join(a.AppDir, a.ServiceName, "internal", "server"), os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
}
