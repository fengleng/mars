package utils

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"github.com/gososy/sorpc/log"
	"os"
)

func GetTomlFromFile(fpath string) (map[string]interface{}, error) {
	var config map[string]interface{}
	_, err := toml.DecodeFile(fpath, &config)
	if err != nil {
		log.Errorf("decode toml fail %s, %s", fpath, err)
		return nil, err
	}
	return config, nil
}
func SetTomlToFile(fpath string, config map[string]interface{}) error {
	var configBuffer bytes.Buffer
	e := toml.NewEncoder(&configBuffer)
	err := e.Encode(config)
	if err != nil {
		log.Errorf("toml encode err %+v", err)
	}
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Errorf("can not generate file %s,Error :%v", fpath, err)
		return err
	}
	if _, err := f.Write(configBuffer.Bytes()); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}
