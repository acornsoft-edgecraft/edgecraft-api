package postgresdb

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
	"github.com/gofrs/uuid"
)

const getAllCloudSQL = `
select 
	tc.cloud_uid,
	tc.cloud_name,
	tc.cloud_type,
	tc.cloud_description,
	tc.cloud_state,
	tc.creator,
	tc.created_at,
	tc.updater,
	tc.updated_at ,
	count(tcn.cloud_node_uid) as node_count
from edgecraft.edgecraft.tbl_cloud tc 
left join edgecraft.edgecraft.tbl_cloud_node tcn 
on tc.cloud_uid = tcn.cloud_uid
group by tc.cloud_uid 
`

const searchCloudSQL = `
SELECT
	cloud_uid,
	cloud_name,
	cloud_type,
	cloud_description,
	cloud_state,
	creator,
	created_at,
	updater,
	updated_at
FROM tbl_cloud
WHERE 1=1
{{- if ne (pointerToUUID .CloudUID) ""}}
 AND cloud_uid =:cloud_uid
{{- end}}
{{- if ne (pointerToUUID .CloudName) ""}}
 AND cloud_name =:cloud_name
{{- end}}
{{- if ne (pointerToString .CloudType) ""}}
 AND cloud_type =:cloud_type
{{- end}}
{{- if ne (pointerToString .CloudDesc) ""}}
 AND cloud_description =:cloud_description
{{- end}}
{{- if ne (pointerToString .CloudStatus) ""}}
 AND cloud_state =:cloud_state
{{- end}}
{{- if gt (pointerToInt .Creator) 0}}
 AND creator =:creator
{{- end}}
{{- if ne (pointerToString .CreatedAt) ""}}
 AND created_at =:created_at
{{- end}}
{{- if ne (pointerToString .Updater) ""}}
 AND updater =:updater
{{- end}}
{{- if ne (pointerToString .UpdatedAt) ""}}
 AND updated_at =:updated_at
{{- end}}
 ORDER BY cloud_name
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
func (db *DB) GetAllCloud() ([]model.ResCloud, error) {
	var res []model.ResCloud
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

// GetSearchCloud - Returns Search Cloud list
func (db *DB) GetSearchCloud(u model.Cloud) ([]model.Cloud, error) {
	var res []model.Cloud
	sql, err := utils.RenderTmpl("queryTemplate", searchCloudSQL, u)
	if err != nil {
		return nil, err
	}
	m, _ := utils.StructToMap(u)
	_, err = db.GetClient().Select(&res, sql, m)
	if err != nil {
		return nil, err
	}
	return res, nil
}
