package mfile

import (
	"testing"

	"github.com/gososy/sorpc/log"
	"github.com/gososy/sorpc/utils"
)

func Test_packMisc(t *testing.T) {
	x := packMisc(DelayTypeNil, 2, 0)
	log.Infof("%x", x)

	a, b, c := unpackMisc(x)
	log.Infof("%v %v %v", a, b, c)
}

func TestNewFileGroup(t *testing.T) {
	fg := NewFileGroup("test", "/home/pinfire/smq", nil)
	err := fg.Init()
	if err != nil {
		log.Errorf("init err:%v", err)
		return
	}
	err = fg.Write(&Item{
		CreatedAt: utils.Now(),
		Data:      []byte("hello"),
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	log.Infof("ok")
}
