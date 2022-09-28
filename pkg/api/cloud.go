package api

import (
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/response"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"

	//mr "github.com/acornsoft-edgecraft/edgecraft-api/pkg/model/response"

	"github.com/labstack/echo/v4"
)

// GetCloudListHandler - 전체 클라우드 리스트
// @Tags Cloud
// @Summary GetCloudList
// @Description Get all cloud list
// @ID GetCloudList
// @Produce json
// @Success 200 {object} response.ReturnData
// @Router /clouds [get]
func (a *API) GetCloudListHandler(c echo.Context) error {
	res, err := a.Db.GetCloudList()
	if err != nil {
		return response.Errorf(c, common.CodeFailedDatabase, err)
	}
	return response.Write(c, nil, res)
}

// GetCloudHandler - 클라우드 상세 정보
// @Tags Cloud
// @Summary GetCloud
// @Description Get specific cloud
// @ID GetCloud
// @Produce json
// @Param cloudUid path string true "cloudUid"
// @Success 200 {object} response.ReturnData
// @Router /clouds/{cloudUid} [get]
func (a *API) GetCloudHandler(c echo.Context) error {
	cloudUid := c.Param("cloudUid")
	if cloudUid == "" {
		return response.ErrorfReqRes(c, cloudUid, common.CodeInvalidParm, nil)
	}

	cloudSet := &model.CloudSet{}

	// Cloud 조회
	cloudTable, err := a.Db.GetCloud(cloudUid)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if cloudTable == nil {
		return response.ErrorfReqRes(c, cloudTable, common.DatabaseFalseData, err)
	}
	cloudTable.ToSet(cloudSet)

	// Cluster 조회
	clusters, err := a.Db.SelectClusters(cloudUid)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if len(clusters) == 0 {
		return response.ErrorfReqRes(c, clusters, common.DatabaseFalseData, err)
	}

	clusters[0].ToSet(cloudSet)

	// Node 조회
	nodes, err := a.Db.SelectNodes(cloudUid, *clusters[0].ClusterUid)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if nodes == nil {
		return response.ErrorfReqRes(c, nodes, common.DatabaseFalseData, err)
	}

	cloudSet.Nodes = &model.NodesInfo{}
	cloudSet.Nodes.FromTable(clusters[0], nodes)

	return response.Write(c, nil, cloudSet)
}

// SetCloudHandler - 클라우드 등록
// @Tags Cloud
// @Summary SetCloud
// @Description Register cloud
// @ID SetCloud
// @Produce json
// @Param cloudSet body model.CloudSet true "Cloud Set"
// @Success 200 {object} response.ReturnData
// @Router /clouds [post]
func (a *API) SetCloudHandler(c echo.Context) error {
	// TODO: 로그인 사용자 정보 활용 방법은?
	var cloudSet model.CloudSet

	err := getRequestData(c.Request(), &cloudSet)
	if err != nil {
		return response.ErrorfReqRes(c, cloudSet, common.CodeInvalidData, err)
	}

	cloudTable, clusterTable, nodeTables := cloudSet.ToTable(false, "system", time.Now())

	// Start. Transaction 얻어옴
	txdb, err := a.Db.BeginTransaction()
	if err != nil {
		return response.ErrorfReqRes(c, cloudSet, common.CodeFailedDatabase, err)
	}

	// Cloud 등록
	err = txdb.InsertCloud(cloudTable)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudSet, common.CodeFailedDatabase, err)
	}

	// Cluster 등록
	err = txdb.InsertCluster(clusterTable)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudSet, common.CodeFailedDatabase, err)
	}

	// Node 등록
	for _, nodeTable := range nodeTables {
		err = txdb.InsertNode(nodeTable)
		if err != nil {
			txErr := txdb.Rollback()
			if txErr != nil {
				logger.Info("DB rollback Failed.", txErr)
			}
			return response.ErrorfReqRes(c, cloudSet, common.CodeFailedDatabase, err)
		}
	}

	txErr := txdb.Commit()
	if txErr != nil {
		logger.Info("DB commit Failed.", txErr)
	}

	return response.Write(c, cloudSet, nil)
}

