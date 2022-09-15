package model

import (
	"strings"
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
	"github.com/gofrs/uuid"
)

// User User
type User struct {
	UserUID  *uuid.UUID `json:"userUid" db:"user_uid, default:uuid_generate_v4()"`
	UserRole *string    `json:"userRole" db:"user_role"`
	Name     *string    `json:"name" db:"user_name"`
	UserID   *string    `json:"userId" db:"user_id"`
	Password *string    `json:"password" db:"password"`
	Email    *string    `json:"email" db:"email"`
	Status   *string    `json:"status" db:"user_state"`
	Creator  *string    `json:"creator" db:"creator"`
	Created  *time.Time `json:"created" db:"created_at"`
	Updater  *string    `json:"updater" db:"updater"`
	Updated  *time.Time `json:"updated" db:"updated_at"`
}

// SearchUser SearchUser
type SearchUser struct {
	User
	UserRoleName      *string `json:"userRoleName" db:"user_role_name"`
	AlarmRecvTypeName *string `json:"alarmRecvTypeName" db:"alarm_recv_type_name"`
}

// UserLoginHis UserLoginHis
type UserLoginHis struct {
	UserLoginUID int64      `json:"userLoginUid,omitempty" db:"user_login_uid"`
	UserUID      *uuid.UUID `json:"userUid" db:"user_uid, default:uuid_generate_v4()"`
	LoginType    string     `json:"loginType,omitempty" db:"login_type"`
	UserID       string     `json:"userId,omitempty" db:"user_id"`
	UserAgent    string     `json:"userAgent,omitempty" db:"user_agent"`
	ConnAddr     string     `json:"connAddr,omitempty" db:"conn_addr"`
	SuccYn       string     `json:"succYn,omitempty" db:"succ_yn"`
	RegDate      time.Time  `json:"regDate,omitempty" db:"reg_date"`
}

// UserRole UserRole
type UserRole struct {
	UserRoleID   *string    `json:"userRoleId,omitempty" db:"user_role_id"`
	UserRoleName *string    `json:"userRoleName,omitempty" db:"user_role_name"`
	UseYn        *string    `json:"useYn,omitempty" db:"use_yn"`
	MenuAuths    []MenuAuth `json:"menuAuths,omitempty" db:"-"`
	Creator      *string    `json:"creator" db:"creator"`
	Created      *time.Time `json:"created" db:"created_at"`
	Updater      *string    `json:"updater" db:"updater"`
	Updated      *time.Time `json:"updated" db:"updated_at"`
}

type SearchUserRole struct {
	UserRole
	MenuRwCnt *int `json:"menuRwCnt,omitempty" db:"menu_rw_cnt"`
	MenuRoCnt *int `json:"menuRoCnt,omitempty" db:"menu_ro_cnt"`
}

// MenuAuth MenuAuth
type MenuAuth struct {
	AuthUID    *int       `json:"authUid,omitempty" db:"auth_uid"`
	MenuID     *string    `json:"menuId,omitempty" db:"menu_id"`
	UserRoleID *string    `json:"userRoleId,omitempty" db:"user_role_id"`
	AttrRw     *string    `json:"attrRw,omitempty" db:"attr_rw"`
	UseYn      *string    `json:"useYn,omitempty" db:"use_yn"`
	Creator    *string    `json:"creator" db:"creator"`
	Created    *time.Time `json:"created" db:"created_at"`
	Updater    *string    `json:"updater" db:"updater"`
	Updated    *time.Time `json:"updated" db:"updated_at"`
	MenuDispYn *string    `json:"menuDispYn,omitempty" db:"-"`
}

// MatchPassword - 비밀번호 매칭 검증
func (u *User) MatchPassword(password string) (int, bool) {
	// Hash 검증
	result := strings.Compare(utils.GetHashStr(password), *u.Password)
	if result != 0 {
		return common.CodeInvalidUser, false
	}

	return common.CodeOK, true
}
