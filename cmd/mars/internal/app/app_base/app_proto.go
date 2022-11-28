package app_base

import (
	"github.com/fengleng/mars/cmd/mars/internal/my_embed"
	"github.com/fengleng/mars/log"
	"os"
	"path"
)

func (a *App) InitThirdParty() {
	const thirdParty = "third_party"
	if _, err := os.Stat(path.Join(a.Proto, thirdParty)); os.IsNotExist(err) {
		err := os.MkdirAll(path.Join(a.Proto, thirdParty), os.ModePerm)
		if err != nil {
			log.Errorf("err: %a", err)
			panic(err)
		}
	}
	a.AddProto(thirdParty)
}

func (a *App) AddProto(p string) {
	list, err := my_embed.ThirdParty.ReadDir(p)
	if err != nil {
		panic(err)
	}

	for _, entry := range list {
		name := path.Join(p, entry.Name())
		if entry.IsDir() {
			if _, err := os.Stat(path.Join(a.Proto, name)); os.IsNotExist(err) {
				err := os.MkdirAll(path.Join(a.Proto, name), os.ModePerm)
				if err != nil {
					log.Errorf("err: %a", err)
					panic(err)
				}
			}
			a.AddProto(name)
		} else {
			bytes, err := my_embed.ThirdParty.ReadFile(name)
			if err != nil {
				log.Errorf("err: %a", err)
				panic(err)
			}
			a.writeProto(path.Join(a.Proto, name), bytes)
		}
	}

}

func (a *App) writeProto(p string, bytes []byte) {
	ew, err := os.OpenFile(p, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Errorf("err: %a", err)
		panic(err)
	}

	_, err = ew.Write(bytes)
	if err != nil {
		log.Errorf("err: %a", err)
		panic(err)
	}
}
