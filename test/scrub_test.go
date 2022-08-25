package scrub

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Have a struct with some sensitive fields.
type testScrub struct {
	Username string
	Password string
	Codes    []string
}

// Structure definitions to test scrubbing functionalities.
// Simple struct
type User struct {
	Username  string
	Password  string
	DbSecrets []string
}

// Nested struct
type Users struct {
	Secret   string
	Keys     []string
	UserInfo []User
}

// Test
func TestScrub(t *testing.T) {
	// Create a struct with some sensitive data.
	T := testScrub{
		Username: "administrator",
		Password: "my_secret_passphrase",
		Codes:    []string{"pass1", "pass2", "pass3"},
	}

	// Create a set of field names to scrub (default is 'password').
	fieldsToScrub := map[string]bool{
		"password": true,
		"codes":    true,
	}
	// Call the util API to get a JSON formatted string with scrubbed field values.
	out := Scrub(&T, fieldsToScrub)
	// Log the scrubbed string without worrying about prying eyes!
	log.Println(out)
}

// TestScrubSimple tests scrubbing on a simple struct with default
// sensitive fields.
func TestScrubSimple(t *testing.T) {
	user := &User{
		Username:  "Shyam Rathi",
		Password:  "a",
		DbSecrets: []string{"db_secret_1", "db_secret_2"},
	}

	userScrubbed := &User{
		Username:  "Shyam Rathi",
		Password:  "********",
		DbSecrets: []string{"db_secret_1", "db_secret_2"},
	}

	validateScrub(t, user, userScrubbed, nil)
}

// TestScrubNested tests scrubbing on a nested complex struct with
// specified sensitive fields.
func TestScrubNested(t *testing.T) {
	users := &Users{
		Secret: "secret_sshhh",
		Keys:   []string{"key_1", "key_2", "key_3"},
		UserInfo: []User{
			{
				Username:  "John Doe",
				Password:  "John_Doe's_Password",
				DbSecrets: []string{"John's_db_secret_1", "John's_db_secret_2"},
			},
			{
				Username:  "Jane Doe",
				Password:  "Jane_Doe's_Password",
				DbSecrets: []string{"Jane's_db_secret_1", "Jane's_db_secret_2"},
			},
		},
	}

	userScrubbed := &Users{
		Secret: "********",
		Keys:   []string{"********", "********", "********"},
		UserInfo: []User{
			{
				Username:  "John Doe",
				Password:  "********",
				DbSecrets: []string{"********", "********"},
			},
			{
				Username:  "Jane Doe",
				Password:  "********",
				DbSecrets: []string{"********", "********"},
			},
		},
	}

	secretFields := map[string]bool{
		"password": true, "keys": true, "secret": true, "dbsecrets": true}
	validateScrub(t, users, userScrubbed, secretFields)
}

// TestScrubNil tests scrubbing on a empty or nil input.
func TestScrubNil(t *testing.T) {
	user := &User{
		Username:  "",
		Password:  "nutanix/4u",
		DbSecrets: []string{},
	}

	userScrubbed := &User{
		Username:  "",
		Password:  "********",
		DbSecrets: []string{},
	}

	// Validate input with empty fields
	validateScrub(t, user, userScrubbed, nil)

	// Validate empty pointer input
	var userEmpty *User
	validateScrub(t, userEmpty, userEmpty, nil)

	// Validate nil input
	validateScrub(t, nil, nil, nil)
}

// TestScrubNestedNil tests scrubbing on a nested complex struct with
// some nil, empty and specified sensitive fields.
func TestScrubNestedNil(t *testing.T) {
	users := &Users{
		Secret: "",
		Keys:   nil,
		UserInfo: []User{
			{
				Username:  "John Doe",
				Password:  "",
				DbSecrets: []string{"John's_db_secret_1", "John's_db_secret_2"},
			},
			{
				Username:  "Jane Doe",
				Password:  "Jane_Doe's_Password",
				DbSecrets: []string{},
			},
		},
	}

	userScrubbed := &Users{
		Secret: "",
		Keys:   nil,
		UserInfo: []User{
			{
				Username:  "John Doe",
				Password:  "",
				DbSecrets: []string{"********", "********"},
			},
			{
				Username:  "Jane Doe",
				Password:  "********",
				DbSecrets: []string{},
			},
		},
	}

	// Test a nested struct with some empty and nil fields.
	secretFields := map[string]bool{
		"password": true, "keys": true, "secret": true, "dbsecrets": true}
	validateScrub(t, users, userScrubbed, secretFields)
}

// validateScrub is a helper function to validate scrubbing functionality on a struct.
func validateScrub(t *testing.T, msg, scrubbedMsg interface{}, secretFields map[string]bool) {
	t.Helper()

	// Get the scrubbed string from util API.
	got := Scrub(msg, secretFields)
	log.Println(got)

	// Compare it against the given scrubbed representaation.
	var b []byte
	b, _ = json.Marshal(scrubbedMsg)
	want := string(b)

	assert.Equal(t, want, got,
		"JSON representation mismatch after scrubbing sensitive fields")
}
