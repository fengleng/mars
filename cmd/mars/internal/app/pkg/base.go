package pkg

import (
	"bytes"
	"github.com/fengleng/mars/log"
	"go/format"
	"html/template"
	"os"
)

func TemplateParse(to, tpl string, data interface{}) {
	file, err := os.OpenFile(to, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)

	}

	tmpl, err := template.New("template").Parse(tpl)
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, data)
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
	bytes := buf.Bytes()
	bytes, err = format.Source(bytes)
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
	_, err = file.Write(bytes)
	if err != nil {
		log.Errorf("err: %s", err)
		panic(err)
	}
}
