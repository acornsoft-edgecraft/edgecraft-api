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
func (db *DB) CreateRegisterCloud(create *model.Cloud) error {
	return db.GetClient().Insert(create)
}

// GetUserRole - Returns a UserRole
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

// UpdateRegisterCloud - saves the given RegisterCloud struct
func (db *DB) UpdateRegisterCloud(req *model.Cloud) (int64, error) {
	// Find and Update
	count, err := db.GetClient().Update(req)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// DeleteRegisterCloud - deletes the RegisterCloud with the given id
func (db *DB) DeleteRegisterCloud(uid uuid.UUID) (int64, error) {
	count, err := db.GetClient().Delete(&model.Cloud{CloudUID: &uid})
	if err != nil {
		return -1, err
	}
	return count, nil
}
