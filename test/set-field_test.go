package test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	mr "github.com/acornsoft-edgecraft/edgecraft-api/pkg/model/response"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
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
	Test *string `json:test`
}

type User2 struct {
	Username *string    `json:"username"`
	Password *string    `json:"password"`
	Codes    *[]string  `json:"codes"`
	Uid      *uuid.UUID `json:uuid`
	ID       *int       `json:ID`
}

type Target struct {
	aaa User2 `json:aaa`
	bbb User1 `json:bbb`
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

func pprint(data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
}
