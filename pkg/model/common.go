/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Label struct {
	Key   string `json:"key" example:""`
	Value string `json:"value" example:""`
}

type Labels []interface{}

// Value Marshal
func (a Labels) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *Labels) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

// UrlInfo - Data for URL
type UrlInfo struct {
	IpAddress string `json:"ip_address" example:""`
	Port      string `json:"port" example:""`
}
