package postgresdb

import "github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"

const getAuthUserSQL = `
  SELECT
	user_uid,
	user_id,
	user_role,
	user_name,
	password,
	email,
	user_state,
	creator,
	created_at,
	updater,
	updated_at
FROM tbl_user tu 
WHERE email = $1
  AND inactive_yn = 'N'
  AND user_state = 'A'
`

// GetUserByEmail - 이메일 기준 사용자 정보 조회 (로그인용)
func (db *DB) GetUserByEmail(email string) (*model.User, error) {
	var res *model.User
	if err := db.GetClient().SelectOne(&res, getAuthUserSQL, email); err != nil {
		return nil, err
	}
	return res, nil
}
