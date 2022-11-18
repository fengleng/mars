package utils

import (
	"database/sql"
	"errors"
	"reflect"
	"sync"
)

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func OrmStruct2Map(s interface{}, skip ...string) map[string]interface{} {
	m := make(map[string]interface{})
	elem := reflect.ValueOf(s).Elem()

	relType := elem.Type()

	for i := 0; i < relType.NumField(); i++ {
		n := relType.Field(i).Name
		if len(n) >= 3 &&
			n[0] == 'X' && n[1] == 'X' && n[2] == 'X' {
			continue
		}
		if n == "DeletedAt" || n == "CreatedAt" || n == "UpdatedAt" || n == "Id" {
			// skip
		} else {
			if len(skip) > 0 {
				ignore := false
				for _, v := range skip {
					if v == n {
						ignore = true
						break
					}
				}
				if ignore {
					continue
				}
			}
			key := Camel2UnderScore(n)
			m[key] = elem.Field(i).Interface()
		}
	}

	return m
}

func OrmStruct2Map4Update(s interface{}, skip ...string) map[string]interface{} {
	m := OrmStruct2Map(s, skip...)

	n := map[string]interface{}{}
	for k, v := range m {
		if !IsBlank(reflect.ValueOf(v)) {
			n[k] = v
		}
	}

	return n
}

var camel2UnderScoreMap = map[string]string{}
var camel2UnderScoreMapMu sync.RWMutex

func OrmStruct2MapLower(s interface{}, skip ...string) map[string]interface{} {
	m := make(map[string]interface{})
	elem := reflect.ValueOf(s).Elem()

	relType := elem.Type()

	var nameUpdate map[string]string

	camel2UnderScoreMapMu.RLock()

	for i := 0; i < relType.NumField(); i++ {
		camelName := relType.Field(i).Name
		if len(camelName) >= 3 &&
			camelName[0] == 'X' && camelName[1] == 'X' && camelName[2] == 'X' {
			continue
		}
		n := camel2UnderScoreMap[camelName]
		if n == "" {
			n = Camel2UnderScore(camelName)
			if nameUpdate == nil {
				nameUpdate = map[string]string{}
			}
			nameUpdate[camelName] = n
		}

		if n == "deleted_at" || n == "created_at" || n == "updated_at" || n == "id" {
			// skip
		} else {
			if len(skip) > 0 {
				ignore := false
				for _, v := range skip {
					if v == n {
						ignore = true
						break
					}
				}
				if ignore {
					continue
				}
			}
			m[n] = elem.Field(i).Interface()
		}
	}

	camel2UnderScoreMapMu.RUnlock()

	if nameUpdate != nil {
		camel2UnderScoreMapMu.Lock()
		for k, v := range nameUpdate {
			camel2UnderScoreMap[k] = v
		}
		camel2UnderScoreMapMu.Unlock()
	}

	return m
}

func OrmStruct2MapLower4Update(s interface{}, skip ...string) map[string]interface{} {
	m := OrmStruct2MapLower(s, skip...)

	n := map[string]interface{}{}
	for k, v := range m {
		if !IsBlank(reflect.ValueOf(v)) {
			n[k] = v
		}
	}

	return n
}

