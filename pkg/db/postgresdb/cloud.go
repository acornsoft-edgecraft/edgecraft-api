package postgresdb

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
)

const getCloudListSQL = `
SELECT
	A.cloud_uid,
	A.name,
	A.type,
	A.state,
	A.description,
	A.created_at,
	COALESCE(B.k8s_version, '') as version,
	(SELECT COUNT(node_uid) FROM tbl_cloud_node WHERE cloud_uid = A.cloud_uid) AS nodeCount
FROM
	"edgecraft"."tbl_cloud" A
	LEFT JOIN "edgecraft"."tbl_cloud_cluster" B ON (B.cloud_uid = A.cloud_uid)
`

// GetCloudList - Returns all Cloud list
func (db *DB) GetCloudList() ([]model.CloudList, error) {
	var res []model.CloudList
	_, err := db.GetClient().Select(&res, getCloudListSQL)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetCloud - Returns a GetCloud
func (db *DB) GetCloud(cloudUid string) (*model.CloudTable, error) {
	data, err := db.GetClient().Get(&model.CloudTable{}, cloudUid)
	if err != nil {
		return nil, err
	}
	if data != nil {
		return data.(*model.CloudTable), nil
	}
	return nil, nil
}

// InsertCloud - Insert a new Cloud
func (db *DB) InsertCloud(cloud *model.CloudTable) error {
	return db.GetClient().Insert(cloud)
}

// UpdateCloud - Update a cloud
func (db *DB) UpdateCloud(cloud *model.CloudTable) (int64, error) {
	count, err := db.GetClient().Update(cloud)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// DeleteCloud - Delete a cloud
func (db *DB) DeleteCloud(cloudUid string) (int64, error) {
	count, err := db.GetClient().Delete(&model.CloudTable{CloudUID: &cloudUid})
	if err != nil {
		return -1, err
	}
	return count, nil
}
