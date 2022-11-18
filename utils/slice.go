package utils

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

func IntersectUint64(a []uint64, b []uint64) []uint64 {
	aMap := map[uint64]bool{}
	bMap := map[uint64]bool{}
	var c []uint64

	for _, i := range a {
		aMap[i] = true
	}

	for _, i := range b {
		bMap[i] = true
	}

	for i := range aMap {
		if bMap[i] {
			c = append(c, i)
		}
	}

	return c
}

// DiffSlice
// Deprecated: 这个函数的功能实际上不是 Diff，功能迁移到了 RemoveSlice 中，如果要做 Diff，请使用 RemoveSlice DiffSliceV2
func DiffSlice(a interface{}, b interface{}) interface{} {
	aReflect := reflect.ValueOf(a)
	bReflect := reflect.ValueOf(b)
	aType := reflect.TypeOf(a)
	aMap := map[interface{}]reflect.Value{}
	bMap := map[interface{}]reflect.Value{}

	if aReflect.Kind() != reflect.Slice || aReflect.Kind() != bReflect.Kind() || aReflect.Type() != bReflect.Type() {
		return nil
	}

	cReflect := reflect.MakeSlice(aType, 0, 0)

	for i := int(0); i < aReflect.Len(); i++ {
		aMap[aReflect.Index(i).Interface()] = aReflect.Index(i)
		// log.Errorf("a", aReflect.Index(i))
	}

	for i := int(0); i < bReflect.Len(); i++ {
		bMap[bReflect.Index(i).Interface()] = bReflect.Index(i)
		// log.Errorf("b", bReflect.Index(i))
	}

	for key, reflectValue := range aMap {
		if !bMap[key].IsValid() {
			// log.Errorf("C", reflectValue)
			cReflect = reflect.Append(cReflect, reflectValue)
		}
	}

	return cReflect.Interface()
}

func IntersectSlice(a interface{}, b interface{}) interface{} {
	aReflect := reflect.ValueOf(a)
	bReflect := reflect.ValueOf(b)
	aType := reflect.TypeOf(a)
	aMap := map[interface{}]reflect.Value{}
	bMap := map[interface{}]reflect.Value{}

	if aReflect.Kind() != reflect.Slice || aReflect.Kind() != bReflect.Kind() || aReflect.Type() != bReflect.Type() {
		return nil
	}

	cReflect := reflect.MakeSlice(aType, 0, 0)

	for i := int(0); i < aReflect.Len(); i++ {
		aMap[aReflect.Index(i).Interface()] = aReflect.Index(i)
		// log.Errorf("a", aReflect.Index(i))
	}

	for i := int(0); i < bReflect.Len(); i++ {
		bMap[bReflect.Index(i).Interface()] = bReflect.Index(i)
		// log.Errorf("b", bReflect.Index(i))
	}

	for key, reflectValue := range aMap {
		if bMap[key].IsValid() {
			// log.Errorf("C", reflectValue)
			cReflect = reflect.Append(cReflect, reflectValue)
		}
	}

	return cReflect.Interface()
}

func UniqueSlice(a interface{}) interface{} {
	aReflect := reflect.ValueOf(a)
	aType := reflect.TypeOf(a)
	aMap := map[interface{}]reflect.Value{}

	if aReflect.Kind() != reflect.Slice {
		return nil
	}

	cReflect := reflect.MakeSlice(aType, 0, 0)

	for i := 0; i < aReflect.Len(); i++ {
		if !aMap[aReflect.Index(i).Interface()].IsValid() {
			aMap[aReflect.Index(i).Interface()] = aReflect.Index(i)
			cReflect = reflect.Append(cReflect, aReflect.Index(i))
		}
	}

	return cReflect.Interface()
}

func UniqueSliceV2(s interface{}) interface{} {
	t := reflect.TypeOf(s)
	if t.Kind() != reflect.Slice {
		panic(fmt.Sprintf("s required slice, but got %v", t))
	}

	vo := reflect.ValueOf(s)

	if vo.Len() < 2 {
		return s
	}

	res := reflect.MakeSlice(t, 0, vo.Len())
	m := map[interface{}]struct{}{}
	for i := 0; i < vo.Len(); i++ {
		el := vo.Index(i)
		eli := el.Interface()
		if _, ok := m[eli]; !ok {
			res = reflect.Append(res, el)
			m[eli] = struct{}{}
		}
	}

	return res.Interface()
}

func RandSlice(slice interface{}) {
	rv := reflect.ValueOf(slice)
	if rv.Type().Kind() != reflect.Slice {
		return
	}

	length := rv.Len()
	if length < 2 {
		return
	}

	swap := reflect.Swapper(slice)
	rand.Seed(time.Now().Unix())
	for i := length - 1; i >= 0; i-- {
		j := rand.Intn(length)
		swap(i, j)
	}
	return
}

// DiffSliceV2 传入两个slice
// 如果 a 或者 b 不为 slice 会 panic
// 如果 a 与 b 的元素类型不一致，也会 panic
// 返回的第一个参数为 a 比 b 多的，类型为 a 的类型
// 返回的第二个参数为 b 比 a 多的，类型为 b 的类型
func DiffSliceV2(a interface{}, b interface{}) (interface{}, interface{}) {
	at := reflect.TypeOf(a)
	if at.Kind() != reflect.Slice {
		panic("a is not slice")
	}

	bt := reflect.TypeOf(b)
	if bt.Kind() != reflect.Slice {
		panic("b is not slice")
	}

	atm := at.Elem()
	btm := bt.Elem()

	if atm.Kind() != btm.Kind() {
		panic("a and b are not same type")
	}

	m := map[interface{}]reflect.Value{}

	bv := reflect.ValueOf(b)
	for i := 0; i < bv.Len(); i++ {
		m[bv.Index(i).Interface()] = bv.Index(i)
	}

	c := reflect.MakeSlice(at, 0, 0)
	d := reflect.MakeSlice(bt, 0, 0)
	av := reflect.ValueOf(a)
	for i := 0; i < av.Len(); i++ {
		if !m[av.Index(i).Interface()].IsValid() {
			c = reflect.Append(c, av.Index(i))
		} else {
			delete(m, av.Index(i).Interface())
		}
	}

	for _, value := range m {
		d = reflect.Append(d, value)
	}

	return c.Interface(), d.Interface()
}

// RemoveSlice 传入两个slice
// 如果 src 或者 rm 不为 slice 会 panic
// 如果 src 与 rm 的元素类型不一致，也会 panic
// 返回的第一个参数为 src 中不在 rm 中的元素，数据类型与 src 一致
func RemoveSlice(src interface{}, rm interface{}) interface{} {
	at := reflect.TypeOf(src)
	if at.Kind() != reflect.Slice {
		panic("a is not slice")
	}

	bt := reflect.TypeOf(src)
	if bt.Kind() != reflect.Slice {
		panic("b is not slice")
	}

	atm := at.Elem()
	btm := bt.Elem()

	if atm.Kind() != btm.Kind() {
		panic("a and b are not same type")
	}

	m := map[interface{}]bool{}

	bv := reflect.ValueOf(rm)
	for i := 0; i < bv.Len(); i++ {
		m[bv.Index(i).Interface()] = true
	}

	c := reflect.MakeSlice(at, 0, 0)
	av := reflect.ValueOf(src)
	for i := 0; i < av.Len(); i++ {
		if !m[av.Index(i).Interface()] {
			c = reflect.Append(c, av.Index(i))
			delete(m, av.Index(i).Interface())
		}
	}

	return c.Interface()
}
