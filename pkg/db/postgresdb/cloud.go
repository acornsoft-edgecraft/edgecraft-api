package postgresdb

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	"github.com/gofrs/uuid"
)

const getAllCloudSQL = `
SELECT *
FROM tbl_cloud c 
`

// RegisterCloud - Registration a new Cloud
func (db *DB) CreateCloud(create *model.Cloud) error {
	return db.GetClient().Insert(create)
}

// GetCloud - Returns a GetCloud
func (db *DB) GetCloud(uid uuid.UUID) (*model.Cloud, error) {
	obj, err := db.GetClient().Get(&model.Cloud{}, uid)
	if err != nil {
		return nil, err
	}
	if obj != nil {
		res := obj.(*model.Cloud)
		return res, nil
	}
	return nil, nil
}

// GetAllCloud - Returns all Cloud list
func (db *DB) GetAllCloud() ([]model.Cloud, error) {
	var res []model.Cloud
	_, err := db.GetClient().Select(&res, getAllCloudSQL)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateCloud - saves the given RegisterCloud struct
func (db *DB) UpdateCloud(req *model.Cloud) (int64, error) {
	// Find and Update
	count, err := db.GetClient().Update(req)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// DeleteCloud - deletes the RegisterCloud with the given id
func (db *DB) DeleteCloud(uid uuid.UUID) (int64, error) {
	count, err := db.GetClient().Delete(&model.Cloud{CloudUID: &uid})
	if err != nil {
		return -1, err
	}
	return count, nil
}
