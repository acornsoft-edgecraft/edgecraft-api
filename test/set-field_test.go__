package test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	mr "github.com/acornsoft-edgecraft/edgecraft-api/pkg/model/response"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/emirpasic/gods/stacks/arraystack"
	"github.com/gofrs/uuid"
)

// Structure definitions to test scrubbing functionalities.
// Simple struct
type User struct {
	Username string `json:"username,omitempty"`
	Password string
	Codes    []string
}

type User1 struct {
	Test *string `json:"test"`
}

type User2 struct {
	Username *string    `json:"username"`
	Password *string    `json:"password"`
	Codes    *[]string  `json:"codes"`
	Uid      *uuid.UUID `json:"uuid"`
	ID       *int       `json:"ID"`
}

type Target struct {
	aaa User2 `json:"aaa"`
	bbb User1 `json:"bbb"`
}

type Personal struct {
	Name *string
	Age  *int
}

type Adde struct {
	Person
}

func NewPerson(name string, age int) *Personal {
	return &Personal{
		Name: &name,
		Age:  &age,
	}
}

// Test
func TestSetField(t *testing.T) {
	// Create a struct with some sensitive data.
	w := User{
		Username: "",
		Password: "my_secret_passphrase",
		Codes:    []string{"pass1", "pass2", "pass3"},
	}

	// Create a set of field names to scrub (default is 'password').
	// fieldsToScrub := map[string]bool{
	// 	"password": true,
	// 	"codes":    true,
	// }

	// Call the util API to get a JSON formatted string with scrubbed field values.
	// out := Scrub(&T, fieldsToScrub)
	SetField(&w, "username", "asdfsadf")

	// Log the scrubbed string without worrying about prying eyes!
	// log.Println(out)

	pprint(w)
}

func TestSetField2(t *testing.T) {
	w := new(mr.RegisterCloud)
	uuid := uuid.Must(uuid.NewV4())

	checkConfig2(w, "CloudName", "asdfasdf")
	checkConfig2(w, "CloudUID", uuid)
	pprint(w)
}

func TestSetField3(t *testing.T) {
	w := new(User2)
	field := reflect.ValueOf(&w).Elem().Field(0)
	SetUnexportedField(field, "124351qewf")
	pprint(w)
}

func TestSetField4(t *testing.T) {
	// w := new(mr.RegisterCloud)
	var w mr.RegisterCloud
	a := "asdfasdf"
	uuid := uuid.Must(uuid.NewV4())
	ww := model.Cloud{
		CloudName: &a,
		CloudUID:  &uuid,
	}

	// Scrub(&w, &ww)
	utils.SetFieldsInStruct(&w, &ww)

	pprint(w)
}

func TestIsEqual(t *testing.T) {
	var aVal []Personal

	for i := 0; i < 10; i++ {
		a := NewPerson("aaa", 18)
		aVal = append(aVal, *a)
	}

	aa := NewPerson("aaa", 18)
	bb := NewPerson("vvv", 30)
	cc := IsEqual(aa, bb)
	fmt.Println("== ", cc)
}

func TestIsArrayStack(t *testing.T) {
	stack := arraystack.New() // empty
	stack.Push(1)             // 1
	stack.Push(2)             // 1, 2
	stack.Values()            // 2, 1 (LIFO order)
	pprint(stack.Values())

	_, _ = stack.Peek() // 2,true
	pprint(stack.Values())
	_, _ = stack.Pop() // 2, true
	pprint(stack.Values())
	_, _ = stack.Pop() // 1, true
	pprint(stack.Values())
	_, _ = stack.Pop() // nil, false (nothing to pop)
	stack.Push(1)      // 1
	pprint(stack.Values())
	stack.Clear() // empty
	stack.Empty() // true
	stack.Size()  // 0
}

