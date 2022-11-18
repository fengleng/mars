package recovery

import (
	"context"
	"fmt"
	"testing"

	"github.com/gososy/sorpc/errors"
	"github.com/gososy/sorpc/log"
)

func TestOnce(t *testing.T) {
	defer func() {
		if recover() != nil {
			t.Error("fail")
		}
	}()

	next := func(ctx context.Context, req interface{}) (interface{}, error) {
		panic("panic reason")
	}
	_, e := Recovery()(next)(context.Background(), "panic")
	t.Logf("succ and reason is %v", e)
}

func TestNotPanic(t *testing.T) {
	next := func(ctx context.Context, req interface{}) (interface{}, error) {
		return req.(string) + "https://go-kratos.dev", nil
	}

	_, e := Recovery(WithHandler(func(ctx context.Context, req, err interface{}) error {
		return errors.InternalServer("RECOVERY", fmt.Sprintf("panic triggered: %v", err))
	}))(next)(context.Background(), "notPanic")
	if e != nil {
		t.Errorf("e isn't nil")
	}
}

// Deprecated: Remove this test with WithLogger method.
func TestWithLogger(t *testing.T) {
	_ = WithLogger(log.DefaultLogger)
}
