package routine

import (
	"context"
	"fmt"
	"strings"

	"github.com/gososy/sorpc/log"
	"github.com/gososy/sorpc/utils"
	"github.com/petermattis/goid"
	uuid "github.com/satori/go.uuid"
)

func genRoutineTaskId() (string, error) {
	u := uuid.NewV4()

	return strings.Replace(u.String(), "-", "", -1)[:12], nil
}
func genNewLogCtx(oldCtx string, id string) string {
	pos := strings.IndexByte(oldCtx, '.')
	var s string
	if pos > 0 {
		s = oldCtx[:pos]
	} else {
		s = oldCtx
	}
	gid := goid.Get()
	return fmt.Sprintf("%s.%s.%d", s, id, gid)
}
func Go(ctx *context.Context, logic func(ctx *context.Context) error) {

	//newCtx := rpc.CloneContext(ctx)
	go func() {
		defer utils.CatchPanic(func(err interface{}) {
			msg := fmt.Sprintf("%s: catch panic in go-routine, err %v", "", err)
			log.Error(msg)
		})
		//routineId := goid.Get()
		err := logic(ctx)

		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
	}()
}
