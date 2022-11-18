package utils

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gososy/sorpc/log"
	"github.com/tidwall/pretty"
	"runtime/debug"
	"strings"
)

var WarningFunc func(string)

func CheckError(err error) {
	if err != nil {
		log.Fatal("catch error:", err)
	}
}
func pb2Json(pb proto.Message) string {
	dat, err := json.Marshal(pb)
	if err != nil {
		return fmt.Sprintf("<err:%s>", err)
	}
	return string(dat)
}
func PrintPb(pb proto.Message) {
	j := pb2Json(pb)
	x := pretty.Pretty([]byte(j))
	log.Info(string(x))
}
func PrettyPrint(i interface{}) {
	buf, err := json.Marshal(i)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	x := pretty.Pretty(buf)
	log.Info(string(x))
}
func PrintPbPrefix(pb proto.Message, args ...interface{}) {
	prefix := fmt.Sprint(args...)
	j := pb2Json(pb)
	log.Info(prefix, ":", j)
}
func PrintPbPrefixFormat(pb proto.Message, format string, args ...interface{}) {
	prefix := fmt.Sprintf(format, args...)
	j := pb2Json(pb)
	log.Info(prefix, ":", j)
}
func PrintStack() {
	st := debug.Stack()
	if len(st) > 0 {
		log.Info("dump stack:")
		lines := strings.Split(string(st), "\n")
		for _, line := range lines {
			log.Info("  ", line)
		}
	} else {
		log.Info("stack is empty")
	}
}
func CatchPanic(panicCallback func(err interface{})) {
	if err := recover(); err != nil {
		log.Errorf("PROCESS PANIC: err %s", err)
		st := debug.Stack()
		if len(st) > 0 {
			log.Errorf("dump stack (%s):", err)
			lines := strings.Split(string(st), "\n")
			for _, line := range lines {
				log.Error("  ", line)
			}
		} else {
			log.Errorf("stack is empty (%s)", err)
		}
		if panicCallback != nil {
			panicCallback(err)
		}
	}
}
