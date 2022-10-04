/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import (
	"time"
)

// CloudTable - Cloud Table 정보
// Type: Common Code CloudTypes 참조
// Status: Common Code Status 참조
type CloudTable struct {
	CloudUID *string    `json:"cloud_uid" db:"cloud_uid"`
	Name     *string    `json:"name" db:"name"`
	Type     *int       `json:"type" db:"type"`
	Desc     *string    `json:"desc" db:"description"`
	Status   *int       `json:"status" db:"state"`
	Creator  *string    `json:"creator" db:"creator"`
	Created  *time.Time `json:"created" db:"created_at"`
	Updater  *string    `json:"updater" db:"updater"`
	Updated  *time.Time `json:"updated" db:"updated_at"`
}

// ToSet - Mapping cloud data to CloudSet
func (ct *CloudTable) ToSet(cloudSet *CloudSet) {
	var cloud *CloudInfo = &CloudInfo{}
	cloud.FromTable(ct)

	cloudSet.Cloud = cloud
}
