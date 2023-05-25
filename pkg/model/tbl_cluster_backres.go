/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import (
	"time"
)

// BackResTable - 클러스터 백업/복원 테이블 정보 (Openstack)
type BackResTable struct {
	CloudUid   *string    `json:"cloud_uid" db:"cloud_uid"`
	ClusterUid *string    `json:"cluster_uid" db:"cluster_uid"`
	BackResUid *string    `json:"backres_uid" db:"backres_uid"`
	Name       *string    `json:"name" db:"name"`
	Type       *string    `json:"type" db:"type"`
	Status     *string    `json:"status" db:"status"` // 'R' - Running, 'C' - Completed, 'F' - Failed, 'P' - PartiallyFailed
	Reason     *string    `json:"reason" db:"reason"`
	BackupName *string    `json:"backup_name" db:"backup_name"` // Restore인 경우만 사용
	Creator    *string    `json:"creator" db:"creator"`
	Created    *time.Time `json:"created_at" db:"created_at"`
	Updater    *string    `json:"updater" db:"updater"`
	Updated    *time.Time `json:"updated_at" db:"updated_at"`
}
