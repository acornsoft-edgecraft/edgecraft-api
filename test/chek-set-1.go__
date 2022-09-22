package test

import (
	"fmt"
	"reflect"

	"github.com/gofrs/uuid"
)

func checkConfig2(cnf interface{}, key string, value interface{}) bool {
	v := reflect.ValueOf(cnf).Elem()
	f := v.FieldByName(key)
	fmt.Printf("type: %T\n", value)
	fmt.Printf("type: %v\n", f.IsValid())

	switch v := value.(type) {
	case int:
		v = value.(int)
		if f.IsValid() {
			if f.IsNil() && f.CanSet() {
				f.Set(reflect.ValueOf(&v))
				return true
			}
		}
	case string:
		v = value.(string)
		if f.IsValid() {
			fmt.Printf("asdf--: %v", v)
			if f.IsNil() && f.CanSet() {
				f.Set(reflect.ValueOf(&v))
				fmt.Printf("asdf--: %v", v)
				return true
			}
		}
	case uuid.UUID:
		v = value.(uuid.UUID)
		if f.IsValid() {
			if f.IsNil() && f.CanSet() {
				f.Set(reflect.ValueOf(&v))
				return true
			}
		}
	}
	return false
}
