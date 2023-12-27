package model

import (
	"strings"
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
)

// User User
type UserTable struct {
	UserUID *string `json:"user_uid" db:"user_uid"`
	Role    *int    `json:"role" db:"role"`
	Name    *string `json:"name" db:"name"`
	// ID       *string    `json:"userId" db:"id"`
	Password *string    `json:"password" db:"password"`
	Email    *string    `json:"email" db:"email"`
	Status   *int       `json:"status" db:"state"`
	Creator  *string    `json:"creator" db:"creator"`
	Created  *time.Time `json:"created" db:"created_at"`
	Updater  *string    `json:"updater" db:"updater"`
	Updated  *time.Time `json:"updated" db:"updated_at"`
}

// MatchPassword - 비밀번호 매칭 검증
func (u *UserTable) MatchPassword(password string) (int, bool) {
	// Hash 검증
	result := strings.Compare(utils.GetHashStr(password), *u.Password)
	if result != 0 {
		return common.CodeInvalidUser, false
	}

	return common.CodeOK, true
}
