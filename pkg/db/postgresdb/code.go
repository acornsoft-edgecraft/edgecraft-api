package postgresdb

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
)

// GetCodeGroupList - 전체 코드 그룹 조회
func (db *DB) GetCodeGroupList() ([]*model.CodeGroupTable, error) {
	var list []*model.CodeGroupTable
	_, err := db.GetClient().Select(&list, getCodeGroupListSQL)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// GetCodeGroup - 코드 그룹 상세 조회
func (db *DB) GetCodeGroup(groupId string) (*model.CodeGroupTable, error) {
	data, err := db.GetClient().Get(&model.CodeGroupTable{}, groupId)
	if err != nil {
		return nil, err
	}
	if data != nil {
		return data.(*model.CodeGroupTable), nil
	}
	return nil, nil
}

// InsertCodeGroup - 코드 그룹 등록
func (db *DB) InsertCodeGroup(cgt *model.CodeGroupTable) error {
	return db.GetClient().Insert(cgt)
}

// UpdateCodeGroup - 코드 그룹 갱신
func (db *DB) UpdateCodeGroup(cgt *model.CodeGroupTable) (int64, error) {
	count, err := db.GetClient().Update(cgt)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// DeleteCodeGroup - 코드 그룹 삭제
func (db *DB) DeleteCodeGroup(groupId string) (int64, error) {
	count, err := db.GetClient().Delete(&model.CodeGroupTable{GroupID: &groupId})
	if err != nil {
		return -1, err
	}
	return count, nil
}

// GetCodeList - 모든 코드 조회
func (db *DB) GetCodeList() ([]*model.CodeTable, error) {
	var list []*model.CodeTable
	_, err := db.GetClient().Select(&list, getCodeListSQL)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// GetCodeListByGroup - 지정한 코드 그룹의 모든 코드 조회
func (db *DB) GetCodeListByGroup(groupId string) ([]*model.CodeTable, error) {
	var list []*model.CodeTable
	_, err := db.GetClient().Select(&list, getCodeByGroupSQL, groupId)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// GetCode - 코드 상세 조회
func (db *DB) GetCode(groupId string, code int) (*model.CodeTable, error) {
	data, err := db.GetClient().Get(&model.CodeTable{}, groupId, code)
	if err != nil {
		return nil, err
	}
	if data != nil {
		return data.(*model.CodeTable), nil
	}
	return nil, nil
}

// InsertCode - 코드 등록
func (db *DB) InsertCode(ct *model.CodeTable) error {
	return db.GetClient().Insert(ct)
}

// UpdateCode - 코드 갱신
func (db *DB) UpdateCode(ct *model.CodeTable) (int64, error) {
	count, err := db.GetClient().Update(ct)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// DeleteCode - 지정한 코드 삭제
func (db *DB) DeleteCode(groupId string, code int) (int64, error) {
	count, err := db.GetClient().Delete(&model.CodeTable{GroupID: &groupId, Code: &code})
	if err != nil {
		return -1, err
	}
	return count, nil
}

// DeleteCodeByGroup - 지정한 코드 그룹의 모든 코드 삭제
func (db *DB) DeleteCodeByGroup(groupId string) (int64, error) {
	count, err := db.GetClient().Delete(&model.CodeTable{GroupID: &groupId})
	if err != nil {
		return -1, err
	}
	return count, nil
	// result, err := db.GetClient().Exec(deleteCodeByGroupSQL, groupId)
	// if err != nil {
	// 	return -1, err
	// }
	// return result.RowsAffected()
}
