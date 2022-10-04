package postgresdb

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
)

// GetClouds - Returns all Cloud list
func (db *DB) GetClouds() ([]model.CloudList, error) {
	var res []model.CloudList
	_, err := db.GetClient().Select(&res, getCloudsSQL)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetCloud - Returns a GetCloud
func (db *DB) GetCloud(cloudId string) (*model.CloudTable, error) {
	data, err := db.GetClient().Get(&model.CloudTable{}, cloudId)
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
func (db *DB) DeleteCloud(cloudId string) (int64, error) {
	count, err := db.GetClient().Delete(&model.CloudTable{CloudUID: &cloudId})
	if err != nil {
		return -1, err
	}
	return count, nil
}
