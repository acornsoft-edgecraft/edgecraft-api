package entities

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type TestAPI struct {
	ID           *uuid.UUID  `json:"id" gorm:"type:uuid;column:id;default:uuid_generate_v3()"`
	CreatedAt    time.Time   `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time   `json:"updated_at" gorm:"column:updated_at"`
	API_url      string      `json:"api_url" gorm:"column:api_url"`
	API_method   string      `json:"api_method" gorm:"column:api_method"`
	API_status   string      `json:"api_status" gorm:"column:api_status"`
	API_response APIResponse `json:"api_response" gorm:"column:api_response"`
	API_JSONB    JSONB       `json:"api_jsonb" gorm:"column:api_jsonb"`
}

//- Start - JSONB Interface for JSONB Field of yourTableName Table
type JSONB []interface{}

// Value Marshal
func (a JSONB) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

//- END -

//- Start - APIResponse struct to describe API attributes.
type APIResponse struct {
	Data    interface{} `json:"data"`
	Error   bool        `json:"error"`
	Message string      `json:"msg"`
}

// Value make the APIResponse struct implement the driver.Valuer interface.
// This method simply returns the JSON-encoded representation of the struct.
func (b APIResponse) Value() (driver.Value, error) {
	return json.Marshal(b)
}

// Scan make the APIResponse struct implement the sql.Scanner interface.
// This method simply decodes a JSON-encoded value into the struct fields.
func (b *APIResponse) Scan(value interface{}) error {
	j, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(j, &b)
}

//- END -
