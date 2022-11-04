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
	SaveOnly    bool             `json:"save_only"`
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

	// Node Table는 Delete & Insert 방식이므로 Update 개념 없음.
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
	CloudUID string `json:"cloud_uid" example:""`
	Name     string `json:"name" example:"test cloud"`
	Type     int    `json:"type" example:"1"`
	Desc     string `json:"desc" example:"Baremtal cloud"`
	Status   int    `json:"status" example:"1"`
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
		if cloudTable.CloudUID == nil {
			cloudTable.CloudUID = utils.StringPtr(ci.CloudUID)
		}
		cloudTable.Updater = utils.StringPtr(user)
		cloudTable.Updated = utils.TimePtr(at)
	} else {
		ci.NewKey()
		cloudTable.CloudUID = utils.StringPtr(ci.CloudUID)
		cloudTable.Creator = utils.StringPtr(user)
		cloudTable.Created = utils.TimePtr(at)
	}

	cloudTable.Name = utils.StringPtr(ci.Name)
	cloudTable.Type = utils.IntPrt(ci.Type)
	cloudTable.Desc = utils.StringPtr(ci.Desc)
	// NOTES: Status가 UI에 존재하지 않기 때문에 기본값으로 처리
	cloudTable.Status = utils.IntPrt(1)
}

// FromTable - 테이블의 정보를 CluoudInfo 정보로 설정
func (ci *CloudInfo) FromTable(cloudTable *CloudTable) {
	ci.CloudUID = *cloudTable.CloudUID
	ci.Name = *cloudTable.Name
	ci.Type = *cloudTable.Type
	ci.Desc = *cloudTable.Desc
	// NOTES: 조회용으로 상태 처리
	ci.Status = *cloudTable.Status
}

// CloudList - List of Clouds
type CloudList struct {
	CloudUID  string    `json:"cloud_uid" db:"cloud_uid"`
	Name      string    `json:"name" db:"name"`
	Type      int       `json:"type" db:"type"`
	Desc      string    `json:"desc" db:"description"`
	Status    int       `json:"status" db:"state"`
	NodeCount int       `json:"nodeCount"`
	Version   string    `json:"version"`
	Created   time.Time `json:"created" db:"created_at"`
}
