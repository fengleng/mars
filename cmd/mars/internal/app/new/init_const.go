package new

import (
	"github.com/fengleng/mars/cmd/mars/internal/app/pkg"
	"github.com/fengleng/mars/log"
	"os"
	"path"
)

var constTemplate = `
package consts
const (
	AppName = "{{.AppName}}"
)`

func (a *AppNew) InitConst() {
	err := os.MkdirAll(path.Join(a.AppDir, "pkg", "consts"), os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
	to := path.Join(a.AppDir, a.ServiceName, "pkg", "const", "const.go")
	_, err = os.Stat(to)
	if !os.IsNotExist(err) {
		return
	}

	pkg.TemplateParse(to, constTemplate, a)

}
