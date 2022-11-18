package kfklog

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func NewLogger() *log.Logger {
	if !FileExists(logPath) {
		err := CreateFile(logPath)
		if err != nil {
			fmt.Printf("err:%v", err)
			return log.New(ioutil.Discard, "[Kafuka] ", log.Lshortfile)
		}
	}
	f, err := os.Open(logPath)
	if err != nil {
		fmt.Printf("err:%v", err)
		return log.New(ioutil.Discard, "[Kafuka] ", log.Lshortfile)
	}

	return log.New(io.MultiWriter(os.Stderr, f), "[Kafuka] ", log.Lshortfile)
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func CreateFile(name string) error {
	fo, err := os.Create(name)
	if err != nil {
		return err
	}
	defer func() {
		fo.Close()
	}()
	return nil
}
