package api

import (
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/response"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	mr "github.com/acornsoft-edgecraft/edgecraft-api/pkg/model/response"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

func (a *API) AllCloudNodeListHandler(c echo.Context) error {
	res, err := a.Db.GetAllCloudNode()
	if err != nil {
		return response.Errorf(c, common.CodeFailedDatabase, err)
	}
	return response.Write(c, nil, res)
}

func (a *API) GetCloudNodeHandler(c echo.Context) error {
	// check param UID
	cloudUid, err := uuid.FromString(c.Param("cloudUid"))
	if err != nil {
		return response.ErrorfReqRes(c, cloudUid, common.CodeInvalidParm, err)
	}

	resCloud, err := a.Db.GetCloud(cloudUid)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if resCloud == nil {
		return response.ErrorfReqRes(c, resCloud, common.DatabaseFalseData, err)
	}

	var res model.Cloud

	return response.Write(c, nil, &res)
}

func (a *API) SelectCloudNodeHandler(c echo.Context) error {
	// check param UID
	cloudUid, err := uuid.FromString(c.Param("cloudUid"))
	if err != nil {
		return response.ErrorfReqRes(c, cloudUid, common.CodeInvalidParm, err)
	}

	resCloud, err := a.Db.GetCloud(cloudUid)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if resCloud == nil {
		return response.ErrorfReqRes(c, resCloud, common.DatabaseFalseData, err)
	}

	var res mr.RegisterCloud

	return response.Write(c, nil, &res)
}

func (a *API) UpdateCloudNodeHandler(c echo.Context) error {
	// check param UID
	cloudUid, err := uuid.FromString(c.Param("cloudUid"))
	if err != nil {
		return response.ErrorfReqRes(c, cloudUid, common.CodeInvalidParm, err)
	}

	now := time.Now()

	var req model.Cloud
	err = getRequestData(c.Request(), &req)
	if err != nil {
		return response.ErrorfReqRes(c, req, common.CodeInvalidData, err)
	}
	req.CloudUID = &cloudUid
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
		return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
	}

	// End. Transaction Commit
	txErr := txdb.Commit()
	if txErr != nil {
		logger.Info("DB commit Failed.", txErr)
	}

	if count == 0 {
		return response.ErrorfReqRes(c, cloudUid, common.DatabaseEmptyData, nil)
	}

	return response.Write(c, nil, count)
}

func (a *API) DeleteCloudNodeHandler(c echo.Context) error {
	// check param UID
	cloudUid, err := uuid.FromString(c.Param("cloudUid"))
	if err != nil {
		return response.ErrorfReqRes(c, cloudUid, common.CodeInvalidParm, err)
	}

	// -- Service Logic
	// Start. Transaction 얻어옴
	txdb, err := a.Db.BeginTransaction()
	if err != nil {
		return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
	}

	// 1. cloud 삭제
	count, err := txdb.DeleteCloud(cloudUid)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB Rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
	}

	// End. Transaction Commit
	txErr := txdb.Commit()
	if txErr != nil {
		logger.Info("DB commit Failed.", txErr)
	}

	if count == 0 {
		return response.ErrorfReqRes(c, cloudUid, common.DatabaseEmptyData, nil)
	}

	return response.Write(c, nil, count)
}
