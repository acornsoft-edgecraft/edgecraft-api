/*
 * Copyright (c) 2022 Nutanix Inc. All rights reserved.
 *
 * Author: Shyamsunder Rathi - shyam.rathi@nutanix.com
 * MIT License
 */

package utils

import (
	"fmt"
	"reflect"
)

// Scrub scrubs all the specified string fields in the 'input' struct
func SetFieldsInStruct(fieldsToRedact interface{}, input interface{}) error {
	if input == nil || fieldsToRedact == nil {
		// Return json representation of 'nil' input
		err := fmt.Errorf("is an nil struct ")
		return err
	}

	// Restore all the scrubbed values back to the original values in the struct.
	setField(fieldsToRedact, "", input)

	return nil
}

func setField(target interface{}, fieldName string, fieldsToScrub interface{}) {

	// if target is not pointer, then immediately return
	// modifying struct's field requires addressable object
	addrValue := reflect.ValueOf(target)
	if addrValue.Kind() != reflect.Ptr {
		return
	}

	targetValue := addrValue.Elem()
	if !targetValue.IsValid() {
		return
	}

	targetType := targetValue.Type()

	// If the field/struct is passed by pointer, then first dereference it to get the
	// underlying value (the pointer must not be pointing to a nil value).
	if targetType.Kind() == reflect.Ptr && !targetValue.IsNil() {
		targetValue = targetValue.Elem()
		if !targetValue.IsValid() {
			return
		}

		targetType = targetValue.Type()
	}

	if targetType.Kind() == reflect.Struct {
		// If target is a struct then recurse on each of its field.
		for i := 0; i < targetType.NumField(); i++ {
			fType := targetType.Field(i)
			fValue := targetValue.Field(i)
			if !fValue.IsValid() {
				continue
			}

			if !fValue.CanAddr() {
				// Cannot take pointer of this field, so can't scrub it.
				continue
			}

			if !fValue.Addr().CanInterface() {
				// This is an unexported or private field (begins with lowercase).
				// We can't take an interface on that or scrub it.
				// UnsafeAddr(), which is unsafe.Pointer, can be used to workaround it,
				// but that is not recommended in Golang.
				continue
			}

			setField(fValue.Addr().Interface(), fType.Name, fieldsToScrub)
		}
		return
	}

	if targetType.Kind() == reflect.Array || targetType.Kind() == reflect.Slice {
		// If target is an array/slice, then recurse on each of its element.
		for i := 0; i < targetValue.Len(); i++ {
			arrValue := targetValue.Index(i)
			if !arrValue.IsValid() {
				continue
			}

			if !arrValue.CanAddr() {
				// Cannot take pointer of this field, so can't scrub it.
				continue
			}

			if !arrValue.Addr().CanInterface() {
				// This is an unexported or private field (begins with lowercase).
				// We can't take an interface on that or scrub it.
				// UnsafeAddr(), which is unsafe.Pointer, can be used to workaround it,
				// but that is not recommended in Golang.
				continue
			}
			setField(arrValue.Addr().Interface(), fieldName, fieldsToScrub)
		}

		return
	}

	// If 'fieldName' is not set, then the API was not called on a struct.
	// Since it is not possible to find the variable name of a non-struct field,
	// we can't compare it with 'fieldsToScrub'.
	if fieldName == "" {
		return
	}

	v := reflect.ValueOf(fieldsToScrub).Elem()
	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		if fieldName == typeField.Name {
			targetValue.Set(reflect.ValueOf(v.Field(i).Interface()))
		}
	}

	// var data interface{}
	// if data = returnValue(fieldsToScrub, fieldName, ""); data != nil {
	// 	fmt.Printf("data : %v\n", data)
	// 	targetValue.Set(reflect.ValueOf(data))
	// }
	// setFieldsValue(targetValue, fieldsToScrub, fieldName, "")

}

func setFieldsValue(x interface{}, fn func(input string)) {
	val := getValue(x)

	numberOfValues := 0
	var getField func(int) reflect.Value

	switch val.Kind() {
	case reflect.String:
		fn(val.String())
	case reflect.Struct:
		numberOfValues = val.NumField()
		getField = val.Field
	case reflect.Slice:
		numberOfValues = val.Len()
		getField = val.Index
	}

	for i := 0; i < numberOfValues; i++ {
		setFieldsValue(getField(i).Interface(), fn)
	}

	// addrValue := reflect.ValueOf(target)
	// if addrValue.Kind() != reflect.Ptr {
	// 	return
	// }

	// targetValue := addrValue.Elem()
	// if !targetValue.IsValid() {
	// 	return
	// }

	// targetType := targetValue.Type()
	// if targetType.Kind() == reflect.Ptr && !targetValue.IsNil() {
	// 	targetValue = targetValue.Elem()
	// 	if !targetValue.IsValid() {
	// 		return
	// 	}

	// 	targetType = targetValue.Type()
	// }

	// // It's possible we can cache this, which is why precompute all these ahead of time.
	// findJsonName := func(t reflect.StructTag) (string, error) {
	// 	if jt, ok := t.Lookup("json"); ok {
	// 		return strings.Split(jt, ",")[0], nil
	// 	}
	// 	return "", nil
	// }

	// if targetType.Kind() == reflect.Struct {
	// 	// If target is a struct then recurse on each of its field.
	// 	for i := 0; i < targetType.NumField(); i++ {
	// 		fType := targetType.Field(i)
	// 		fValue := targetValue.Field(i)
	// 		if !fValue.IsValid() {
	// 			continue
	// 		}

	// 		if !fValue.CanAddr() {
	// 			// Cannot take pointer of this field, so can't scrub it.
	// 			continue
	// 		}

	// 		if !fValue.Addr().CanInterface() {
	// 			// This is an unexported or private field (begins with lowercase).
	// 			// We can't take an interface on that or scrub it.
	// 			// UnsafeAddr(), which is unsafe.Pointer, can be used to workaround it,
	// 			// but that is not recommended in Golang.
	// 			continue
	// 		}

	// 		tag := fType.Tag
	// 		jname, _ := findJsonName(tag)

	// 		if jname != "" && jname != "-" {
	// 			setFieldsValue(setValue, fValue.Addr().Interface(), fieldName, fType.Name)
	// 		}
	// 	}
	// 	return
	// }

	// if name == "" {
	// 	return
	// }

	// if fieldName == name {
	// 	fmt.Println("fieldName: ", name)
	// 	fmt.Printf("--type : %v\n", targetValue.Interface())
	// 	setValue.Set(reflect.ValueOf(targetValue.Interface()))
	// 	return
	// }
}

func getValue(x interface{}) reflect.Value {
	val := reflect.ValueOf(x)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return val
}

func walk(x interface{}, fn func(input string)) {
	val := getValue(x)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		switch field.Kind() {
		case reflect.String:
			fn(field.String())
		case reflect.Struct:
			walk(field.Interface(), fn)
		}
	}
}
