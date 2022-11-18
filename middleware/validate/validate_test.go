package validate

import (
	"context"
	"errors"
	"testing"

	kratoserrors "github.com/gososy/sorpc/errors"
	"github.com/gososy/sorpc/middleware"
)

// protoVali implement validate.validator
type protoVali struct {
	name  string
	age   int
	isErr bool
}

func (v protoVali) Validate() error {
	if v.name == "" || v.age < 0 {
		return errors.New("err")
	}
	return nil
}

func TestTable(t *testing.T) {
	var mock middleware.Handler = func(ctx context.Context, req interface{}) (interface{}, error) { return nil, nil }

	tests := []protoVali{
		{"v1", 365, false},
		{"v2", -1, true},
		{"", 365, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := Validator()(mock)
			_, err := v(context.Background(), test)
			if want, have := test.isErr, kratoserrors.IsBadRequest(err); want != have {
				t.Errorf("fail data %v, want %v, have %v", test, want, have)
			}
		})
	}
}
