package postgresdb

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	"github.com/gofrs/uuid"
)

const getAllCloudNodeSQL = `
SELECT *
FROM tbl_cloud_node c 
`

// RegisterCloud - Registration a new Cloud
func (db *DB) CreateCloudNode(create *model.CloudNode) error {
	return db.GetClient().Insert(create)
}

// GetUserRole - Returns a UserRole
func (db *DB) GetCloudNode(uid uuid.UUID) (*model.CloudNode, error) {
	obj, err := db.GetClient().Get(&model.CloudNode{}, uid)
	if err != nil {
		return nil, err
	}
	if obj != nil {
		res := obj.(*model.CloudNode)
		return res, nil
	}
	return nil, nil
}

// GetAllCloud - Returns all Cloud list
func (db *DB) GetAllCloudNode() ([]model.CloudNode, error) {
	var res []model.CloudNode
	_, err := db.GetClient().Select(&res, getAllCloudNodeSQL)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateCloud - saves the given RegisterCloud struct
func (db *DB) UpdateCloudNode(req *model.CloudNode) (int64, error) {
	// Find and Update
	count, err := db.GetClient().Update(req)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// DeleteCloud - deletes the RegisterCloud with the given id
func (db *DB) DeleteCloudNode(uid uuid.UUID) (int64, error) {
	count, err := db.GetClient().Delete(&model.CloudNode{CloudNodeUid: &uid})
	if err != nil {
		return -1, err
	}
	return count, nil
}
