package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/gososy/sorpc/log"
)

func IsBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

func ReverseAnySlice(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

func ToPbList(list interface{}) []proto.Message {
	v := reflect.ValueOf(list)
	if v.Kind() != reflect.Slice {
		return nil
	}
	var pbList []proto.Message
	for i := 0; i < v.Len(); i++ {
		el := v.Index(i).Interface()
		if pb, ok := el.(proto.Message); ok {
			pbList = append(pbList, pb)
		}
	}
	return pbList
}

func ToPb(i interface{}) proto.Message {
	log.Infof("%+v", i)
	v := reflect.ValueOf(i)
	el := v.Interface()
	if pb, ok := el.(proto.Message); ok {
		return pb
	}
	return nil
}

func InSliceUint32(i uint32, list []uint32) bool {
	for _, v := range list {
		if i == v {
			return true
		}
	}
	return false
}

func InSliceUint64(i uint64, list []uint64) bool {
	for _, v := range list {
		if i == v {
			return true
		}
	}
	return false
}

func InSliceStr(s string, list []string) bool {
	for _, v := range list {
		if s == v {
			return true
		}
	}
	return false
}

func GetSliceIndexStr(s string, list []string) int {
	for k, v := range list {
		if v == s {
			return k
		}
	}
	return -1
}

func Min(list ...interface{}) interface{} {
	length := len(list)
	if length >= 1 {
		min := list[0]
		minV := reflect.ValueOf(min)
		for i := 1; i < length; i++ {
			sourceV := reflect.ValueOf(list[i])
			if minV.Kind() != sourceV.Kind() {
				return nil
			}
			switch minV.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if sourceV.Int() < minV.Int() {
					min = list[i]
					minV = sourceV
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				if sourceV.Uint() < minV.Uint() {
					min = list[i]
					minV = sourceV
				}
			case reflect.Float32, reflect.Float64:
				if sourceV.Float() < minV.Float() {
					min = list[i]
					minV = sourceV
				}
			default:
				return nil
			}
		}
		return min
	}
	return nil
}

func JoinUint32(a []uint32, sep string) string {
	switch len(a) {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%d", a[0])
	}
	var b []string
	for _, i := range a {
		b = append(b, fmt.Sprintf("%d", i))
	}
	return strings.Join(b, sep)
}

func JoinUint64(a []uint64, sep string) string {
	switch len(a) {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%d", a[0])
	}
	var b []string
	for _, i := range a {
		b = append(b, fmt.Sprintf("%d", i))
	}
	return strings.Join(b, sep)
}
