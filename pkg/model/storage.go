/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// StorageClassInfo - Data for Storage Class
type StorageClassInfo struct {
	UseCeph bool   `json:"use_ceph" example:"false"`
	Label1  string `json:"label1" example:""`
	Label2  string `json:"label2" example:""`
	Label3  string `json:"label3" example:""`
}

// ToTable - Storage Class 정보를 테이블로 설정
func (sci *StorageClassInfo) ToTable(clusterTable *ClusterTable) {

}

// FromTable - 테이블 정보를 Storage Class 정보로 설정
func (sci *StorageClassInfo) FromTable(clusterTable *ClusterTable) {

}

// StorageClass - Data for Storage Class
type StorageClass struct {
	Use_ceph bool     `json:"use_ceph" db:"-, default:false"`
	Labels   []*Label `json:"labels"`
}

// ToTable - Storage Class 정보를 테이블로 설정
func (sci *StorageClass) ToTable(clusterTable *ClusterTable) {

}

// FromTable - 테이블 정보를 Storage Class 정보로 설정
func (sci *StorageClass) FromTable(clusterTable *ClusterTable) {

}

// Value Marshal
func (a StorageClass) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *StorageClass) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

// TODO: StorageClass 와 StorageClassInfo 차이는? 화면에선 StorageClassInfo를 DB에서는 StorageClass를 사용 중.
