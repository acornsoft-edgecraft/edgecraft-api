package postgresdb

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
)

const getAllCloudSQL = `
SELECT *
FROM tbl_cloud c 
`

// GetAllConnect - Returns all Connect list
func (db *DB) GetAllCloud() ([]model.Cloud, error) {
	var res []model.Cloud
	_, err := db.GetClient().Select(&res, getAllCloudSQL)
	if err != nil {
		return nil, err
	}
	return res, nil
}
