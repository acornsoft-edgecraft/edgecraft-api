package test

import (
	"fmt"
	"reflect"
	"strings"
)

func SetField(item interface{}, fieldName string, value interface{}) error {
	v := reflect.ValueOf(item).Elem()
	fmt.Printf("%v\n", v.FieldByName(fieldName))
	if !v.CanAddr() {
		return nil
	}
	// It's possible we can cache this, which is why precompute all these ahead of time.
	findJsonName := func(t reflect.StructTag) (string, error) {
		if jt, ok := t.Lookup("json"); ok {
			return strings.Split(jt, ",")[0], nil
		}
		return "", nil
	}
	fieldNames := map[string]int{}
	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		tag := typeField.Tag
		jname, _ := findJsonName(tag)
		fieldNames[jname] = i
	}

	fieldNum, ok := fieldNames[fieldName]
	if !ok {
		return fmt.Errorf("field %s does not exist within the provided item", fieldName)
	}
	fieldVal := v.Field(fieldNum)
	fieldVal.Set(reflect.ValueOf(value))
	return nil
}
