package postgresdb

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	"github.com/gofrs/uuid"
)

const getAllCloudClusterSQL = `
SELECT *
FROM tbl_cloud_cluster c 
`

// RegisterCloud - Registration a new Cloud
func (db *DB) CreateCloudCluster(create *model.CloudCluster) error {
	return db.GetClient().Insert(create)
}

// GetUserRole - Returns a UserRole
func (db *DB) GetCloudCluster(uid uuid.UUID, i string) (*model.CloudCluster, error) {

	obj, err := db.GetClient().Get(&model.CloudCluster{}, nil, uid)
	if err != nil {
		return nil, err
	}
	if obj != nil {
		res := obj.(*model.CloudCluster)
		return res, nil
	}
	return nil, nil
}

// GetAllCloud - Returns all Cloud list
func (db *DB) GetAllCloudCluster() ([]model.CloudCluster, error) {
	var res []model.CloudCluster
	_, err := db.GetClient().Select(&res, getAllCloudClusterSQL)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateCloud - saves the given RegisterCloud struct
func (db *DB) UpdateCloudCluster(req *model.CloudCluster) (int64, error) {
	// Find and Update
	count, err := db.GetClient().Update(req)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// DeleteCloud - deletes the RegisterCloud with the given id
func (db *DB) DeleteCloudCluster(uid uuid.UUID) (int64, error) {
	count, err := db.GetClient().Delete(&model.CloudCluster{CloudClusterUid: &uid})
	if err != nil {
		return -1, err
	}
	return count, nil
}
