/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import (
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
	"github.com/gofrs/uuid"
)

type BackResParam struct {
	BackResId string `json:"backres_uid"`
	Name      string `json:"name"`
}

type BackResResult struct {
	ClusterInfo OsClusterInfo `json:"cluster"`
	List        []BackResInfo `json:"list"`
}

type BackResInfo struct {
	CloudUid   string    `json:"cloud_uid" db:"cloud_uid"`
	ClusterUid string    `json:"cluster_uid" db:"cluster_uid"`
	BackResUid string    `json:"backres_uid" db:"backres_uid"`
	Name       string    `json:"name" db:"name"`
	Type       string    `json:"type" db:"type"`
	Status     string    `json:"status" db:"status"`
	Reason     string    `json:"reasen" db:"reaseon"`
	BackupName string    `json:"backup_name" db:"backup_name"` // Restore인 경우 사용할 Backup 명
	Created    time.Time `json:"created" db:"created_at"`
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
		BackupName: utils.StringPtr(bri.BackupName),
		Creator:    utils.StringPtr(user),
		Created:    utils.TimePtr(at),
	}
}

// FromTable - Table정보를 Backup/Restore Info로 전환
func (bri *BackResInfo) FromTable(tbl BackResTable) {
	bri.CloudUid = *tbl.CloudUid
	bri.ClusterUid = *tbl.ClusterUid
	bri.BackResUid = *tbl.BackResUid
	bri.Name = *tbl.Name
	bri.Type = *tbl.Type
	bri.Status = *tbl.Status
	bri.Reason = *tbl.Reason
	bri.BackupName = *tbl.BackupName
	bri.Created = *tbl.Created
}

// NewBackResInfo - 지정된 클러스터 정보로 백업/복원 정보를 생성한다.
func NewBackResInfo(cloudId, clusterId, name, backupName string, isBackup bool) *BackResInfo {
	backresInfo := &BackResInfo{
		CloudUid:   cloudId,
		ClusterUid: clusterId,
		Name:       name,
	}

	if backupName != "" {
		backresInfo.BackupName = backupName
	}

	if isBackup {
		backresInfo.Type = "B"
	} else {
		backresInfo.Type = "R"
	}

	backresInfo.NewKey()

	backresInfo.Status = "R"

	return backresInfo
}
