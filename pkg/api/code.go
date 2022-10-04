/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package api

import (
	"strconv"
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/response"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	"github.com/labstack/echo/v4"
)

/************************
 * Code Group
 ************************/

// GetCodeGroupListHandler - 전체 코드그룹 리스트
// @Tags CodeGroup
// @Summary GetCodeGroupList
// @Description Get all code-group list
// @ID GetCodeGroupList
// @Produce json
// @Success 200 {object} response.ReturnData
// @Router /codes/groups [get]
func (a *API) GetCodeGroupListHandler(c echo.Context) error {
	list, err := a.Db.GetCodeGroupList()
	if err != nil {
		return response.Errorf(c, common.CodeFailedDatabase, err)
	}
	return response.Write(c, nil, list)
}

// GetCodeGroupHandler - 코드그룹 상세 조회
// @Tags CodeGroup
// @Summary GetCodeGroup
// @Description Get code group
// @ID GetCodeGroup
// @Produce json
// @Param groupId path string true "Code Group ID"
// @Success 200 {object} response.ReturnData
// @Router /codes/groups/{groupId} [get]
func (a *API) GetCodeGroupHandler(c echo.Context) error {
	groupId := c.Param("groupId")
	if groupId == "" {
		return response.ErrorfReqRes(c, groupId, common.CodeInvalidParm, nil)
	}

	data, err := a.Db.GetCodeGroup(groupId)
	if err != nil {
		return response.Errorf(c, common.CodeFailedDatabase, err)
	}
	if data == nil {
		return response.Errorf(c, common.DatabaseFalseData, err)
	}
	return response.Write(c, nil, data)
}

// SetCodeGroupHandler - 코드그룹 등록
// @Tags CodeGroup
// @Summary SetCodeGroup
// @Description Register code group
// @ID SetCodeGroup
// @Produce json
// @Param codeGroup body model.CodeGroup true "Code Group"
// @Success 200 {object} response.ReturnData
// @Router /codes/groups [post]
func (a *API) SetCodeGroupHandler(c echo.Context) error {
	// TODO: 로그인 사용자 정보 활용 방법은?
	var codeGroup model.CodeGroup

	err := getRequestData(c.Request(), &codeGroup)
	if err != nil {
		return response.ErrorfReqRes(c, codeGroup, common.CodeInvalidData, err)
	}

	var codeGroupTable *model.CodeGroupTable = &model.CodeGroupTable{}
	codeGroup.ToTable(codeGroupTable, false, "system", time.Now())

	// Start. Transaction 얻어옴
	txdb, err := a.Db.BeginTransaction()
	if err != nil {
		return response.ErrorfReqRes(c, codeGroup, common.CodeFailedDatabase, err)
	}

	// CodeGroup 등록
	err = txdb.InsertCodeGroup(codeGroupTable)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, codeGroup, common.CodeFailedDatabase, err)
	}

	txErr := txdb.Commit()
	if txErr != nil {
		logger.Info("DB commit Failed.", txErr)
	}

	return response.Write(c, codeGroup, nil)
}

// UpdateCodeGroupHandler - 코드그룹 수정
// @Tags CodeGroup
// @Summary UpdateCodeGroup
// @Description Update code group
// @ID UpdateCodeGroup
// @Produce json
// @Param codeGroup body model.CodeGroup true "Code Group"
// @Success 200 {object} response.ReturnData
// @Router /codes/groups [put]
func (a *API) UpdateCodeGroupHandler(c echo.Context) error {
	// TODO: 로그인 사용자 정보 활용 방법은?
	var codeGroup *model.CodeGroup

	err := getRequestData(c.Request(), &codeGroup)
	if err != nil {
		return response.ErrorfReqRes(c, codeGroup, common.CodeInvalidData, err)
	}

	// Exists Check
	codeGroupTable, err := a.Db.GetCodeGroup(codeGroup.GroupID)
	if err != nil {
		return response.ErrorfReqRes(c, codeGroup, common.CodeFailedDatabase, err)
	} else if codeGroupTable == nil {
		return response.ErrorfReqRes(c, codeGroup, common.DatabaseFalseData, err)
	}

	codeGroup.ToTable(codeGroupTable, true, "system", time.Now())

	// Start. Transaction 얻어옴
	txdb, err := a.Db.BeginTransaction()
	if err != nil {
		return response.ErrorfReqRes(c, codeGroup, common.CodeFailedDatabase, err)
	}

	// CodeGroup 갱신
	cnt, err := txdb.UpdateCodeGroup(codeGroupTable)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, codeGroup, common.CodeFailedDatabase, err)
	}
	if cnt == 0 {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, codeGroup, common.DatabaseFalseData, nil)
	}

	txErr := txdb.Commit()
	if txErr != nil {
		logger.Info("DB commit Failed.", txErr)
	}

	return response.Write(c, codeGroup, nil)
}

