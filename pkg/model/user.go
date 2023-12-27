package model

import (
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
	"github.com/gofrs/uuid"
)

type User struct {
	UserUID string `json:"user_uid" db:"user_uid"`
	Role    int    `json:"role" db:"role"`
	Name    string `json:"name" db:"name"`
	// ID      string    `json:"userId" db:"id"`
	Password string    `json:"password" db:"password"`
	Email    string    `json:"email" db:"email"`
	Status   int       `json:"status" db:"state"`
	Created  time.Time `json:"created" db:"created_at"`
}

func (u *User) NewKey() {
	if u.UserUID == "" {
		u.UserUID = uuid.Must(uuid.NewV4()).String()
	}
}

// ToTable - User 정보를 테이블 정보로 설정
func (u *User) ToTable(ut *UserTable, isUpdate bool, user string, at time.Time) {
	if isUpdate {
		ut.Updater = utils.StringPtr(user)
		ut.Updated = utils.TimePtr(at)
	} else {
		u.NewKey()
		ut.Password = utils.StringPtr(utils.GetHashStr(u.Password))
		ut.Creator = utils.StringPtr(user)
		ut.Created = utils.TimePtr(at)
	}
	ut.UserUID = utils.StringPtr(u.UserUID)
	ut.Role = utils.IntPrt(u.Role)
	ut.Name = utils.StringPtr(u.Name)
	ut.Email = utils.StringPtr(u.Email)
	ut.Status = utils.IntPrt(1)
}
