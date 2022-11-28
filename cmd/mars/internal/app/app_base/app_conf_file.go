package app_base

import (
	"github.com/fengleng/mars/log"
	"os"
	"path"
)

var appConfFile = `port: 3333
etcd:
  - 127.0.0.1:2379
`

func (a *App) InitAppConfFile() {
	to := path.Join(a.AppDir, a.ServiceName, "cmd", "config.yaml")
	_, err := os.Stat(to)
	if !os.IsNotExist(err) {
		return
	}
	file, err := os.OpenFile(to, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
	_, err = file.Write([]byte(appConfFile))
	if err != nil {
		log.Errorf("err: %s", err)
		os.Exit(1)
	}
}
