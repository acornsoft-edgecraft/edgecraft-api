package postgresdb

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
	"github.com/gofrs/uuid"
)

const getAllCodeSQL = `
SELECT *
FROM tbl_code c 
`

const searchCodeSQL = `
SELECT *
FROM tbl_code tc
WHERE 1=1
{{- if ne (isNull .CodeUid) false }}
 AND code_uid = :code_uid
{{- end}}
{{- if ne (isNull .CodeGroupUID) false }}
 AND code_group_uid = :code_group_uid
{{- end}}
{{- if ne (isNull .CodeID) false }}
 AND code_id = :code_id
{{- end}}
{{- if ne (isNull .CodeName) false }}
 AND code_name = :code_name
{{- end}}
{{- if ne (isNull .CodeDescription) false }}
 AND code_description = :code_description
{{- end}}
{{- if ne (isNull .CodeDisplayOrder) false }}
 AND code_display_order = :code_display_order
{{- end}}
{{- if ne (isNull  .UseYn) false}}
 AND use_yn = :use_yn
{{- end}}
{{- if ne (isNull  .Creator) false}}
 AND creator = :creator
{{- end}}
{{- if ne (isNull  .Updater) false}}
 AND updater = :updater
{{- end}}
 ORDER BY ode_display_order
`

const deleteCodeByGroupUIDSQL = `
DELETE
FROM tbl_code
WHERE code_uid = :code_uid
`

// RegisterCloud - Registration a new Cloud
func (db *DB) CreateCode(create *model.Code) error {
	return db.GetClient().Insert(create)
}

// GetCode - Returns a GetCode
func (db *DB) GetCode(uid uuid.UUID) (*model.Code, error) {
	obj, err := db.GetClient().Get(&model.Code{}, uid)
	if err != nil {
		return nil, err
	}
	if obj != nil {
		res := obj.(*model.Code)
		return res, nil
	}
	return nil, nil
}

// GetAllCode - Returns all Code list
func (db *DB) GetAllCode() ([]model.Code, error) {
	var res []model.Code
	_, err := db.GetClient().Select(&res, getAllCodeSQL)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateCloud - saves the given RegisterCloud struct
func (db *DB) UpdateCode(req *model.Code) (int64, error) {
	// Find and Update
	// fmt.Println("-- UpdateCloud --")
	// utils.Print(req)
	count, err := db.GetClient().Update(req)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// SelectCode - Returns a SelectCode
func (db *DB) SelectCode(uid uuid.UUID) ([]model.Code, error) {
	obj, err := db.GetClient().Get(&model.Code{}, uid)
	if err != nil {
		return nil, err
	}
	if obj != nil {
		res := obj.([]model.Code)
		return res, nil
	}
	return nil, nil
}

// DeleteCode - deletes the RegisterCloud with the given id
func (db *DB) DeleteCode(uid uuid.UUID) (int64, error) {
	count, err := db.GetClient().Delete(&model.Code{CodeUID: &uid})
	if err != nil {
		return -1, err
	}
	return count, nil
}

// DeleteCodeByGroupID - GroupUID 로 다중건 삭제
func (db *DB) DeleteCodeByGroupUID(uid uuid.UUID) (int64, error) {
	m := map[string]interface{}{"eventSetUid": uid}
	result, err := db.GetClient().Exec(deleteCodeByGroupUIDSQL, m)
	if err != nil {
		return -1, err
	}
	count, err := result.RowsAffected()
	return count, err
}

// SearchCode - Returns Search Code list
func (db *DB) SearchCode(u model.Code) ([]model.Code, error) {
	var res []model.Code
	sql, err := utils.RenderTmpl("queryTemplate", searchCodeSQL, u)
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
