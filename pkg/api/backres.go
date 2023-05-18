/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package api

import (
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/response"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/db"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	"github.com/labstack/echo/v4"
)

/*******************************
 ** Backup for Cluster
 *******************************/

// SetBackupHandler - 클러스터의 백업 실행 (Velero)
// @Tags        Openstack-Cluster-Backup
// @Summary     SetBackup
// @Description 클러스터의 백업 실행 (Velero)
// @ID          SetBackup
// @Produce     json
// @Param       cloudId 	path     string true "Cloud ID"
// @Param       clusterId 	path     string true "Cluster ID"
// @Param       name	 	path     string true "Name"
// @Success     200     {object} response.ReturnData
// @Router      /clouds/{cloudId}/clusters/{clusterId}/backup [post]
func (a *API) SetBackupHandler(c echo.Context) error {
	// 파라미터 검사
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}
	clusterId := c.Param("clusterId")
	if clusterId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}
	name := c.Param("name")
	if clusterId == "" {
		return response.ErrorfReqRes(c, name, common.CodeInvalidParm, nil)
	}

	// 클러스터 정보 조회
	clusterTable, err := a.Db.GetOpenstackCluster(cloudId, clusterId)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	}
	if clusterTable == nil {
		return response.ErrorfReqRes(c, clusterTable, common.ClusterNotFound, nil)
	}

	// 클러스터 상태 조회
	if *clusterTable.Status != common.StatusProvisioned {
		return response.ErrorfReqRes(c, nil, common.BackResOnlyProvisioned, err)
	}

	// 백업 정보 생성
	backresInfo := model.NewBackResInfo(cloudId, clusterId, name, true)

	// 백업 실행
	// err = kubemethod.ApplyBackup(backresInfo)
	// if err != nil {
	// 	return response.ErrorfReqRes(c, clusterTable, common.BackResFailed, err)
	// }

	// 데이터베이스 저장
	backresTable := backresInfo.ToTable("system", time.Now())
	err = a.Db.TransactionScope(func(txDB db.DB) error {
		err = txDB.InsertBackRes(backresTable)
		if err != nil {
			return err
		}

		return nil
	})

	return response.WriteWithCode(c, nil, common.BackResExecuing, nil)
}

// GetBackupListHandler - 클러스터의 백업 리스트
// @Tags        Openstack-Cluster-Backup
// @Summary     GetBackupList
// @Description 클러스터의 백업 리스트
// @ID          GetBackupList
// @Produce     json
// @Param       cloudId 	path     string true "Cloud ID"
// @Success     200     {object} response.ReturnData
// @Router      /clouds/{cloudId}/backup [get]
func (a *API) GetBackupListHandler(c echo.Context) error {
	// 파라미터 검사
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}

	result, err := a.Db.GetBackupList(cloudId)
	if err != nil {
		return response.Errorf(c, common.CodeFailedDatabase, err)
	}

	return response.Write(c, nil, result)
}

// DeleteBackupHandler - 클러스터의 백업 삭제
// @Tags        Openstack-Cluster-Backup
// @Summary     DeleteBackup
// @Description 클러스터의 백업 삭제
// @ID          DeleteBackup
// @Produce     json
// @Param       cloudId 	path     string true "Cloud ID"
// @Param       clusterId 	path     string true "Cluster ID"
// @Param       backresId 	path     string true "BackRes ID"
// @Success     200     {object} response.ReturnData
// @Router      /clouds/{cloudId}/clusters/{clusterId}/backup [delete]
func (a *API) DeleteBackupHandler(c echo.Context) error {
	// TODO: 로그인 사용자 정보 활용 방법은?
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}

	// Start. Transaction 얻어옴
	txdb, err := a.Db.BeginTransaction()
	if err != nil {
		return response.ErrorfReqRes(c, cloudId, common.CodeFailedDatabase, err)
	}

	// Cloud 삭제
	cnt, err := txdb.DeleteCloud(cloudId)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudId, common.CodeFailedDatabase, err)
	}
	if cnt == 0 {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudId, common.DatabaseFalseData, nil)
	}

	// Cluster 삭제
	cnt, err = txdb.DeleteClusters(cloudId)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudId, common.CodeFailedDatabase, err)
	}
	if cnt == 0 {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudId, common.DatabaseFalseData, nil)
	}

	// Cloud Nodes 삭제
	cnt, err = txdb.DeleteCloudNodes(cloudId)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudId, common.CodeFailedDatabase, err)
	}
	if cnt == 0 {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudId, common.DatabaseFalseData, nil)
	}

	txErr := txdb.Commit()
	if txErr != nil {
		logger.Info("DB commit Failed.", txErr)
	}

	return response.Write(c, cloudId, nil)
}

/*******************************
 ** Restore for Cluster
 *******************************/

