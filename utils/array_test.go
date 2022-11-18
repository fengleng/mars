package utils

import (
	"testing"
)

func TestJoinIntegerSliceToString(t *testing.T) {
	uint64List := []uint64{1, 2}
	if str, err := JoinIntegerSliceToString(uint64List, ","); err != nil {
		t.Error(err)
	} else if str != "1,2" {
		t.Error("str:" + str)
	}
}
