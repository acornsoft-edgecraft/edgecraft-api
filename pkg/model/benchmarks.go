/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import (
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
	"github.com/gofrs/uuid"
)

type OpenstackBenchmarksList struct {
	CloudUid      *string    `json:"cloud_uid" db:"cloud_uid"`
	ClusterUid    *string    `json:"cluster_uid" db:"cluster_uid"`
	BenchmarksUid *string    `json:"benchmarks_uid" db:"benchmarks_uid"`
	Totals        *Totals    `json:"totals" db:"totals"`
	Status        *int       `json:"status" db:"state"`
	Reason        *string    `json:"reason" db:"reason"`
	Created       *time.Time `json:"created" db:"created_at"`
}

type OpenstackBenchmarksSet struct {
	CloudUid        string    `json:"cloud_uid" db:"cloud_uid"`
	ClusterUid      string    `json:"cluster_uid" db:"cluster_uid"`
	BenchmarksUid   string    `json:"benchmarks_uid" db:"benchmarks_uid"`
	CisVersion      string    `json:"cis_version" db:"cis_version"`
	DetectedVersion string    `json:"detected_version" db:"detected_version"`
	Results         string    `json:"results" db:"results"`
	Totals          string    `json:"totals" db:"totals"`
	Status          int       `json:"status" db:"state"`
	Reason          string    `json:"reason" db:"reason"`
	Created         time.Time `json:"created" db:"created_at"`
}

func (ob *OpenstackBenchmarksSet) NewKey() {
	ob.BenchmarksUid = uuid.Must(uuid.NewV4()).String()
}

func (ob *OpenstackBenchmarksSet) ToTable(cloudId *string, clusterId *string, user string, at time.Time) (benchmarksTable *OpenstackBenchmarksTable) {
	benchmarksTable = &OpenstackBenchmarksTable{
		CloudUid:      cloudId,
		ClusterUid:    clusterId,
		BenchmarksUid: utils.StringPtr(ob.BenchmarksUid),
		Status:        utils.IntPrt(1),
		Creator:       utils.StringPtr(user),
		Created:       utils.TimePtr(at),
	}
	return
}
