package model

import "github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"

type (
	// Login - 사용자 로그인 정보
	Login struct {
		Email    *string `json:"email" example:"ccambo@acornsoft.io"`
		Password *string `json:"password" example:"1234abcd@Acorn"`
	}

	LoginInfo struct {
		UserID   *string `json:"userId" db:"user_id"`
		UserRole *string `json:"userRole" db:"user_role"`
		Name     *string `json:"name" db:"user_name"`
		Email    *string `json:"email" db:"email"`
		Status   *string `json:"status" db:"user_state"`
		// TODO: Auth-Menu 연결 필요
	}

	// TokenResponse - 사용자
	TokenResponse struct {
		AccessToken  *string    `json:"accessToken"`
		RefreshToken *string    `json:"refreshToken"`
		User         *LoginInfo `json:"user"`
	}
)

// Validate - 로그인 정보 검증
func (l *Login) Validate() (int, bool) {
	if l.Email == nil || l.Password == nil {
		return common.CodeInvalidUser, false
	}
	return common.CodeOK, true
}
