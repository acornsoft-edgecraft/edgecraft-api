package postgresdb

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
	"github.com/gofrs/uuid"
)

const getAllCodeGroupSQL = `
SELECT *
FROM tbl_code_group tcg 
`

const selectCodeGroupSQL = `
SELECT *
FROM tbl_code_group c
WHERE
group_code_uid = $1
`
const searchCodeGroupSQL = `
SELECT *
FROM tbl_code_group c
WHERE 1=1
{{- if ne (isNull .CodeGroupUID) false }}
AND code_group_uid = :code_group_uid
{{- end }}
{{- if ne (isNull .CodeGroupName) false }}
AND code_group_name = :code_group_name
{{- end }}
{{- if ne (isNull .CodeGroupDescription) false }}
AND code_group_description =:code_group_description
{{- end }}
{{- if ne (isNull .UseYn) false }}
AND use_yn =:use_yn
{{- end }}
{{- if ne (isNull .Creator) false }}
AND creator =:creator
{{- end }}
{{- if ne (isNull .Updater) false }}
AND updater =:updater
{{- end }}
`

// RegisterCloud - Registration a new Cloud
func (db *DB) CreateCodeGroup(create *model.CodeGroup) error {
	return db.GetClient().Insert(create)
}

// GetCode - Returns a GetCode
func (db *DB) GetCodeGroup(uid uuid.UUID) (*model.CodeGroup, error) {
	obj, err := db.GetClient().Get(&model.CodeGroup{}, uid)
	if err != nil {
		return nil, err
	}
	if obj != nil {
		res := obj.(*model.CodeGroup)
		return res, nil
	}
	return nil, nil
}

// GetAllCode - Returns all Code list
func (db *DB) GetAllCodeGroup() ([]model.CodeGroup, error) {
	var res []model.CodeGroup
	_, err := db.GetClient().Select(&res, getAllCodeGroupSQL)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateCodeGroup - saves the given CodeGroup struct
func (db *DB) UpdateCodeGroup(req *model.CodeGroup) (int64, error) {
	// Find and Update
	count, err := db.GetClient().Update(req)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// SelectCloudNode - Returns a matching value for cloud clusters
func (db *DB) SelectCodeGroup(uid uuid.UUID) ([]model.CodeGroup, error) {
	var res []model.CodeGroup
	_, err := db.GetClient().Select(&res, selectCodeGroupSQL, uid)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// DeleteCode - deletes the RegisterCloud with the given id
func (db *DB) DeleteCodeGroup(uid uuid.UUID) (int64, error) {
	count, err := db.GetClient().Delete(&model.CodeGroup{CodeGroupUID: &uid})
	if err != nil {
		return -1, err
	}
	return count, nil
}

// SearchCodeGroup - Returns Search CodeGroup list
func (db *DB) SearchCodeGroup(u model.CodeGroup) ([]model.CodeGroup, error) {
	var res []model.CodeGroup
	sql, err := utils.RenderTmpl("queryTemplate", searchCodeGroupSQL, u)
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
