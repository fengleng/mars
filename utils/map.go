package utils

import (
	"reflect"

	"github.com/elliotchance/pie/pie"
)

func MapKeysInt32(m interface{}) []int32 {
	t := reflect.ValueOf(m)
	if t.Kind() != reflect.Map {
		panic("required map type")
	}

	keyType := t.Type().Key()
	if keyType.Kind() != reflect.Uint32 {
		panic("map key type required uint32")
	}

	var result []int32
	for _, v := range t.MapKeys() {
		result = append(result, int32(v.Uint()))
	}

	return result
}

func MapKeysUint32(m interface{}) []uint32 {
	t := reflect.ValueOf(m)
	if t.Kind() != reflect.Map {
		panic("required map type")
	}

	keyType := t.Type().Key()
	if keyType.Kind() != reflect.Uint32 {
		panic("map key type required uint32")
	}

	var result []uint32
	for _, v := range t.MapKeys() {
		result = append(result, uint32(v.Uint()))
	}

	return result
}

func MapKeysUint64(m interface{}) []uint64 {
	t := reflect.ValueOf(m)
	if t.Kind() != reflect.Map {
		panic("required map type")
	}

	keyType := t.Type().Key()
	if keyType.Kind() != reflect.Uint64 {
		panic("map key type required uint64")
	}

	var result []uint64
	for _, v := range t.MapKeys() {
		result = append(result, v.Uint())
	}

	return result
}

func MapKeysString(m interface{}) pie.Strings {
	t := reflect.ValueOf(m)
	if t.Kind() != reflect.Map {
		panic("required map type")
	}

	keyType := t.Type().Key()
	if keyType.Kind() != reflect.String {
		panic("map key type required string")
	}

	var result []string
	for _, v := range t.MapKeys() {
		result = append(result, v.String())
	}

	return result
}

func MapValues(m interface{}) interface{} {
	vo := reflect.ValueOf(m)
	if vo.Kind() != reflect.Map {
		panic("required map type")
	}

	elType := vo.Type().Elem()
	list := reflect.MakeSlice(reflect.SliceOf(elType), 0, 0)

	for _, key := range vo.MapKeys() {
		list = reflect.Append(list, vo.MapIndex(key))
	}

	return list.Interface()
}
