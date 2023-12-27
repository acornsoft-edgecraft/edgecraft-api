package postgresdb

import "github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"

// GetUserList - 사용자 목록
func (db *DB) GetUserList() ([]*model.UserTable, error) {
	var res []*model.UserTable
	_, err := db.GetClient().Select(&res, getUserListSQL)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetUser - 사용자 조회
func (db *DB) GetUser(userId string) (*model.UserTable, error) {
	data, err := db.GetClient().Get(&model.UserTable{}, userId)
	if err != nil {
		return nil, err
	}
	if data != nil {
		return data.(*model.UserTable), nil
	}
	return nil, nil
}

// InsertUser - 사용자 등록
func (db *DB) InsertUser(user *model.UserTable) error {
	return db.GetClient().Insert(user)
}

// UpdateUser - 사용자 수정
func (db *DB) UpdateUser(ut *model.UserTable) (int64, error) {
	count, err := db.GetClient().Update(ut)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// DeleteUser - 사용자 삭제
func (db *DB) DeleteUser(userId string) (int64, error) {
	count, err := db.GetClient().Delete(&model.UserTable{UserUID: &userId})
	if err != nil {
		return -1, err
	}
	return count, nil
}
