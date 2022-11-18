package utils

import (
	"encoding/json"

	"github.com/gososy/sorpc/log"
	yaml "gopkg.in/yaml.v3"
)

func yamlConvert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = yamlConvert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = yamlConvert(v)
		}
	}
	return i
}
func YamlDecode(code string) (interface{}, error) {
	var body interface{}
	if err := yaml.Unmarshal([]byte(code), &body); err != nil {
		return nil, err
	}
	body = yamlConvert(body)
	return body, nil
}
func Yaml2Json(code string) (string, error) {
	body, err := YamlDecode(code)
	if err != nil {
		log.Errorf("YamlDecode err:%v: %s", err, code)
		return "", err
	}
	b, err := json.Marshal(body)
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}
	return string(b), nil
}
