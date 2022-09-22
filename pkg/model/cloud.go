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
func (ci *CloudSet) ToTable() (cloudTable *CloudTable, clusterTable *ClusterTable, nodeTables []*NodeTable) {
	var now = time.Now()
	var creator = "system"

	// Cloud Table
	cloudTable = &CloudTable{}
	ci.Cloud.ToTable(cloudTable)
	cloudTable.Creator = creator
	cloudTable.Created = now

	// Cluster Table
	clusterTable = &ClusterTable{}
	ci.Cluster.ToTable(clusterTable)
	ci.EtcdStorage.ToTable(clusterTable)
	clusterTable.CloudUid = cloudTable.CloudUID
	clusterTable.Creator = creator
	clusterTable.Created = now

	nodeTables = ci.Nodes.ToTable(clusterTable)
	for _, node := range nodeTables {
		node.CloudUid = cloudTable.CloudUID
		node.ClusterUid = clusterTable.ClusterUid
		node.Creator = creator
		node.Created = now
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
func (ci *CloudInfo) ToTable(cloudTable *CloudTable) {
	ci.NewKey()
	utils.CopyTo(&cloudTable, ci)
	cloudTable.Status = "1"
}

// FromTable - 테이블의 정보를 CluoudInfo 정보로 설정
func (ci *CloudInfo) FromTable(cloudTable *CloudTable) {
	ci.CloudUID = cloudTable.CloudUID
	ci.Name = cloudTable.Name
	ci.Type = cloudTable.Type
	ci.Desc = cloudTable.Desc
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

// used nullstring type
// type Cloud struct {
// 	CloudUID    uuid.UUID        `json:"cloudtUid" db:"cloud_uid, default:uuid_generate_v4()"`
// 	CloudName   utils.NullString `json:"cloudName" db:"cloud_name"`
// 	CloudType   utils.NullString `json:"cloudTpye" db:"cloud_type"`
// 	CloudDesc   utils.NullString `json:"cloudDesc" db:"cloud_description"`
// 	CloudStatus utils.NullString `json:"cloudStatus" db:"cloud_state"`
// 	Creator     utils.NullString `json:"creator" db:"creator"`
// 	CreatedAt   utils.NullTime   `json:"createdAt" db:"created_at"`
// 	Updater     utils.NullString `json:"updater" db:"updater"`
// 	UpdatedAt   utils.NullTime   `json:"updatedAt" db:"updated_at"`
// }

// type Cloud struct {
// 	CloudUID    uuid.UUID `json:"cloudtUid" db:"cloud_uid, default:uuid_generate_v4()"`
// 	CloudName   MyString  `json:"cloudName" db:"cloud_name"`
// 	CloudType   MyString  `json:"cloudTpye" db:"cloud_type"`
// 	CloudDesc   MyString  `json:"cloudDesc" db:"cloud_description"`
// 	CloudStatus MyString  `json:"cloudStatus" db:"cloud_state"`
// 	Creator     MyString  `json:"creator" db:"creator"`
// 	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
// 	Updater     MyString  `json:"updater" db:"updater"`
// 	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
// 	TestInt     int       `json:"testInt" db:"test_int"`
// }

// type MyString string

// const MyStringNull MyString = "\x00"

// // implements driver.Valuer, will be invoked automatically when written to the db
// func (s MyString) Value() (driver.Value, error) {
// 	if s == MyStringNull {
// 		return nil, nil
// 	}
// 	return []byte(s), nil
// }

// // implements sql.Scanner, will be invoked automatically when read from the db
// func (s *MyString) Scan(src interface{}) error {
// 	switch v := src.(type) {
// 	case string:
// 		*s = MyString(v)
// 	case []byte:
// 		*s = MyString(v)
// 	case nil:
// 		*s = MyStringNull
// 	}
// 	return nil
// }
