package model

import "time"

type Cloud struct {
	CloudUID    *int       `json:"cloudtUid" db:"cloud_uid"`
	CloudName   *string    `json:"cloudName" db:"cloud_name"`
	CloudTpye   *int       `json:"cloudTpye" db:"cloud_type"`
	CloudDesc   *string    `json:"cloudDesc" db:"cloud_description"`
	CloudStatus *int       `json:"cloudStatus" db:"cloud_state"`
	Creator     *int       `json:"creator" db:"creator"`
	CreatedAt   *time.Time `json:"createdAt" db:"created_at"`
	Updater     *int       `json:"updater" db:"updater"`
	UpdatedAt   *time.Time `json:"updatedAt" db:"updated_at"`
	CompletedAt *time.Time `json:"completedAt" db:"completed_at"`
}
