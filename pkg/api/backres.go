/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package api

import (
	"errors"
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/response"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/db"
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
// @Param       cloudId 		path     string true "Cloud ID"
// @Param       clusterId 		path     string true "Cluster ID"
// @Param       backresParam	body	 model.BackResParam	true "BackResParam"
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
	var backresParam model.BackResParam
	err := getRequestData(c.Request(), &backresParam)
	if err != nil {
		return response.ErrorfReqRes(c, backresParam, common.CodeInvalidData, err)
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

	// 동일한 백업이 존재하는지 검증
	exists, err := a.Db.CheckBackResDuplicate(backresParam.Name)
	if err != nil {
		return response.ErrorfReqRes(c, backresParam.Name, common.CodeFailedDatabase, nil)

	} else if exists {
		return response.ErrorfReqRes(c, backresParam.Name, common.BackResDuplicated, nil)
	}

	// 백업 정보 생성
	backresInfo := model.NewBackResInfo(cloudId, clusterId, backresParam.Name, "", true)

	// 백업 실행
	err = ProvisioningBackRes(a.Worker, a.Db, *clusterTable.Name, "velero", backresInfo)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.BackResJobFailed, err)
	}

	// 데이터베이스 저장
	backresTable := backresInfo.ToTable("system", time.Now())
	err = a.Db.TransactionScope(func(txDB db.DB) error {
		err = txDB.InsertBackRes(backresTable)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return response.ErrorfReqRes(c, backresTable, common.CodeFailedDatabase, err)
	}

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
// @Router      /clouds/{cloudId}/clusters/{clusterId}/backup/{backresId} [delete]
func (a *API) DeleteBackupHandler(c echo.Context) error {
	// TODO: 로그인 사용자 정보 활용 방법은?
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}
	clusterId := c.Param("clusterId")
	if clusterId == "" {
		return response.ErrorfReqRes(c, clusterId, common.CodeInvalidParm, nil)
	}
	backresId := c.Param("backresId")
	if backresId == "" {
		return response.ErrorfReqRes(c, backresId, common.CodeInvalidParm, nil)
	}

	// 백업 CR 삭제
	// err = kubemethod.ApplyBackup(backresInfo)
	// if err != nil {
	// 	return response.ErrorfReqRes(c, clusterTable, common.BackResFailed, err)
	// }

	// 데이터베이스 삭제
	err := a.Db.TransactionScope(func(txDB db.DB) error {
		cnt, err := txDB.DeleteBackRes(cloudId, clusterId, backresId)
		if err != nil {
			return err
		}
		if cnt == 0 {
			return errors.New("cannot find backup for deleting")
		}

		return nil
	})

	if err != nil {
		return response.ErrorfReqRes(c, backresId, common.DatabaseFalseData, err)
	}

	return response.Write(c, backresId, nil)
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
// @Param       cloudId 		path     string true "Cloud ID"
// @Param       clusterId 		path     string true "Cluster ID"
// @Param       backresParam	body	 model.BackResParam	true "BackResParam"
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
	var backresParam model.BackResParam
	err := getRequestData(c.Request(), &backresParam)
	if err != nil {
		return response.ErrorfReqRes(c, backresParam, common.CodeInvalidData, err)
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
	backupTable, err := a.Db.GetBackup(backresParam.BackResId)
	if err != nil {
		return response.ErrorfReqRes(c, backresParam.BackResId, common.BackupNotAvailable, err)
	}
	if backupTable == nil {
		return response.ErrorfReqRes(c, clusterTable, common.BackupNotFound, nil)
	}

	// 복원 정보 생성
	backresInfo := model.NewBackResInfo(cloudId, clusterId, backresParam.Name, *backupTable.Name, false)

	// 복원 실행
	err = ProvisioningBackRes(a.Worker, a.Db, *clusterTable.Name, "velero", backresInfo)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.BackResJobFailed, err)
	}

	// 데이터베이스 저장
	backresTable := backresInfo.ToTable("system", time.Now())
	err = a.Db.TransactionScope(func(txDB db.DB) error {
		err = txDB.InsertBackRes(backresTable)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return response.ErrorfReqRes(c, backresTable, common.CodeFailedDatabase, err)
	}

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
// @Router      /clouds/{cloudId}/clusters/{clusterId}/restore/{backresId} [delete]
func (a *API) DeleteRestoreHandler(c echo.Context) error {
	// TODO: 로그인 사용자 정보 활용 방법은?
	// TODO: 로그인 사용자 정보 활용 방법은?
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}
	clusterId := c.Param("clusterId")
	if clusterId == "" {
		return response.ErrorfReqRes(c, clusterId, common.CodeInvalidParm, nil)
	}
	backresId := c.Param("backresId")
	if backresId == "" {
		return response.ErrorfReqRes(c, backresId, common.CodeInvalidParm, nil)
	}

	// 복원 CR 삭제
	// err = kubemethod.ApplyBackup(backresInfo)
	// if err != nil {
	// 	return response.ErrorfReqRes(c, clusterTable, common.BackResFailed, err)
	// }

	// 데이터베이스 삭제
	err := a.Db.TransactionScope(func(txDB db.DB) error {
		cnt, err := txDB.DeleteBackRes(cloudId, clusterId, backresId)
		if err != nil {
			return err
		}
		if cnt == 0 {
			return errors.New("cannot find restore for deleting")
		}

		return nil
	})

	if err != nil {
		return response.ErrorfReqRes(c, backresId, common.DatabaseFalseData, err)
	}

	return response.Write(c, backresId, nil)
}
