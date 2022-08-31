package test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
)

type Person struct {
	Name    string
	Profile Profile
}

type Profile struct {
	Age  int
	City string
}

type Home struct {
	HostName *string
	Address  *string
	Group    Group
}

type Group struct {
	Name *string
	Age  *int
}

func TestWalk(t *testing.T) {
	aa := "asdfasdf"
	bb := 33
	cases := []struct {
		Name          string
		Input         interface{}
		ExpectedCalls []string
	}{
		{
			"Struct with one string field",
			struct {
				Name string
			}{"Chris"},
			[]string{"Chris"},
		},
		{
			"Struct with two string fields",
			struct {
				Name string
				City string
			}{"Chris", "London"},
			[]string{"Chris", "London"},
		},
		{
			"Struct with non string field",
			struct {
				Name string
				Age  int
			}{"Chris", 33},
			[]string{"Chris"},
		},
		{
			"Nested fields",
			struct {
				Name    string
				Profile struct {
					Age  int
					City string
				}
			}{"Chris", struct {
				Age  int
				City string
			}{33, "London"}},
			[]string{"Chris", "London"},
		},
		{
			"Nested fields use Type",
			Person{
				"Chris",
				Profile{33, "London"},
			},
			[]string{"Chris", "London"},
		},
		{
			"Pointers to things",
			&Person{
				"Chris",
				Profile{33, "London"},
			},
			[]string{"Chris", "London"},
		},
		{
			"Pointers to things II",
			&Group{
				&aa,
				&bb,
			},
			[]string{"Chris", "London"},
		},
		{
			"Slices",
			[]Profile{
				{33, "London"},
				{34, "Reykjavík"},
			},
			[]string{"London", "Reykjavík"},
		},
		{
			"Arrays",
			[2]Profile{
				{33, "London"},
				{34, "Reykjavík"},
			},
			[]string{"London", "Reykjavík"},
		},
		{
			"Maps",
			map[string]string{
				"Foo": "Bar",
				"Baz": "Boz",
			},
			[]string{"Bar", "Boz"},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			var got []interface{}
			walk(test.Input, func(input interface{}) {
				got = append(got, input)
			})

			utils.Print(got)
			if !reflect.DeepEqual(got, test.ExpectedCalls) {
				t.Errorf("got %v, want %v", got, test.ExpectedCalls)
			}
		})
	}

}

func TestSet11(t *testing.T) {
	// var group Group
	// var home Home

	a := "asdf"
	b := "192.1234.1234"
	c := "name"
	d := 233

	dd := &Group{
		Name: &c,
		Age:  &d,
	}

	input := &Home{
		HostName: &a,
		Address:  &b,
		Group:    *dd,
	}

	utils.Print(input)
	testSetWalk(input, "", input, "")
	fmt.Println("After")
	fmt.Println()
	utils.Print(input)
}

func testSetWalk(x interface{}, fieldName string, b interface{}, inputName string) {
	val := getValue(x)

	numberOfValues := 0
	var getField func(int) reflect.Value

	// inputVal := getValue(b)
	// numberOfInputValues := 0
	// var getInputField func(int) reflect.Value

	switch val.Kind() {
	case reflect.Struct:
		numberOfValues = val.NumField()
		getField = val.Field
	case reflect.Slice, reflect.Array:
		numberOfValues = val.Len()
		getField = val.Index
	case reflect.Map:
		for _, key := range val.MapKeys() {
			testSetWalk(val.MapIndex(key).Interface(), fieldName, b, inputName)
		}
	}

	for i := 0; i < numberOfValues; i++ {
		testSetWalk(getField(i).Interface(), getField(i).Type().Name(), b, inputName)
	}
	fmt.Println(fieldName)
}

func walk(x interface{}, fn func(input interface{})) {
	val := getValue(x)

	numberOfValues := 0
	var getField func(int) reflect.Value

	switch val.Kind() {
	// case reflect.String:
	// 	fn(val.String())
	// case reflect.Int:
	// 	fn(val.Int())
	case reflect.Struct:
		numberOfValues = val.NumField()
		getField = val.Field
	case reflect.Slice, reflect.Array:
		numberOfValues = val.Len()
		getField = val.Index
	case reflect.Map:
		for _, key := range val.MapKeys() {
			walk(val.MapIndex(key).Interface(), fn)
		}
	}

	for i := 0; i < numberOfValues; i++ {
		walk(getField(i).Interface(), fn)
	}

	// fmt.Println(val)
}

func getValue(x interface{}) reflect.Value {
	val := reflect.ValueOf(x)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return val
}