// DeleteCodeGroupHandler - 코드그룹 삭제
// @Tags CodeGroup
// @Summary DeleteCodeGroup
// @Description Delete code group and codes belong to
// @ID DeleteCodeGroup
// @Produce json
// @Param groupId path string true "Code Group ID"
// @Success 200 {object} response.ReturnData
// @Router /codes/groups/{groupId} [delete]
func (a *API) DeleteCodeGroupHandler(c echo.Context) error {
	// TODO: 로그인 사용자 정보 활용 방법은?
	groupId := c.Param("groupId")
	if groupId == "" {
		return response.ErrorfReqRes(c, groupId, common.CodeInvalidParm, nil)
	}

	// 존재여부 검증
	codeGroupTable, err := a.Db.GetCodeGroup(groupId)
	if err != nil {
		return response.ErrorfReqRes(c, groupId, common.CodeFailedDatabase, err)
	} else if codeGroupTable == nil {
		return response.ErrorfReqRes(c, groupId, common.DatabaseFalseData, err)
	}

	// Start. Transaction 얻어옴
	txdb, err := a.Db.BeginTransaction()
	if err != nil {
		return response.ErrorfReqRes(c, groupId, common.CodeFailedDatabase, err)
	}

	// CodeGroup 삭제
	cnt, err := txdb.DeleteCodeGroup(groupId)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, groupId, common.CodeFailedDatabase, err)
	}
	if cnt == 0 {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, groupId, common.DatabaseFalseData, nil)
	}

	// Code 삭제
	_, err = txdb.DeleteCodeByGroup(groupId)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, groupId, common.CodeFailedDatabase, err)
	}

	txErr := txdb.Commit()
	if txErr != nil {
		logger.Info("DB commit Failed.", txErr)
	}

	return response.Write(c, groupId, nil)
}

/************************
 * Code
 ************************/

// GetCodeListHandler - 전체 코드 리스트
// @Tags Code
// @Summary GetCodeList
// @Description Get all code list
// @ID GetCodeList
// @Produce json
// @Success 200 {object} response.ReturnData
// @Router /codes [get]
func (a *API) GetCodeListHandler(c echo.Context) error {
	list, err := a.Db.GetCodeList()
	if err != nil {
		return response.Errorf(c, common.CodeFailedDatabase, err)
	}
	if len(list) == 0 {
		return response.Errorf(c, common.DatabaseFalseData, err)
	}

	return response.Write(c, nil, list)
}

// GetCodeListByGroupHandler - 그룹에 속하는 코드 리스트
// @Tags Code
// @Summary GetCodeListByGroup
// @Description Get codes by group
// @ID GetCodeListByGroup
// @Produce json
// @Param groupId path string true "Code Group ID"
// @Success 200 {object} response.ReturnData
// @Router /codes/{groupId} [get]
func (a *API) GetCodeListByGroupHandler(c echo.Context) error {
	groupId := c.Param("groupId")
	if groupId == "" {
		return response.ErrorfReqRes(c, groupId, common.CodeInvalidParm, nil)
	}

	list, err := a.Db.GetCodeListByGroup(groupId)
	if err != nil {
		return response.Errorf(c, common.CodeFailedDatabase, err)
	}
	if len(list) == 0 {
		return response.Errorf(c, common.DatabaseFalseData, err)
	}

	return response.Write(c, nil, list)
}

// GetCodeHandler - 코드 상세 조회
// @Tags Code
// @Summary GetCode
// @Description Get code
// @ID GetCode
// @Produce json
// @Param groupId path string true "Code Group ID"
// @Param code path int true "Code"
// @Success 200 {object} response.ReturnData
// @Router /codes/{groupId}/{code} [get]
func (a *API) GetCodeHandler(c echo.Context) error {
	groupId := c.Param("groupId")
	if groupId == "" {
		return response.ErrorfReqRes(c, groupId, common.CodeInvalidParm, nil)
	}

	code, err := strconv.Atoi(c.Param("code"))
	if err != nil || code <= 0 {
		return response.ErrorfReqRes(c, code, common.CodeInvalidParm, nil)
	}

	data, err := a.Db.GetCode(groupId, code)
	if err != nil {
		return response.Errorf(c, common.CodeFailedDatabase, err)
	}
	if data == nil {
		return response.Errorf(c, common.DatabaseFalseData, err)
	}

	return response.Write(c, nil, data)
}