func CopySameFields(fromValue interface{}, toValue interface{}) (err error) {
	var isSlice bool
	var amount = 1

	var deepFields func(reflectType reflect.Type) []reflect.StructField
	var set func(to, from reflect.Value) bool

	indirectType := func(reflectType reflect.Type) reflect.Type {
		for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
			reflectType = reflectType.Elem()
		}
		return reflectType
	}

	indirect := func(reflectValue reflect.Value) reflect.Value {
		for reflectValue.Kind() == reflect.Ptr {
			reflectValue = reflectValue.Elem()
		}
		return reflectValue
	}

	deepFields = func(reflectType reflect.Type) []reflect.StructField {
		var fields []reflect.StructField

		if reflectType = indirectType(reflectType); reflectType.Kind() == reflect.Struct {
			for i := 0; i < reflectType.NumField(); i++ {
				v := reflectType.Field(i)
				if v.Anonymous {
					fields = append(fields, deepFields(v.Type)...)
				} else {
					fields = append(fields, v)
				}
			}
		}

		return fields
	}

	set = func(to, from reflect.Value) bool {
		if from.IsValid() {
			if to.Kind() == reflect.Ptr {
				//set `to` to nil if from is nil
				if from.Kind() == reflect.Ptr && from.IsNil() {
					to.Set(reflect.Zero(to.Type()))
					return true
				} else if to.IsNil() {
					to.Set(reflect.New(to.Type().Elem()))
				}
				to = to.Elem()
			}

			if from.Type().ConvertibleTo(to.Type()) {
				to.Set(from.Convert(to.Type()))
			} else if scanner, ok := to.Addr().Interface().(sql.Scanner); ok {
				err = scanner.Scan(from.Interface())
				if err != nil {
					return false
				}
			} else if from.Kind() == reflect.Ptr {
				return set(to, from.Elem())
			} else {
				return false
			}
		}
		return true
	}

	from := indirect(reflect.ValueOf(fromValue))
	to := indirect(reflect.ValueOf(toValue))

	if !to.CanAddr() {
		return errors.New("copy to value is unaddressable")
	}

	// Return is from value is invalid
	if !from.IsValid() {
		return
	}

	fromType := indirectType(from.Type())
	toType := indirectType(to.Type())

	// Just set it if possible to assign
	// And need to do copy anyway if the type is struct
	if fromType.Kind() != reflect.Struct && from.Type().AssignableTo(to.Type()) {
		to.Set(from)
		return
	}

	if fromType.Kind() != reflect.Struct || toType.Kind() != reflect.Struct {
		return
	}

	if to.Kind() == reflect.Slice {
		isSlice = true
		if from.Kind() == reflect.Slice {
			amount = from.Len()
		}
	}

	for i := 0; i < amount; i++ {
		var dest, source reflect.Value

		if isSlice {
			// source
			if from.Kind() == reflect.Slice {
				source = indirect(from.Index(i))
			} else {
				source = indirect(from)
			}
			// dest
			dest = indirect(reflect.New(toType).Elem())
		} else {
			source = indirect(from)
			dest = indirect(to)
		}

		// check source
		if source.IsValid() {
			fromTypeFields := deepFields(fromType)
			//fmt.Printf("%#v", fromTypeFields)
			// Copy from field to field or method
			for _, field := range fromTypeFields {
				name := field.Name

				if fromField := source.FieldByName(name); fromField.IsValid() {
					// has field
					if toField := dest.FieldByName(name); toField.IsValid() {
						if toField.CanSet() {
							if !set(toField, fromField) {
								if err = CopySameFields(fromField.Interface(), toField.Addr().Interface()); err != nil {
									return err
								}
							}
						}
					} else {
						// try to set to method
						var toMethod reflect.Value
						if dest.CanAddr() {
							toMethod = dest.Addr().MethodByName(name)
						} else {
							toMethod = dest.MethodByName(name)
						}

						if toMethod.IsValid() && toMethod.Type().NumIn() == 1 && fromField.Type().AssignableTo(toMethod.Type().In(0)) {
							toMethod.Call([]reflect.Value{fromField})
						}
					}
				}
			}

			// Copy from method to field
			for _, field := range deepFields(toType) {
				name := field.Name

				var fromMethod reflect.Value
				if source.CanAddr() {
					fromMethod = source.Addr().MethodByName(name)
				} else {
					fromMethod = source.MethodByName(name)
				}

				if fromMethod.IsValid() && fromMethod.Type().NumIn() == 0 && fromMethod.Type().NumOut() == 1 {
					if toField := dest.FieldByName(name); toField.IsValid() && toField.CanSet() {
						values := fromMethod.Call([]reflect.Value{})
						if len(values) >= 1 {
							set(toField, values[0])
						}
					}
				}
			}
		}
		if isSlice {
			if dest.Addr().Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest.Addr()))
			} else if dest.Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest))
			}
		}
	}
	return
}
