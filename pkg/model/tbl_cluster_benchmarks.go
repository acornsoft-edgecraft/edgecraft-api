/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import (
	"time"
)

// OpenstackBenchmarksTable - 클러스터 Benchmarks 테이블 정보 (Openstack)
type OpenstackBenchmarksTable struct {
	CloudUid        *string    `json:"cloud_uid" db:"cloud_uid"`
	ClusterUid      *string    `json:"cluster_uid" db:"cluster_uid"`
	BenchmarksUid   *string    `json:"benchmarks_uid" db:"benchmarks_uid"`
	CisVersion      *string    `json:"cis_version" db:"cis_version"`
	DetectedVersion *string    `json:"detected_version" db:"detected_version"`
	Results         *string    `json:"results" db:"results"`
	Totals          *string    `json:"totals" db:"totals"`
	SuccessYn       *bool      `json:"success_yn" db:"success_yn"`
	Reason          *string    `json:"reason" db:"reason"`
	Creator         *string    `json:"creator" db:"creator"`
	Created         *time.Time `json:"created_at" db:"created_at"`
}
