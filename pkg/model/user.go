package model

import "time"

// User User
type User struct {
	UserUID       *int       `json:"userUid" db:"user_uid"`
	UserRoleID    *string    `json:"userRoleId" db:"user_role_id"`
	Name          *string    `json:"name" db:"name"`
	UserID        *string    `json:"userId" db:"user_id"`
	UserPassword  *string    `json:"userPassword" db:"user_password"`
	Email         *string    `json:"email" db:"email"`
	Mobile        *string    `json:"mobile" db:"mobile"`
	Status        *string    `json:"status" db:"status"`
	AlarmRecvType *string    `json:"alarmRecvType" db:"alarm_recv_type"`
	RegUserUID    *int       `json:"regUserUid" db:"reg_user_uid"`
	RegDate       *time.Time `json:"regDate" db:"reg_date"`
	EdtUserUID    *int       `json:"edtUserUid" db:"edt_user_uid"`
	EdtDate       *time.Time `json:"edtDate" db:"edt_date"`
	SearchName    *string    `json:"searchName" db:"-"`
}

// SearchUser SearchUser
type SearchUser struct {
	User
	UserRoleName      *string `json:"userRoleName" db:"user_role_name"`
	AlarmRecvTypeName *string `json:"alarmRecvTypeName" db:"alarm_recv_type_name"`
}

// UserLoginHis UserLoginHis
type UserLoginHis struct {
	UserLoginUID int64     `json:"userLoginUid,omitempty" db:"user_login_uid"`
	UserUID      int       `json:"userUid,omitempty" db:"user_uid"`
	LoginType    string    `json:"loginType,omitempty" db:"login_type"`
	UserID       string    `json:"userId,omitempty" db:"user_id"`
	UserAgent    string    `json:"userAgent,omitempty" db:"user_agent"`
	ConnAddr     string    `json:"connAddr,omitempty" db:"conn_addr"`
	SuccYn       string    `json:"succYn,omitempty" db:"succ_yn"`
	RegDate      time.Time `json:"regDate,omitempty" db:"reg_date"`
}

// UserRole UserRole
type UserRole struct {
	UserRoleID   *string    `json:"userRoleId,omitempty" db:"user_role_id"`
	UserRoleName *string    `json:"userRoleName,omitempty" db:"user_role_name"`
	UseYn        *string    `json:"useYn,omitempty" db:"use_yn"`
	MenuAuths    []MenuAuth `json:"menuAuths,omitempty" db:"-"`
	RegUserUID   *int       `json:"regserUid" db:"reg_user_uid"`
	RegDate      *time.Time `json:"regDate" db:"reg_date"`
	EdtUserUID   *int       `json:"edtUserUid" db:"edt_user_uid"`
	EdtDate      *time.Time `json:"edtDate" db:"edt_date"`
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
	RegUserUID *int       `json:"regUserUid,omitempty" db:"reg_user_uid"`
	RegDate    *time.Time `json:"regDate,omitempty" db:"reg_date"`
	EdtUserUID *int       `json:"edtUserUid,omitempty" db:"edt_user_uid"`
	EdtDate    *time.Time `json:"edtDate,omitempty" db:"edt_date"`
	MenuDispYn *string    `json:"menuDispYn,omitempty" db:"-"`
}
