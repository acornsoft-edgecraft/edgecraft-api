package postgresdb

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
)

const getAllCloudSQL = `
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

// const searchCloudSQL = `
// SELECT
// 	cloud_uid,
// 	cloud_name,
// 	cloud_type,
// 	cloud_description,
// 	cloud_state,
// 	creator,
// 	created_at,
// 	updater,
// 	updated_at
// FROM tbl_cloud
// WHERE 1=1
// {{- if ne (pointerToUUID .CloudUID) ""}}
//  AND cloud_uid =:cloud_uid
// {{- end}}
// {{- if ne (pointerToUUID .CloudName) ""}}
//  AND cloud_name =:cloud_name
// {{- end}}
// {{- if ne (pointerToString .CloudType) ""}}
//  AND cloud_type =:cloud_type
// {{- end}}
// {{- if ne (pointerToString .CloudDesc) ""}}
//  AND cloud_description =:cloud_description
// {{- end}}
// {{- if ne (pointerToString .CloudStatus) ""}}
//  AND cloud_state =:cloud_state
// {{- end}}
// {{- if gt (pointerToInt .Creator) 0}}
//  AND creator =:creator
// {{- end}}
// {{- if ne (pointerToString .CreatedAt) ""}}
//  AND created_at =:created_at
// {{- end}}
// {{- if ne (pointerToString .Updater) ""}}
//  AND updater =:updater
// {{- end}}
// {{- if ne (pointerToString .UpdatedAt) ""}}
//  AND updated_at =:updated_at
// {{- end}}
//  ORDER BY cloud_name
// `

// GetAllCloud - Returns all Cloud list
func (db *DB) GetAllCloud() ([]model.CloudList, error) {
	var res []model.CloudList
	_, err := db.GetClient().Select(&res, getAllCloudSQL)
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

// // UpdateCloud - saves the given RegisterCloud struct
// func (db *DB) UpdateCloud(req *model.CloudTable) (int64, error) {
// 	// Find and Update
// 	count, err := db.GetClient().Update(req)
// 	if err != nil {
// 		return -1, err
// 	}
// 	return count, nil
// }

// // DeleteCloud - deletes the RegisterCloud with the given id
// func (db *DB) DeleteCloud(uid uuid.UUID) (int64, error) {
// 	count, err := db.GetClient().Delete(&model.CloudTable{CloudUID: &uid})
// 	if err != nil {
// 		return -1, err
// 	}
// 	return count, nil
// }

// // GetSearchCloud - Returns Search Cloud list
// func (db *DB) GetSearchCloud(u model.CloudTable) ([]model.CloudTable, error) {
// 	var res []model.CloudTable
// 	sql, err := utils.RenderTmpl("queryTemplate", searchCloudSQL, u)
// 	if err != nil {
// 		return nil, err
// 	}
// 	m, _ := utils.StructToMap(u)
// 	_, err = db.GetClient().Select(&res, sql, m)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return res, nil
// }
