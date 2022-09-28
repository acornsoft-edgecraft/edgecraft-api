/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import "time"

type CodeGroupTable struct {
	GroupID     *string    `json:"group_id" db:"group_id"`
	Description *string    `json:"desc" db:"description"`
	UseYn       *bool      `json:"use_yn" db:"use_yn"`
	Creator     *string    `json:"creator" db:"creator"`
	Created     *time.Time `json:"created" db:"created_at"`
	Updater     *string    `json:"updater" db:"updater"`
	Updated     *time.Time `json:"updated" db:"updated_at"`
}

type CodeTable struct {
	GroupID      *string    `json:"group_id" db:"group_id"`
	Code         *int       `json:"code" db:"code"`
	Name         *string    `json:"name" db:"name"`
	DisplayOrder *int       `json:"display_order" db:"display_order"`
	Description  *string    `json:"desc" db:"description"`
	UseYn        *bool      `json:"use_yn" db:"use_yn"`
	Creator      *string    `json:"creator" db:"creator"`
	Created      *time.Time `json:"created" db:"created_at"`
	Updater      *string    `json:"updater" db:"updater"`
	Updated      *time.Time `json:"updated" db:"updated_at"`
}
