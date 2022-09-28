package model

import "time"

type CodeGroup struct {
	GroupID     string `json:"group_id" example:"TestGroup"`
	Description string `json:"desc" example:"Code Group Testing"`
	UseYn       bool   `json:"use_yn" example:"true"`
}

// ToTable - CodeGroup 정보를 테이블 정보로 설정
func (cg *CodeGroup) ToTable(cgt *CodeGroupTable, isUpdate bool, user string, at time.Time) {
	cgt.GroupID = &cg.GroupID
	cgt.Description = &cg.Description
	cgt.UseYn = &cg.UseYn
	if isUpdate {
		cgt.Updater = &user
		cgt.Updated = &at
	} else {
		cgt.Creator = &user
		cgt.Created = &at
	}
}

// FromTable - 테이블 정보를 CodeGroup 정보로 설정
func (cg *CodeGroup) FromTable(cgt *CodeGroupTable) {
	cg.GroupID = *cgt.GroupID
	cg.Description = *cgt.Description
	cg.UseYn = *cgt.UseYn
}

type Code struct {
	GroupID      string `json:"group_id" example:"TestGroup"`
	Code         int    `json:"code" example:"1"`
	Name         string `json:"name" example:"TestCode #1"`
	DisplayOrder int    `json:"display_order" example:"0"`
	Description  string `json:"desc" example:"Test Code for Testing"`
	UseYn        bool   `json:"use_yn" example:"true"`
}

// ToTable - CodeGroup 정보를 테이블 정보로 설정
func (c *Code) ToTable(ct *CodeTable, isUpdate bool, user string, at time.Time) {
	ct.GroupID = &c.GroupID
	ct.Code = &c.Code
	ct.Name = &c.Name
	ct.DisplayOrder = &c.DisplayOrder
	ct.Description = &c.Description
	ct.UseYn = &c.UseYn
	if isUpdate {
		ct.Updater = &user
		ct.Updated = &at
	} else {
		ct.Creator = &user
		ct.Created = &at
	}
}

// FromTable - 테이블 정보를 CodeGroup 정보로 설정
func (c *Code) FromTable(ct *CodeTable) {
	c.GroupID = *ct.GroupID
	c.Code = *ct.Code
	c.Name = *ct.Name
	c.DisplayOrder = *ct.DisplayOrder
	c.Description = *ct.Description
	c.UseYn = *ct.UseYn
}