func TestIsArrayList(t *testing.T) {
	var aa []interface{}
	for i := 0; i < 10; i++ {
		aa = append(aa, uuid.Must(uuid.NewV4()))
	}
	bb := []interface{}{
		"f35ec01a-60e0-40c5-8581-93b18f5cb75f",
		"67965b3e-8118-4872-82bd-3a15b8c6263d",
		"27b1ac06-de9c-4c64-9b77-e6c943044a51",
		"6e6d78e0-b26d-41bb-8f5c-abc85ab4b8f8",
	}

	pprint(aa)
	pprint(bb)

	cc := arraylist.New()
	cc.Add("f35ec01a-60e0-40c5-8581-93b18f5cb75f")
	cc.Add("67965b3e-8118-4872-82bd-3a15b8c6263d")
	cc.Add("27b1ac06-de9c-4c64-9b77-e6c943044a51")
	cc.Add("6e6d78e0-b26d-41bb-8f5c-abc85ab4b8f8")
	cc.Add("f35ec01a-60e0-40c5-8581-93b18f5cb75f")

	pprint(cc)

	set := hashset.New()
	for i := 0; i < 10; i++ {
		set.Add(uuid.Must(uuid.NewV4()))
	}
	set.Add(bb...)
	pprint("--set--")
	pprint(set)

	set1 := hashset.New()

	set1.Add("95b15da8-65ef-49a2-a15c-dfcbb6f87ddf")
	set1.Add("afa45bb6-3194-47bb-9648-f7a556801e79")
	set1.Add("86f0df12-9b73-44cd-b453-4f8a4f729230")
	set1.Add("2bd0145c-633b-4a15-9f8f-df3edff817d4")
	set1.Add("f35ec01a-60e0-40c5-8581-93b18f5cb75f")
	set1.Add("95b15da8-65ef-49a2-a15c-dfcbb6f87ddf")
	set1.Add("de9486b4-9eff-4c35-9163-9abfd443e5d8")
	set1.Add("db44e496-69c1-4e03-b8f7-4c52328910db")
	set1.Add("9f8fb895-276f-4950-b2ce-e9dc5b76d7a3")
	pprint("--set1--")
	pprint(set1)

	set2 := hashset.New()
	set2.Add("95b15da8-65ef-49a2-a15c-dfcbb6f87ddf")
	set2.Add("afa45bb6-3194-47bb-9648-f7a556801e79")
	set2.Add("86f0df12-9b73-44cd-b453-4f8a4f729230")
	set2.Add("2bd0145c-633b-4a15-9f8f-df3edff817d4")
	set2.Add("f35ec01a-60e0-40c5-8581-93b18f5cb75f")
	set2.Add("95b15da8-65ef-49a2-a15c-dfcbb6f87ddf")

	pprint(set2)
	aaa := set1.Contains(set2)
	pprint(aaa)

	set1.Remove(set2)
	pprint(set1)

	tt := arraylist.New()
	tt.Add("95b15da8-65ef-49a2-a15c-dfcbb6f87ddf")
	tt.Add("afa45bb6-3194-47bb-9648-f7a556801e79")
	tt.Add("86f0df12-9b73-44cd-b453-4f8a4f729230")
	tt.Add("2bd0145c-633b-4a15-9f8f-df3edff817d4")
	tt.Add("f35ec01a-60e0-40c5-8581-93b18f5cb75f")
	tt.Add("95b15da8-65ef-49a2-a15c-dfcbb6f87ddf")
	tt.Add("de9486b4-9eff-4c35-9163-9abfd443e5d8")
	tt.Add("db44e496-69c1-4e03-b8f7-4c52328910db")
	tt.Add("9f8fb895-276f-4950-b2ce-e9dc5b76d7a3")

	gg := arraylist.New()
	gg.Add("95b15da8-65ef-49a2-a15c-dfcbb6f87ddf")
	gg.Add("afa45bb6-3194-47bb-9648-f7a556801e79")
	gg.Add("86f0df12-9b73-44cd-b453-4f8a4f729230")
	gg.Add("2bd0145c-633b-4a15-9f8f-df3edff817d4")
	gg.Add("f35ec01a-60e0-40c5-8581-93b18f5cb75f")

	pprint(tt.Contains(gg))

	pprint(tt)
	pprint(gg)

	hhh, _ := gg.Get(0)
	ff := tt.Contains(hhh)

	pprint("--result--")
	pprint(ff)
	pprint(hhh)

	pprint("--remove--")
	vv := hashset.New()
	vv.Add("86f0df12-9b73-44cd-b453-4f8a4f729230")
	vv.Add("95b15da8-65ef-49a2-a15c-dfcbb6f87ddf")
	// vvv := vv.Values()
	// set1.Remove(vvv)

	// set1.Remove("86f0df12-9b73-44cd-b453-4f8a4f729230", "95b15da8-65ef-49a2-a15c-dfcbb6f87ddf")
	pprint(set1)

	pprint("--HashMap--")
	m := hashmap.New()
	m.Put(0, "95b15da8-65ef-49a2-a15c-dfcbb6f87ddf")
	m.Put(1, "afa45bb6-3194-47bb-9648-f7a556801e79")
	m.Put(2, "afa45bb6-3194-47bb-9648-f7a556801e79")
	m.Put(3, "86f0df12-9b73-44cd-b453-4f8a4f729230")
	m.Put(4, "2bd0145c-633b-4a15-9f8f-df3edff817d4")
	m.Put(5, "95b15da8-65ef-49a2-a15c-dfcbb6f87ddf")
	m.Put(6, "de9486b4-9eff-4c35-9163-9abfd443e5d8")
	m.Put(7, "db44e496-69c1-4e03-b8f7-4c52328910db")
	m.Put(8, "9f8fb895-276f-4950-b2ce-e9dc5b76d7a3")
	pprint(m)

	mm := hashmap.New()
	mm.Put(0, "95b15da8-65ef-49a2-a15c-dfcbb6f87ddf")
	mm.Put(1, "afa45bb6-3194-47bb-9648-f7a556801e79")
	mm.Put(2, "afa45bb6-3194-47bb-9648-f7a556801e79")
	mm.Put(3, "86f0df12-9b73-44cd-b453-4f8a4f729230")
	mm.Put(4, "2bd0145c-633b-4a15-9f8f-df3edff817d4")
	mm.Put(5, "95b15da8-65ef-49a2-a15c-dfcbb6f87ddf")
	mm.Put(6, "de9486b4-9eff-4c35-9163-9abfd443e5d8")
	mm.Put(7, "db44e496-69c1-4e03-b8f7-4c52328910db")
	mm.Put(8, "9f8fb895-276f-4950-b2ce-e9dc5b76d7a3")
	pprint(mm)

	list1 := arraylist.New()
	list2 := arraylist.New()

	list1.Add("95b15da8-65ef-49a2-a15c-dfcbb6f87ddf")
	list1.Add("afa45bb6-3194-47bb-9648-f7a556801e79")
	list1.Add("afa45bb6-3194-47bb-9648-f7a556801e79")
	list1.Add("86f0df12-9b73-44cd-b453-4f8a4f729230")

	list2.Add("afa45bb6-3194-47bb-9648-f7a556801e79")

}

