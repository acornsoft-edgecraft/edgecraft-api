/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
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
	// change to StorageClass
	var storageClass *StorageClass = &StorageClass{}
	storageClass.UseCeph = sci.UseCeph
	if sci.Label1 != "" {
		keyVal := strings.Split(sci.Label1, "=")
		var label *Label = &Label{Key: keyVal[0], Value: keyVal[1]}
		storageClass.Labels = append(storageClass.Labels, label)
	}
	if sci.Label2 != "" {
		keyVal := strings.Split(sci.Label2, "=")
		var label *Label = &Label{Key: keyVal[0], Value: keyVal[1]}
		storageClass.Labels = append(storageClass.Labels, label)
	}
	if sci.Label3 != "" {
		keyVal := strings.Split(sci.Label3, "=")
		var label *Label = &Label{Key: keyVal[0], Value: keyVal[1]}
		storageClass.Labels = append(storageClass.Labels, label)
	}

	clusterTable.StorageClass = storageClass
}

// FromTable - 테이블 정보를 Storage Class 정보로 설정
func (sci *StorageClassInfo) FromTable(clusterTable *ClusterTable) {
	if clusterTable.StorageClass != nil {
		sci.UseCeph = clusterTable.StorageClass.UseCeph
		for i, label := range clusterTable.StorageClass.Labels {
			if i == 0 {
				sci.Label1 = label.Key + "=" + label.Value
			} else if i == 1 {
				sci.Label2 = label.Key + "=" + label.Value
			} else if i == 2 {
				sci.Label3 = label.Key + "=" + label.Value
			}
		}
	}
}

// StorageClass - Data for Storage Class
type StorageClass struct {
	UseCeph bool     `json:"use_ceph"`
	Labels  []*Label `json:"labels"`
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
