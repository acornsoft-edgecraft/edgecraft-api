package model

import (
	"time"

	"github.com/gofrs/uuid"
)

// used pointer
type Cloud struct {
	CloudUID    *uuid.UUID `json:"cloudUid" db:"cloud_uid, default:uuid_generate_v4()"`
	CloudName   *string    `json:"name" db:"cloud_name"`
	CloudType   *string    `json:"type" db:"cloud_type"`
	CloudDesc   *string    `json:"desc" db:"cloud_description"`
	CloudStatus *string    `json:"cloudStatus" db:"cloud_state"`
	Creator     *string    `json:"creator" db:"creator"`
	CreatedAt   *time.Time `json:"createdAt" db:"created_at"`
	Updater     *string    `json:"updater" db:"updater"`
	UpdatedAt   *time.Time `json:"updatedAt" db:"updated_at"`
}

// used nullstring type
// type Cloud struct {
// 	CloudUID    uuid.UUID        `json:"cloudtUid" db:"cloud_uid, default:uuid_generate_v4()"`
// 	CloudName   utils.NullString `json:"cloudName" db:"cloud_name"`
// 	CloudType   utils.NullString `json:"cloudTpye" db:"cloud_type"`
// 	CloudDesc   utils.NullString `json:"cloudDesc" db:"cloud_description"`
// 	CloudStatus utils.NullString `json:"cloudStatus" db:"cloud_state"`
// 	Creator     utils.NullString `json:"creator" db:"creator"`
// 	CreatedAt   utils.NullTime   `json:"createdAt" db:"created_at"`
// 	Updater     utils.NullString `json:"updater" db:"updater"`
// 	UpdatedAt   utils.NullTime   `json:"updatedAt" db:"updated_at"`
// }

// type Cloud struct {
// 	CloudUID    uuid.UUID `json:"cloudtUid" db:"cloud_uid, default:uuid_generate_v4()"`
// 	CloudName   MyString  `json:"cloudName" db:"cloud_name"`
// 	CloudType   MyString  `json:"cloudTpye" db:"cloud_type"`
// 	CloudDesc   MyString  `json:"cloudDesc" db:"cloud_description"`
// 	CloudStatus MyString  `json:"cloudStatus" db:"cloud_state"`
// 	Creator     MyString  `json:"creator" db:"creator"`
// 	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
// 	Updater     MyString  `json:"updater" db:"updater"`
// 	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
// 	TestInt     int       `json:"testInt" db:"test_int"`
// }

// type MyString string

// const MyStringNull MyString = "\x00"

// // implements driver.Valuer, will be invoked automatically when written to the db
// func (s MyString) Value() (driver.Value, error) {
// 	if s == MyStringNull {
// 		return nil, nil
// 	}
// 	return []byte(s), nil
// }

// // implements sql.Scanner, will be invoked automatically when read from the db
// func (s *MyString) Scan(src interface{}) error {
// 	switch v := src.(type) {
// 	case string:
// 		*s = MyString(v)
// 	case []byte:
// 		*s = MyString(v)
// 	case nil:
// 		*s = MyStringNull
// 	}
// 	return nil
// }