func TestIsHashSet(t *testing.T) {
	aa := []interface{}{
		"f35ec01a-60e0-40c5-8581-93b18f5cb75f",
		"67965b3e-8118-4872-82bd-3a15b8c6263d",
		"27b1ac06-de9c-4c64-9b77-e6c943044a51",
		"6e6d78e0-b26d-41bb-8f5c-abc85ab4b8f8",
		"9f8fb895-276f-4950-b2ce-e9dc5b76d7a3",
	}
	bb := []interface{}{
		"f35ec01a-60e0-40c5-8581-93b18f5cb75f",
		"6e6d78e0-b26d-41bb-8f5c-abc85ab4b8f8",
		"9f8fb895-276f-4950-b2ce-e9dc5b76d7a3",
		"86f0df12-9b73-44cd-b453-4f8a4f729230",
	}

	set1 := hashset.New()
	set2 := hashset.New()

	set1.Add(aa...)
	set2.Add(bb...)

	pprint("--hashset--")
	pprint(set1)
	pprint(set2)

	set1.Remove(bb...)
	pprint("--remove--")
	pprint(set1.Values())

	pprint("--union--")
	union := set1.Union(set2)
	pprint(union)

	pprint("--Difference--")
	difference := set1.Difference(set2)
	pprint(difference)

}

func pprint(data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b), "\n")
}

func IsEqual(A, B interface{}) bool {
	// defer execTime("IsEqual took", time.Now())
	// Find out the type of A & B is Person or not
	if _, ok := A.(*Personal); ok {
		if _, ok := B.(*Personal); ok {
			if A.(*Personal).Name == B.(*Personal).Name {
				fmt.Println("asdfasdfas")
				return A.(*Personal).Age == B.(*Personal).Age
			} else {
				return false
			}
		}
		return false
	}
	return false
}

func TestCustomComparator(t *testing.T) {

	type Custom struct {
		id   int
		name string
	}

	byID := func(a, b interface{}) int {
		c1 := a.(Custom)
		c2 := b.(Custom)
		switch {
		case c1.id > c2.id:
			return 1
		case c1.id < c2.id:
			return -1
		default:
			return 0
		}
	}

	// o1,o2,expected
	tests := [][]interface{}{
		{Custom{1, "a"}, Custom{1, "a"}, 0},
		{Custom{1, "a"}, Custom{2, "b"}, -1},
		{Custom{2, "b"}, Custom{1, "a"}, 1},
		{Custom{1, "a"}, Custom{1, "b"}, 0},
	}

	for _, test := range tests {
		actual := byID(test[0], test[1])
		expected := test[2]
		pprint(actual)
		fmt.Println("----------")
		// pprint(expected)
		if actual != expected {
			t.Errorf("Got %v expected %v", actual, expected)
		}
	}
}
