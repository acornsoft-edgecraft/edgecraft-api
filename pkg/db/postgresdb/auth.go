package postgresdb

import "github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"

// GetUserByEmail - 이메일 기준 사용자 정보 조회 (로그인용)
func (db *DB) GetUserByEmail(email string) (*model.UserTable, error) {
	var res *model.UserTable
	if err := db.GetClient().SelectOne(&res, getAuthUserSQL, email); err != nil {
		return nil, err
	}
	return res, nil
}
