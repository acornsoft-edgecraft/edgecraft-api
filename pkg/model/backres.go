/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import (
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
	"github.com/gofrs/uuid"
)

// type OpenstackBackResList struct {
// 	CloudUid      *string    `json:"cloud_uid" db:"cloud_uid"`
// 	ClusterUid    *string    `json:"cluster_uid" db:"cluster_uid"`
// 	BenchmarksUid *string    `json:"benchmarks_uid" db:"benchmarks_uid"`
// 	Totals        *string    `json:"totals" db:"totals"`
// 	SuccessYn     *bool      `json:"success_yn" db:"success_yn"`
// 	Reason        *string    `json:"reason" db:"reason"`
// 	Created       *time.Time `json:"created" db:"created_at"`
// }

// type OpenstackBackResSet struct {
// 	CloudUid        string    `json:"cloud_uid" db:"cloud_uid"`
// 	ClusterUid      string    `json:"cluster_uid" db:"cluster_uid"`
// 	BenchmarksUid   string    `json:"benchmarks_uid" db:"benchmarks_uid"`
// 	CisVersion      string    `json:"cis_version" db:"cis_version"`
// 	DetectedVersion string    `json:"detected_version" db:"detected_version"`
// 	Results         string    `json:"results" db:"results"`
// 	Totals          string    `json:"totals" db:"totals"`
// 	SuccessYn       bool      `json:"success_yn" db:"success_yn"`
// 	Reason          string    `json:"reason" db:"reason"`
// 	Created         time.Time `json:"created" db:"created_at"`
// }

// func (ob *OpenstackBackResSet) NewKey() {
// 	ob.BackResUidBenchmarksUid = uuid.Must(uuid.NewV4()).String()
// }

// func (ob *OpenstackBenchmarksSet) ToTable(cloudId *string, clusterId *string, user string, at time.Time) (benchmarksTable *OpenstackBenchmarksTable) {
// 	benchmarksTable = &OpenstackBenchmarksTable{
// 		CloudUid:      cloudId,
// 		ClusterUid:    clusterId,
// 		BenchmarksUid: utils.StringPtr(ob.BenchmarksUid),
// 		SuccessYn:     utils.BoolPtr(false),
// 		Creator:       utils.StringPtr(user),
// 		Created:       utils.TimePtr(at),
// 	}
// 	return
// }

type BackResInfo struct {
	CloudUid   string    `json:"cloud_uid" db:"cloud_uid"`
	ClusterUid string    `json:"cluster_uid" db:"cluster_uid"`
	BackResUid string    `json:"backres_uid" db:"backres_uid"`
	Name       string    `json:"name" db:"name"`
	Type       string    `json:"type" db:"type"`
	Status     string    `json:"status" db:"status"`
	Reason     string    `json:"reasen" db:"reaseon"`
	Created    time.Time `json:"created_at" db:"created_at"`
}

// NewKey - Backup/Restore Unique Key 생성
func (bri *BackResInfo) NewKey() {
	bri.BackResUid = uuid.Must(uuid.NewV4()).String()
}

// ToTable - Backup/Restore 정보를 Table 정보로 전환
func (bri *BackResInfo) ToTable(user string, at time.Time) *BackResTable {
	return &BackResTable{
		CloudUid:   utils.StringPtr(bri.CloudUid),
		ClusterUid: utils.StringPtr(bri.ClusterUid),
		BackResUid: utils.StringPtr(bri.BackResUid),
		Name:       utils.StringPtr(bri.Name),
		Type:       utils.StringPtr(bri.Type),
		Status:     utils.StringPtr(bri.Status),
		Reason:     utils.StringPtr(bri.Reason),
		Creator:    utils.StringPtr(user),
		Created:    utils.TimePtr(at),
	}
}

// NewBackResInfo - 지정된 클러스터 정보로 백업/복원 정보를 생성한다.
func NewBackResInfo(cloudId string, clusterId string, name string, isBackup bool) *BackResInfo {
	backresInfo := &BackResInfo{
		CloudUid:   cloudId,
		ClusterUid: clusterId,
		Name:       name,
	}

	if isBackup {
		backresInfo.Type = "B"
	} else {
		backresInfo.Type = "C"
	}

	backresInfo.NewKey()

	return backresInfo
}
