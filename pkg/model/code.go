package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type CodeGroup struct {
	CodeGroupUID         *uuid.UUID `json:"code_group_uid" db:"code_group_uid, primarykey, default:uuid_generate_v4()"`
	CodeGroupName        *string    `json:"code_group_name" db:"code_group_name"`
	CodeGroupDescription *string    `json:"code_group_description" db:"code_group_description"`
	UseYn                *bool      `json:"use_yn" db:"use_yn"`
	Creator              *string    `json:"creator" db:"creator"`
	CreatedAt            *time.Time `json:"created_at" db:"created_at"`
	Updater              *string    `json:"updater" db:"updater"`
	UpdatedAt            *time.Time `json:"updated_at" db:"updated_at"`
}

type Code struct {
	CodeUID          *uuid.UUID `json:"code_uid" db:"code_uid, primarykey, default:uuid_generate_v4()"`
	CodeGroupUID     *uuid.UUID `json:"code_group_uid" db:"code_group_uid"`
	CodeID           *string    `json:"code_id" db:"code_id"`
	CodeName         *string    `json:"code_name" db:"code_name"`
	CodeDescription  *string    `json:"code_description" db:"code_description"`
	CodeDisplayOrder *int       `json:"code_display_order" db:"code_display_order"`
	UseYn            *bool      `json:"use_yn" db:"use_yn"`
	Creator          *string    `json:"creator" db:"creator"`
	CreatedAt        *time.Time `json:"created_at" db:"created_at"`
	Updater          *string    `json:"updater" db:"updater"`
	UpdatedAt        *time.Time `json:"updated_at" db:"updated_at"`
}
