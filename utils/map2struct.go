package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

/*
type User struct {
	Name string
	Age  int8
	Date time.Time
}

func main() {

	data := make(map[string]interface{})
	data["Name"] = "张三"
	data["Age"] = 26
	data["Date"] = "2015-09-29 00:00:00"

	result := &User{}
	err := Map2Struct(data, result)
	fmt.Println(err, fmt.Sprintf("%+v", *result))
	t := time.Now()
	ret := Struct2Map(*result)
	fmt.Println(time.Now().Sub(t), ret)

}
*/

//用map填充结构
func Map2StructByTag(data map[string]interface{}, obj interface{}, tag string) error {
	if len(tag) == 0 {
		return Map2Struct(data, obj)
	}
	for k, v := range data {
		name := NameByTagName(obj, tag, k)
		err := SetField(obj, name, v)
		if err != nil {
			return err
		}
	}
	return nil
}

//用map填充结构
func Map2Struct(data map[string]interface{}, obj interface{}) error {
	for k, v := range data {
		err := SetField(obj, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

//用map的值替换结构的值
func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()        //结构体属性值
	structFieldValue := structValue.FieldByName(name) //结构体单个属性值

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type() //结构体的类型
	val := reflect.ValueOf(value)              //map值的反射值

	var err error
	if structFieldType != val.Type() {
		val, err = TypeConversion(fmt.Sprintf("%v", value), structFieldValue.Type().Name()) //类型转换
		if err != nil {
			return err
		}
	}

	structFieldValue.Set(val)
	return nil
}

//类型转换
func TypeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 8)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 32)
		return reflect.ValueOf(int32(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 32)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "uint32" {
		i, err := strconv.ParseInt(value, 10, 32)
		return reflect.ValueOf(uint32(i)), err
	} else if ntype == "uint64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(uint64(i)), err
	}

	//else if .......增加其他一些类型的转换

	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
}

//NameByTagName fieldbyname 内部也是遍历，这里直接遍历了
func NameByTagName(obj interface{}, tag, name string) string {
	if len(tag) == 0 {
		return name
	}
	t := reflect.TypeOf(obj).Elem()
	fieldnum := t.NumField()
	for i := 0; i < fieldnum; i++ {
		field := t.Field(i)
		tagList := strings.Split(field.Tag.Get(tag), ",")
		for _, v := range tagList {
			if v == name {
				return field.Name
			}
		}
	}
	return name
}
