package model

import (
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
)

type CodeGroup struct {
	GroupID     string `json:"group_id" example:"TestGroup"`
	Description string `json:"desc" example:"Code Group Testing"`
	UseYn       bool   `json:"use_yn" example:"true"`
}

// ToTable - CodeGroup 정보를 테이블 정보로 설정
func (cg *CodeGroup) ToTable(cgt *CodeGroupTable, isUpdate bool, user string, at time.Time) {
	cgt.GroupID = utils.StringPtr(cg.GroupID)
	cgt.Description = utils.StringPtr(cg.Description)
	cgt.UseYn = utils.BoolPtr(cg.UseYn)
	if isUpdate {
		cgt.Updater = utils.StringPtr(user)
		cgt.Updated = utils.TimePtr(at)
	} else {
		cgt.Creator = utils.StringPtr(user)
		cgt.Created = utils.TimePtr(at)
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
	ct.GroupID = utils.StringPtr(c.GroupID)
	ct.Code = utils.IntPrt(c.Code)
	ct.Name = utils.StringPtr(c.Name)
	ct.DisplayOrder = utils.IntPrt(c.DisplayOrder)
	ct.Description = utils.StringPtr(c.Description)
	ct.UseYn = utils.BoolPtr(c.UseYn)
	if isUpdate {
		ct.Updater = utils.StringPtr(user)
		ct.Updated = utils.TimePtr(at)
	} else {
		ct.Creator = utils.StringPtr(user)
		ct.Created = utils.TimePtr(at)
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
