package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	//"github.com/elliotchance/pie/pie"
	"github.com/gososy/sorpc/log"
)

func PluckUint64(list interface{}, fieldName string) []uint64 {
	var result []uint64
	vo := reflect.ValueOf(list)
	switch vo.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < vo.Len(); i++ {
			elem := vo.Index(i)
			for elem.Kind() == reflect.Ptr {
				elem = elem.Elem()
			}
			if elem.Kind() != reflect.Struct {
				err := errors.New("element not struct")
				panic(err)
			}
			f := elem.FieldByName(fieldName)
			if !f.IsValid() {
				err := fmt.Errorf("struct missed field %s", fieldName)
				panic(err)
			}
			if f.Kind() != reflect.Uint64 {
				err := fmt.Errorf("struct element %s type required uint64", fieldName)
				panic(err)
			}
			result = append(result, f.Uint())
		}
	default:
		err := errors.New("required list of struct type")
		panic(err)
	}
	return result
}

func PluckInt64(list interface{}, fieldName string) []uint64 {
	var result []uint64
	vo := reflect.ValueOf(list)
	if vo.IsZero() || vo.IsNil() || !vo.IsValid() {
		return result
	}
	switch vo.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < vo.Len(); i++ {
			elem := vo.Index(i)
			for elem.Kind() == reflect.Ptr {
				elem = elem.Elem()
			}
			if elem.Kind() != reflect.Struct {
				err := errors.New("element not struct")
				panic(err)
			}
			f := elem.FieldByName(fieldName)
			if !f.IsValid() {
				err := fmt.Errorf("struct missed field %s", fieldName)
				panic(err)
			}
			if f.Kind() != reflect.Int64 {
				err := fmt.Errorf("struct element %s type required int64", fieldName)
				panic(err)
			}
			result = append(result, f.Uint())
		}
	default:
		err := errors.New("required list of struct type")
		panic(err)
	}
	return result
}

func PluckUint64Map(list interface{}, fieldName string) map[uint64]bool {
	out := PluckUint64(list, fieldName)
	res := map[uint64]bool{}
	for _, v := range out {
		res[v] = true
	}
	return res
}
func PluckUint32(list interface{}, fieldName string) []uint32 {
	var result []uint32
	vo := reflect.ValueOf(list)
	switch vo.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < vo.Len(); i++ {
			elem := vo.Index(i)
			for elem.Kind() == reflect.Ptr {
				elem = elem.Elem()
			}
			if elem.Kind() != reflect.Struct {
				err := errors.New("element not struct")
				panic(err)
			}
			f := elem.FieldByName(fieldName)
			if !f.IsValid() {
				err := fmt.Errorf("struct missed field %s", fieldName)
				panic(err)
			}
			if f.Kind() != reflect.Uint32 {
				err := fmt.Errorf("struct element %s type required uint32", fieldName)
				panic(err)
			}
			result = append(result, uint32(f.Uint()))
		}
	default:
		err := errors.New("required list of struct type")
		panic(err)
	}
	return result
}
func PluckUint32Map(list interface{}, fieldName string) map[uint32]bool {
	out := PluckUint32(list, fieldName)
	res := map[uint32]bool{}
	for _, v := range out {
		res[v] = true
	}
	return res
}
func PluckString(list interface{}, fieldName string) []string {
	var result []string
	vo := reflect.ValueOf(list)
	switch vo.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < vo.Len(); i++ {
			elem := vo.Index(i)
			for elem.Kind() == reflect.Ptr {
				elem = elem.Elem()
			}
			if elem.Kind() != reflect.Struct {
				err := errors.New("element not struct")
				panic(err)
			}
			f := elem.FieldByName(fieldName)
			if !f.IsValid() {
				err := fmt.Errorf("struct missed field %s", fieldName)
				panic(err)
			}
			if f.Kind() != reflect.String {
				err := fmt.Errorf("struct element %s type required string", fieldName)
				panic(err)
			}
			result = append(result, f.String())
		}
	default:
		err := errors.New("required list of struct type")
		panic(err)
	}
	return result
}
func PluckStringMap(list interface{}, fieldName string) map[string]bool {
	out := PluckString(list, fieldName)
	res := map[string]bool{}
	for _, v := range out {
		res[v] = true
	}
	return res
}
func SubSlice(obj interface{}, startIdx, endIdx interface{}) (subSlice interface{}) {
	vo := reflect.ValueOf(obj)
	vk := vo.Kind()
	if vk != reflect.Slice && vk != reflect.Array {
		panic("obj required slice or array type")
	}
	list := reflect.MakeSlice(reflect.SliceOf(vo.Type().Elem()), 0, 0)
	subSlice = list.Interface()
	start, end := Interface2Int(startIdx), Interface2Int(endIdx)
	if start < 0 || end < 0 {
		log.Warnf("invalid slice index start %d end %d", start, end)
		return
	}
	if start >= end {
		if start > end {
			log.Warnf("invalid slice index start %d >= end %d", start, end)
		}
		return
	}
	length := vo.Len()
	for i := start; i < length && i < end; i++ {
		list = reflect.Append(list, vo.Index(i))
	}
	subSlice = list.Interface()
	return
}

// list 是 []StructType
// res 是 *map[fieldType]StructType
func KeyBy(list interface{}, fieldName string, res interface{}) {
	// 取下 field type
	vo := EnsureIsSliceOrArray(list)
	elType := vo.Type().Elem()
	t := elType
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("slice or array element required struct type, but got %v", t))
	}
	var keyType reflect.Type
	if sf, ok := t.FieldByName(fieldName); ok {
		keyType = sf.Type
	} else {
		panic(fmt.Sprintf("not found field %s", fieldName))
	}
	m := reflect.MakeMap(reflect.MapOf(keyType, elType))
	resVo := reflect.ValueOf(res)
	if resVo.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("invalid res type %v, required *map[key]val", resVo.Type()))
	}
	resVo = resVo.Elem()
	EnsureIsMapType(resVo, keyType, elType)
	l := vo.Len()
	for i := 0; i < l; i++ {
		el := vo.Index(i)
		elDef := el
		for elDef.Kind() == reflect.Ptr {
			elDef = elDef.Elem()
		}
		f := elDef.FieldByName(fieldName)
		if !f.IsValid() {
			continue
		}
		m.SetMapIndex(f, el)
	}
	resVo.Set(m)
}

// list 转化 map ，用于判断某个值是否存在 list 中
func SliceToUint64Map(list []uint64) map[uint64]bool {
	result := map[uint64]bool{}
	for _, v := range list {
		result[v] = true
	}
	return result
}

func JoinIntegerSliceToString(list interface{}, sep string) (string, error) {
	listValueReflect := reflect.ValueOf(list)
	switch listValueReflect.Index(0).Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		var intList []int
		for i := 0; i < listValueReflect.Len(); i++ {
			intList = append(intList, Interface2Int(listValueReflect.Index(i).Interface()))
		}

		var intSlice2String string
		for _, item := range intList {
			intSlice2String = intSlice2String + strconv.Itoa(item) + sep
		}
		return strings.TrimRight(intSlice2String, sep), nil
	default:
		return "", errors.New(fmt.Sprintf("Argument 1 expect integer type, given %t", list))
	}
}
