package atexit

import (
	"testing"
	"time"
)

func TestAtExit(t *testing.T) {

	Register(func() {
		t.Log("hahhahah")
		time.Sleep(time.Second * 2)
	})

	for {
		t.Log("lllllllll")
		time.Sleep(time.Second * 10)
	}

}
