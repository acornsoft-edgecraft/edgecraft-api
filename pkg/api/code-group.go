package api

import (
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/response"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	mr "github.com/acornsoft-edgecraft/edgecraft-api/pkg/model/response"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

func (a *API) AllCodeGroupListHandler(c echo.Context) error {
	res, err := a.Db.GetAllCodeGroup()
	if err != nil {
		return response.Errorf(c, common.CodeFailedDatabase, err)
	}
	return response.Write(c, nil, res)
}

func (a *API) RegisterCodeGroupHandler(c echo.Context) error {
	now := time.Now()

	var req []model.CodeGroup
	err := getRequestData(c.Request(), &req)
	if err != nil {
		return response.ErrorfReqRes(c, req, common.CodeInvalidData, err)
	}

	// -- Service Logic
	// Start. Transaction 얻어옴
	txdb, err := a.Db.BeginTransaction()
	if err != nil {
		return response.ErrorfReqRes(c, req, common.CodeFailedDatabase, err)
	}

	// CodeGroup 등록
	var codeGroup []model.CodeGroup = req
	for _, data := range codeGroup {
		data.CreatedAt = &now
		err = txdb.CreateCodeGroup(&data)
		if err != nil {
			txErr := txdb.Rollback()
			if txErr != nil {
				logger.Info("DB Rollback Failed.", txErr)
			}
			return response.ErrorfReqRes(c, req, common.CodeFailedDatabase, err)
		}
	}

	// End. Transaction Commit
	txErr := txdb.Commit()
	if txErr != nil {
		logger.Info("DB commit Failed.", txErr)
	}

	return response.Write(c, req, codeGroup)
}

func (a *API) GetCodeGroupHandler(c echo.Context) error {
	// check param UID
	codeGroupUid, err := uuid.FromString(c.Param("group_code_uid"))
	if err != nil {
		return response.ErrorfReqRes(c, codeGroupUid, common.CodeInvalidParm, err)
	}

	resCloud, err := a.Db.GetCloud(codeGroupUid)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if resCloud == nil {
		return response.ErrorfReqRes(c, resCloud, common.DatabaseFalseData, err)
	}

	var res model.Cloud

	return response.Write(c, nil, &res)
}

func (a *API) SelectCodeGroupHandler(c echo.Context) error {
	// check param UID
	codeGroupUid, err := uuid.FromString(c.Param("group_code_uid"))
	if err != nil {
		return response.ErrorfReqRes(c, codeGroupUid, common.CodeInvalidParm, err)
	}

	resCloud, err := a.Db.GetCloud(codeGroupUid)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if resCloud == nil {
		return response.ErrorfReqRes(c, resCloud, common.DatabaseFalseData, err)
	}

	var res mr.RegisterCloud

	return response.Write(c, nil, &res)
}

func (a *API) CodeGroupSearchHandler(c echo.Context) error {
	var req model.CodeGroup
	err := getRequestData(c.Request(), &req)
	if err != nil {
		return response.ErrorfReqRes(c, req, common.CodeInvalidData, err)
	}

	utils.Print(&req)

	res, err := a.Db.SearchCodeGroup(req)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if res == nil {
		return response.ErrorfReqRes(c, res, common.DatabaseFalseData, err)
	}
	return response.Write(c, nil, res)
}

func (a *API) UpdateCodeGroupHandler(c echo.Context) error {
	// check param UID
	codeGroupUid, err := uuid.FromString(c.Param("group_code_uid"))
	if err != nil {
		return response.ErrorfReqRes(c, codeGroupUid, common.CodeInvalidParm, err)
	}

	now := time.Now()

	var req model.Cloud
	err = getRequestData(c.Request(), &req)
	if err != nil {
		return response.ErrorfReqRes(c, req, common.CodeInvalidData, err)
	}
	req.CloudUID = &codeGroupUid
	req.UpdatedAt = &now

	// -- Service Logic
	// Start. Transaction 얻어옴
	txdb, err := a.Db.BeginTransaction()
	if err != nil {
		return response.ErrorfReqRes(c, req, common.CodeFailedDatabase, err)
	}

	// Clud 등록 업데이트
	count, err := txdb.UpdateCloud(&req)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB Rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, codeGroupUid, common.CodeFailedDatabase, err)
	}

	// End. Transaction Commit
	txErr := txdb.Commit()
	if txErr != nil {
		logger.Info("DB commit Failed.", txErr)
	}

	if count == 0 {
		return response.ErrorfReqRes(c, codeGroupUid, common.DatabaseEmptyData, nil)
	}

	return response.Write(c, nil, count)
}

func (a *API) DeleteCodeGroupHandler(c echo.Context) error {
	// check param UID
	codeGroupUid, err := uuid.FromString(c.Param("group_code_uid"))
	if err != nil {
		return response.ErrorfReqRes(c, codeGroupUid, common.CodeInvalidParm, err)
	}

	// -- Service Logic
	// Start. Transaction 얻어옴
	txdb, err := a.Db.BeginTransaction()
	if err != nil {
		return response.ErrorfReqRes(c, codeGroupUid, common.CodeFailedDatabase, err)
	}

	// 1. cloud 삭제
	count, err := txdb.DeleteCloud(codeGroupUid)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB Rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, codeGroupUid, common.CodeFailedDatabase, err)
	}

	// End. Transaction Commit
	txErr := txdb.Commit()
	if txErr != nil {
		logger.Info("DB commit Failed.", txErr)
	}

	if count == 0 {
		return response.ErrorfReqRes(c, codeGroupUid, common.DatabaseEmptyData, nil)
	}

	return response.Write(c, nil, count)
}