// SetCodeHandler - 코드 등록
// @Tags Code
// @Summary SetCode
// @Description Register code
// @ID SetCode
// @Produce json
// @Param code body model.Code true "Code"
// @Success 200 {object} response.ReturnData
// @Router /codes [post]
func (a *API) SetCodeHandler(c echo.Context) error {
	// TODO: 로그인 사용자 정보 활용 방법은?
	var code *model.Code

	err := getRequestData(c.Request(), &code)
	if err != nil {
		return response.ErrorfReqRes(c, code, common.CodeInvalidData, err)
	}

	// 중복 검증
	codeTable, err := a.Db.GetCode(code.GroupID, code.Code)
	if err != nil {
		return response.ErrorfReqRes(c, code, common.CodeFailedDatabase, err)
	}
	if codeTable != nil {
		return response.ErrorfReqRes(c, code, common.DatabaseExistData, err)
	}

	codeTable = &model.CodeTable{}
	code.ToTable(codeTable, false, "system", time.Now())

	// Start. Transaction 얻어옴
	txdb, err := a.Db.BeginTransaction()
	if err != nil {
		return response.ErrorfReqRes(c, code, common.CodeFailedDatabase, err)
	}

	// Code 등록
	err = txdb.InsertCode(codeTable)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, code, common.CodeFailedDatabase, err)
	}

	txErr := txdb.Commit()
	if txErr != nil {
		logger.Info("DB commit Failed.", txErr)
	}

	return response.Write(c, code, nil)
}

// UpdateCodeHandler - 코드 수정
// @Tags Code
// @Summary UpdateCode
// @Description Update code
// @ID UpdateCode
// @Produce json
// @Param code body model.Code true "Code"
// @Success 200 {object} response.ReturnData
// @Router /codes [put]
func (a *API) UpdateCodeHandler(c echo.Context) error {
	// TODO: 로그인 사용자 정보 활용 방법은?
	var code model.Code

	err := getRequestData(c.Request(), &code)
	if err != nil {
		return response.ErrorfReqRes(c, code, common.CodeInvalidData, err)
	}

	// 존재여부 검증
	codeTable, err := a.Db.GetCode(code.GroupID, code.Code)
	if err != nil {
		return response.ErrorfReqRes(c, code, common.CodeFailedDatabase, err)
	} else if codeTable == nil {
		return response.ErrorfReqRes(c, code, common.DatabaseFalseData, err)
	}

	// 변경정보 갱신
	code.ToTable(codeTable, true, "system", time.Now())

	// Start. Transaction 얻어옴
	txdb, err := a.Db.BeginTransaction()
	if err != nil {
		return response.ErrorfReqRes(c, code, common.CodeFailedDatabase, err)
	}

	// Code 갱신
	cnt, err := txdb.UpdateCode(codeTable)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, code, common.CodeFailedDatabase, err)
	}
	if cnt == 0 {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, code, common.DatabaseFalseData, nil)
	}

	txErr := txdb.Commit()
	if txErr != nil {
		logger.Info("DB commit Failed.", txErr)
	}

	return response.Write(c, code, nil)
}

// DeleteCodeHandler - 코드 삭제
// @Tags Code
// @Summary DeleteCode
// @Description Delete code
// @ID DeleteCode
// @Produce json
// @Param groupId path string true "Code Group ID"
// @Param code path int true "Code"
// @Success 200 {object} response.ReturnData
// @Router /codes/{groupId}/{code} [delete]
func (a *API) DeleteCodeHandler(c echo.Context) error {
	// TODO: 로그인 사용자 정보 활용 방법은?
	groupId := c.Param("groupId")
	if groupId == "" {
		return response.ErrorfReqRes(c, groupId, common.CodeInvalidParm, nil)
	}

	code, err := strconv.Atoi(c.Param("code"))
	if err != nil || code <= 0 {
		return response.ErrorfReqRes(c, code, common.CodeInvalidParm, nil)
	}

	// 존재여부 검증
	codeTable, err := a.Db.GetCode(groupId, code)
	if err != nil {
		return response.ErrorfReqRes(c, code, common.CodeFailedDatabase, err)
	} else if codeTable == nil {
		return response.ErrorfReqRes(c, code, common.DatabaseFalseData, err)
	}

	// Start. Transaction 얻어옴
	txdb, err := a.Db.BeginTransaction()
	if err != nil {
		return response.ErrorfReqRes(c, groupId, common.CodeFailedDatabase, err)
	}

	// Code 삭제
	cnt, err := txdb.DeleteCode(groupId, code)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, groupId, common.CodeFailedDatabase, err)
	}
	if cnt == 0 {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, groupId, common.DatabaseFalseData, nil)
	}

	txErr := txdb.Commit()
	if txErr != nil {
		logger.Info("DB commit Failed.", txErr)
	}

	return response.Write(c, groupId, nil)
}
