package model

import (
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
	"github.com/gofrs/uuid"
)

// CloudSet - Data set of Cluoud
type CloudSet struct {
	Cloud       *CloudInfo       `json:"cloud"`
	Cluster     *ClusterInfo     `json:"cluster"`
	Nodes       *NodesInfo       `json:"nodes"`
	EtcdStorage *EtcdStorageInfo `json:"etcd_storage"`
	OpenStack   *OpenstackInfo   `json:"openstack"`
}

// ToTable - CloudInfo를 대상 Table 정보에 매핑 처리
func (ci *CloudSet) ToTable(isUpdate bool, user string, at time.Time) (cloudTable *CloudTable, clusterTable *ClusterTable, nodeTables []*NodeTable) {
	// Cloud Table
	cloudTable = &CloudTable{}
	ci.Cloud.ToTable(cloudTable, isUpdate, user, at)

	// Cluster Table
	clusterTable = &ClusterTable{}
	clusterTable.CloudUid = cloudTable.CloudUID
	ci.Cluster.ToTable(clusterTable, isUpdate, user, at)
	ci.EtcdStorage.ToTable(clusterTable)

	// Node는 Delete & Insert 방식이므로 Update 개념 없음.
	nodeTables = ci.Nodes.ToTable(clusterTable, false, user, at)
	for _, node := range nodeTables {
		node.CloudUid = cloudTable.CloudUID
		node.ClusterUid = clusterTable.ClusterUid
	}

	// TODO: Openstack

	return
}

// CloudInfo - Data of Cloud
type CloudInfo struct {
	CloudUID string `json:"cloud_uid" db:"-" example:""`
	Name     string `json:"name" example:"test cloud"`
	Type     string `json:"type" example:"1"`
	Desc     string `json:"desc" example:"Baremtal cloud"`
}

// NewKey - Make new UUID V4
func (ci *CloudInfo) NewKey() {
	if ci.CloudUID == "" {
		ci.CloudUID = uuid.Must(uuid.NewV4()).String()
	}
}

// ToTable - CloudInfo정보를 Table 정보로 설정
func (ci *CloudInfo) ToTable(cloudTable *CloudTable, isUpdate bool, user string, at time.Time) {
	if isUpdate {
		*cloudTable.Updater = user
		*cloudTable.Updated = at
	} else {
		ci.NewKey()
		*cloudTable.Creator = user
		*cloudTable.Created = at
	}

	utils.CopyTo(&cloudTable, ci)
	*cloudTable.Status = "1"
}

// FromTable - 테이블의 정보를 CluoudInfo 정보로 설정
func (ci *CloudInfo) FromTable(cloudTable *CloudTable) {
	ci.CloudUID = *cloudTable.CloudUID
	ci.Name = *cloudTable.Name
	ci.Type = *cloudTable.Type
	ci.Desc = *cloudTable.Desc
}

// CloudList - List of Clouds
type CloudList struct {
	CloudUID  string    `json:"cloud_uid" db:"cloud_uid"`
	Name      string    `json:"name" db:"name"`
	Type      string    `json:"type" db:"type"`
	Desc      string    `json:"desc" db:"description"`
	Status    string    `json:"status" db:"state"`
	NodeCount int       `json:"nodeCount" db:"nodecount"`
	Version   string    `json:"version" db:"version"`
	Created   time.Time `json:"created" db:"created_at"`
}
