package utils

import (
	"fmt"
	"reflect"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
)

// Iterate through Fields of the Struct
func TypeOfStruct(input interface{}, fieldsToRedact interface{}) (interface{}, error) {
	if input == nil {
		// Return json representation of 'nil' input
		logger.Errorf("input struct nil.")
		return nil, fmt.Errorf(" input struct nil")
	}

	// Get a JSON marshalled string from the scrubb string to return.
	// var b []byte
	// b, _ = json.Marshal(input)

	// Restore all the scrubbed values back to the original values in the struct.
	scrubInternal(fieldsToRedact, "", input)

	// Return the scrubbed string
	return fieldsToRedact, nil
}

func scrubInternal(target interface{}, fieldName string, input interface{}) {

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

			scrubInternal(fValue.Addr().Interface(), fType.Name, input)
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

			scrubInternal(arrValue.Addr().Interface(), fieldName, input)
		}

		return
	}

	// If 'fieldName' is not set, then the API was not called on a struct.
	// Since it is not possible to find the variable name of a non-struct field,
	// we can't compare it with 'fieldsToScrub'.
	if fieldName == "" {
		return
	}

	fmt.Println("sdfgsdfgsd : ", fieldName)

	// if reflect.ValueOf(input).Elem().FieldByName(fieldName) == (reflect.Value{}) {
	// 	fmt.Printf("asdf %s \n", fieldName)
	// }

	// if reflect.ValueOf(input).Elem().FieldByName(fieldName) != (reflect.Value{}) {
	// 	fmt.Printf("asdf %s \n", fieldName)
	// }

}