// SetRestoreHandler - 클러스터의 복원 실행 (Velero)
// @Tags        Openstack-Cluster-Restore
// @Summary     SetRestore
// @Description 클러스터의 복원 실행 (Velero)
// @ID          SetRestore
// @Produce     json
// @Param       cloudId 	path     string true "Cloud ID"
// @Param       clusterId 	path     string true "Cluster ID"
// @Param       backresId 	path     string true "BackRes ID"
// @Param       name	 	path     string true "Name"
// @Success     200     {object} response.ReturnData
// @Router      /clouds/{cloudId}/clusters/{clusterId}/restore [post]
func (a *API) SetRestoreHandler(c echo.Context) error {
	// 파라미터 검사
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}
	clusterId := c.Param("clusterId")
	if clusterId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}
	backresId := c.Param("backresId")
	if backresId == "" {
		return response.ErrorfReqRes(c, backresId, common.CodeInvalidParm, nil)
	}
	name := c.Param("name")
	if clusterId == "" {
		return response.ErrorfReqRes(c, name, common.CodeInvalidParm, nil)
	}

	// 클러스터 정보 조회
	clusterTable, err := a.Db.GetOpenstackCluster(cloudId, clusterId)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	}
	if clusterTable == nil {
		return response.ErrorfReqRes(c, clusterTable, common.ClusterNotFound, nil)
	}

	// 클러스터 상태 조회
	if *clusterTable.Status != common.StatusProvisioned {
		return response.ErrorfReqRes(c, nil, common.BackResOnlyProvisioned, err)
	}

	// 기존 백업 정보 조회
	// backup, err := a.Db.GetBackup(cloudId, clusterId, backresId)
	// if err != nil {
	// 	return response.ErrorfReqRes(c, backresId, common.BackupNotAvailable, err)
	// }

	// 복원 정보 생성
	backresInfo := model.NewBackResInfo(cloudId, clusterId, name, true)

	// 복원 실행
	// err = kubemethod.ApplyRestore(backresInfo, backup.Name)
	// if err != nil {
	// 	return response.ErrorfReqRes(c, clusterTable, common.BackResFailed, err)
	// }

	// 데이터베이스 저장
	backresTable := backresInfo.ToTable("system", time.Now())
	err = a.Db.TransactionScope(func(txDB db.DB) error {
		err = txDB.InsertBackRes(backresTable)
		if err != nil {
			return err
		}

		return nil
	})

	return response.WriteWithCode(c, nil, common.BackResExecuing, nil)
}

// GetRestoreListHandler - 클러스터의 복원 리스트
// @Tags        Openstack-Cluster-Restore
// @Summary     GetRestoreList
// @Description 클러스터의 복원 리스트
// @ID          GetRestoreList
// @Produce     json
// @Param       cloudId 	path     string true "Cloud ID"
// @Success     200     {object} response.ReturnData
// @Router      /clouds/{cloudId}/restore [get]
func (a *API) GetRestoreListHandler(c echo.Context) error {
	// 파라미터 검사
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}

	result, err := a.Db.GetRestoreList(cloudId)
	if err != nil {
		return response.Errorf(c, common.CodeFailedDatabase, err)
	}

	return response.Write(c, nil, result)
}

// DeleteRestoreHandler - 클러스터의 복원 삭제
// @Tags        Openstack-Cluster-Restore
// @Summary     DeleteRestore
// @Description 클러스터의 복원 삭제
// @ID          DeleteRestore
// @Produce     json
// @Param       cloudId 	path     string true "Cloud ID"
// @Param       clusterId 	path     string true "Cluster ID"
// @Param       backresId 	path     string true "BackRes ID"
// @Success     200     {object} response.ReturnData
// @Router      /clouds/{cloudId}/clusters/{clusterId}/restore [delete]
func (a *API) DeleteRestoreHandler(c echo.Context) error {
	// TODO: 로그인 사용자 정보 활용 방법은?
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}

	// Start. Transaction 얻어옴
	txdb, err := a.Db.BeginTransaction()
	if err != nil {
		return response.ErrorfReqRes(c, cloudId, common.CodeFailedDatabase, err)
	}

	// Cloud 삭제
	cnt, err := txdb.DeleteCloud(cloudId)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudId, common.CodeFailedDatabase, err)
	}
	if cnt == 0 {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudId, common.DatabaseFalseData, nil)
	}

	// Cluster 삭제
	cnt, err = txdb.DeleteClusters(cloudId)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudId, common.CodeFailedDatabase, err)
	}
	if cnt == 0 {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudId, common.DatabaseFalseData, nil)
	}

	// Cloud Nodes 삭제
	cnt, err = txdb.DeleteCloudNodes(cloudId)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudId, common.CodeFailedDatabase, err)
	}
	if cnt == 0 {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudId, common.DatabaseFalseData, nil)
	}

	txErr := txdb.Commit()
	if txErr != nil {
		logger.Info("DB commit Failed.", txErr)
	}

	return response.Write(c, cloudId, nil)
}
