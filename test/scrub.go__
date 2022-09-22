package test

import (
	"fmt"
	"reflect"
)

// Scrub scrubs all the specified string fields in the 'input' struct
func Scrub(input interface{}, fieldsToScrub interface{}) interface{} {
	if input == nil || fieldsToScrub == nil {
		// Return json representation of 'nil' input
		return "null"
	}

	// Restore all the scrubbed values back to the original values in the struct.
	scrubInternal(input, "", fieldsToScrub)

	// Return the scrubbed string
	return input
}

func scrubInternal(target interface{}, fieldName string, fieldsToScrub interface{}) {

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

			scrubInternal(fValue.Addr().Interface(), fType.Name, fieldsToScrub)
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
			scrubInternal(arrValue.Addr().Interface(), fieldName, fieldsToScrub)
		}

		return
	}

	// If 'fieldName' is not set, then the API was not called on a struct.
	// Since it is not possible to find the variable name of a non-struct field,
	// we can't compare it with 'fieldsToScrub'.
	if fieldName == "" {
		return
	}

	fmt.Println("fieldName: ", fieldName)

	v := reflect.ValueOf(fieldsToScrub).Elem()
	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		if fieldName == typeField.Name {
			targetValue.Set(reflect.ValueOf(v.Field(i).Interface()))
		}
	}
}
