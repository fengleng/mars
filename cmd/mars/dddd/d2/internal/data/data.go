package data

import (
	"github.com/fengleng/mars/log"
)

type Data struct {
}

func NewData() (*Data, func(), error) {
	cleanup := func() {
		log.Info("exit")
	}
	return &Data{}, cleanup, nil
}
