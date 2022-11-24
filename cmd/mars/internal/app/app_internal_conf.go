package app

import (
	"github.com/fengleng/mars/log"
	"os"
	"path"
)

var appInternalConf = `package conf

import "github.com/fengleng/mars/config"

var (
	Conf config.Config

	SvcConf config.Config
)
`

func (a *App) initAppInternalConf() {
	to := path.Join(a.AppDir, a.ServiceName, "internal", "conf", "conf.go")
	_, err := os.Stat(to)
	if !os.IsNotExist(err) {
		return
	}
	file, err := os.OpenFile(to, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
	_, err = file.Write([]byte(appInternalConf))
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
}
