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

// UserUID       *int       `json:"userUid" db:"user_uid"`
// 	UserRoleID    *string    `json:"userRoleId" db:"user_role_id"`
// 	Name          *string    `json:"name" db:"name"`
// 	UserID        *string    `json:"userId" db:"user_id"`
// 	UserPassword  *string    `json:"userPassword" db:"user_password"`
// 	Email         *string    `json:"email" db:"email"`
// 	Mobile        *string    `json:"mobile" db:"mobile"`
// 	Status        *string    `json:"status" db:"status"`
// 	AlarmRecvType *string    `json:"alarmRecvType" db:"alarm_recv_type"`
// 	RegUserUID    *int       `json:"regUserUid" db:"reg_user_uid"`
// 	RegDate       *time.Time `json:"regDate" db:"reg_date"`
// 	EdtUserUID    *int       `json:"edtUserUid" db:"edt_user_uid"`
// 	EdtDate       *time.Time `json:"edtDate" db:"edt_date"`
// 	SearchName    *string    `json:"searchName" db:"-"`

// 	type UserRole struct {
// 		UserRoleID   *string    `json:"userRoleId,omitempty" db:"user_role_id"`
// 		UserRoleName *string    `json:"userRoleName,omitempty" db:"user_role_name"`
// 		UseYn        *string    `json:"useYn,omitempty" db:"use_yn"`
// 		MenuAuths    []MenuAuth `json:"menuAuths,omitempty" db:"-"`
// 		RegUserUID   *int       `json:"regserUid" db:"reg_user_uid"`
// 		RegDate      *time.Time `json:"regDate" db:"reg_date"`
// 		EdtUserUID   *int       `json:"edtUserUid" db:"edt_user_uid"`
// 		EdtDate      *time.Time `json:"edtDate" db:"edt_date"`
// 	}

// 	type SearchUserRole struct {
// 		UserRole
// 		MenuRwCnt *int `json:"menuRwCnt,omitempty" db:"menu_rw_cnt"`
// 		MenuRoCnt *int `json:"menuRoCnt,omitempty" db:"menu_ro_cnt"`
// 	}

// 	// MenuAuth MenuAuth
// 	type MenuAuth struct {
// 		AuthUID    *int       `json:"authUid,omitempty" db:"auth_uid"`
// 		MenuID     *string    `json:"menuId,omitempty" db:"menu_id"`
// 		UserRoleID *string    `json:"userRoleId,omitempty" db:"user_role_id"`
// 		AttrRw     *string    `json:"attrRw,omitempty" db:"attr_rw"`
// 		UseYn      *string    `json:"useYn,omitempty" db:"use_yn"`
// 		RegUserUID *int       `json:"regUserUid,omitempty" db:"reg_user_uid"`
// 		RegDate    *time.Time `json:"regDate,omitempty" db:"reg_date"`
// 		EdtUserUID *int       `json:"edtUserUid,omitempty" db:"edt_user_uid"`
// 		EdtDate    *time.Time `json:"edtDate,omitempty" db:"edt_date"`
// 		MenuDispYn *string    `json:"menuDispYn,omitempty" db:"-"`
// 	}