// UpdateCloudHandler - 클라우드 갱신
// @Tags Cloud
// @Summary UpdateCloud
// @Description Update cloud
// @ID UpdateCloud
// @Produce json
// @Param cloudUid path string true "CloudUid"
// @Param cloudSet body model.CloudSet true "Cloud Set"
// @Success 200 {object} response.ReturnData
// @Router /clouds/{cloudUid} [put]
func (a *API) UpdateCloudHandler(c echo.Context) error {
	// TODO: 로그인 사용자 정보 활용 방법은?
	cloudUid := c.Param("cloudUid")
	if cloudUid == "" {
		return response.ErrorfReqRes(c, cloudUid, common.CodeInvalidParm, nil)
	}

	var cloudSet model.CloudSet
	var at time.Time = time.Now()

	err := getRequestData(c.Request(), &cloudSet)
	if err != nil {
		return response.ErrorfReqRes(c, cloudSet, common.CodeInvalidData, err)
	}

	cloudTable, clusterTable, nodeTables := cloudSet.ToTable(true, "system", at)

	// Start. Transaction 얻어옴
	txdb, err := a.Db.BeginTransaction()
	if err != nil {
		return response.ErrorfReqRes(c, cloudSet, common.CodeFailedDatabase, err)
	}

	// Cloud 갱신
	cnt, err := txdb.UpdateCloud(cloudTable)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudSet, common.CodeFailedDatabase, err)
	}
	if cnt == 0 {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudSet, common.DatabaseFalseData, err)
	}

	// Cluster 갱신
	cnt, err = txdb.UpdateCluster(clusterTable)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudSet, common.CodeFailedDatabase, err)
	}
	if cnt == 0 {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudSet, common.DatabaseFalseData, err)
	}

	// 기존 Nodes 삭제
	cnt, err = txdb.DeleteNodes(*clusterTable.CloudUid, *clusterTable.ClusterUid)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudSet, common.CodeFailedDatabase, err)
	}
	if cnt == 0 {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudSet, common.DatabaseFalseData, err)
	}

	// Nodes 추가
	for _, nodeTable := range nodeTables {
		err = txdb.InsertNode(nodeTable)
		if err != nil {
			txErr := txdb.Rollback()
			if txErr != nil {
				logger.Info("DB rollback Failed.", txErr)
			}
			return response.ErrorfReqRes(c, cloudSet, common.CodeFailedDatabase, err)
		}
	}

	txErr := txdb.Commit()
	if txErr != nil {
		logger.Info("DB commit Failed.", txErr)
	}

	return response.Write(c, cloudSet, nil)
}

// DeleteCloudHandler - 클라우드 삭제
// @Tags Cloud
// @Summary DeleteCloud
// @Description Delete cloud
// @ID DeleteCloud
// @Produce json
// @Param cloudUid path string true "CloudUid"
// @Success 200 {object} response.ReturnData
// @Router /clouds/{cloudUid} [delete]
func (a *API) DeleteCloudHandler(c echo.Context) error {
	// TODO: 로그인 사용자 정보 활용 방법은?
	cloudUid := c.Param("cloudUid")
	if cloudUid == "" {
		return response.ErrorfReqRes(c, cloudUid, common.CodeInvalidParm, nil)
	}

	// Start. Transaction 얻어옴
	txdb, err := a.Db.BeginTransaction()
	if err != nil {
		return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
	}

	// Cloud 삭제
	cnt, err := txdb.DeleteCloud(cloudUid)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
	}
	if cnt == 0 {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudUid, common.DatabaseFalseData, nil)
	}

	// Cluster 삭제
	cnt, err = txdb.DeleteCloudClusters(cloudUid)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
	}
	if cnt == 0 {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudUid, common.DatabaseFalseData, nil)
	}

	// Cloud Nodes 삭제
	cnt, err = txdb.DeleteCloudNodes(cloudUid)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
	}
	if cnt == 0 {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudUid, common.DatabaseFalseData, nil)
	}

	txErr := txdb.Commit()
	if txErr != nil {
		logger.Info("DB commit Failed.", txErr)
	}

	return response.Write(c, cloudUid, nil)
}
