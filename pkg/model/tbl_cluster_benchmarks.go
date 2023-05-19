/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// OpenstackBenchmarksTable - 클러스터 Benchmarks 테이블 정보 (Openstack)
type OpenstackBenchmarksTable struct {
	CloudUid        *string    `json:"cloud_uid" db:"cloud_uid"`
	ClusterUid      *string    `json:"cluster_uid" db:"cluster_uid"`
	BenchmarksUid   *string    `json:"benchmarks_uid" db:"benchmarks_uid"`
	CisVersion      *string    `json:"cis_version" db:"cis_version"`
	DetectedVersion *string    `json:"detected_version" db:"detected_version"`
	Results         *Outputs   `json:"results" db:"results"`
	Totals          *Totals    `json:"totals" db:"totals"`
	Status          *int       `json:"status" db:"state"`
	Reason          *string    `json:"reason" db:"reason"`
	Creator         *string    `json:"creator" db:"creator"`
	Created         *time.Time `json:"created_at" db:"created_at"`
}

type Outputs []Output

// Value Marshal
func (o Outputs) Value() (driver.Value, error) {
	return json.Marshal(o)
}

// Scan Unmarshal
func (o *Outputs) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &o)
}

type Output struct {
	NodeName string   `json:"node_name"`
	Controls Controls `json:"controls"`
	Totals   Totals   `json:"totals"`
}

type Controls []Control

// Value Marshal
func (c Controls) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Scan Unmarshal
func (c *Controls) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}

type Control struct {
	Id              string `json:"id"`
	Version         string `json:"version"`
	DetectedVersion string `json:"deteted_version"`
	Text            string `json:"text"`
	NodeType        string `json:"node_type"`
	Tests           Tests  `json:"tests"`
}

type Tests []Test

// Value Marshal
func (t Tests) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// Scan Unmarshal
func (t *Tests) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type Test struct {
	Section string  `json:"section"`
	Type    string  `json:"type"`
	Pass    int     `json:"pass"`
	Fail    int     `json:"fail"`
	Warn    int     `json:"warn"`
	Info    int     `json:"info"`
	Desc    string  `json:"desc"`
	Results Results `json:"results"`
}

type Results []Result

// Value Marshal
func (r Results) Value() (driver.Value, error) {
	return json.Marshal(r)
}

// Scan Unmarshal
func (r *Results) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &r)
}

type Result struct {
	TestNumber     string   `json:"test_number"`
	TestDesc       string   `json:"test_desc"`
	Audit          string   `json:"audit"`
	Type           string   `json:"type"`
	Remediation    string   `json:"remediation"`
	TestInfo       []string `json:"test_info"`
	Status         string   `json:"status"`
	ActualValue    string   `json:"actual_value"`
	Scored         bool     `json:"scored"`
	IsMultiple     bool     `json:"isMultiple"`
	ExpectedResult string   `json:"expected_result"`
	Reason         string   `json:"reason"`
}

type Totals struct {
	TotalPass int `json:"total_pass"`
	TotalFail int `json:"total_fail"`
	TotalWarn int `json:"total_warn"`
	TotalInfo int `json:"total_info"`
}

// Value Marshal
func (t Totals) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// Scan Unmarshal
func (t *Totals) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
