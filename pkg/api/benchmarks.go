/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package api

import (
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/kubemethod"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/response"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/db"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	"github.com/labstack/echo/v4"
)

/*******************************
 ** Benchmarks for Openstack
 *******************************/

// SetBenchmarksHandler - 클러스터의 CIS Benchmarks 실행 (Openstack)
// @Tags        Openstack-Cluster-Benchmarks
// @Summary     SetBenchmarks
// @Description 클러스터에 CIS Benchmarks 실행 (Openstack)
// @ID          SetBenchmarks
// @Produce     json
// @Param       cloudId 	path     string true "Cloud ID"
// @Param       clusterId 	path     string true "Cluster ID"
// @Success     200     {object} response.ReturnData
// @Router      /clouds/{cloudId}/clusters/{clusterId}/benchmarks [post]
func (a *API) SetBenchmarksHandler(c echo.Context) error {
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}
	clusterId := c.Param("clusterId")
	if clusterId == "" {
		return response.ErrorfReqRes(c, clusterId, common.CodeInvalidParm, nil)
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
		return response.ErrorfReqRes(c, nil, common.BenchmarksOnlyProvisioned, err)
	}

	benchmarksSet := &model.OpenstackBenchmarksSet{}
	benchmarksSet.NewKey()

	// benchmarks 실행
	err = kubemethod.SetEdgeBenchmarks(*clusterTable, benchmarksSet.BenchmarksUid, a.Config.Benchmarks)
	if err != nil {
		return response.ErrorfReqRes(c, clusterTable, common.BenchmarksSetFailed, err)
	}

	benchmarksTable := benchmarksSet.ToTable(clusterTable.CloudUid, clusterTable.ClusterUid, "system", time.Now())

	// 트랜잭션 구간 처리
	err = a.Db.TransactionScope(func(txDB db.DB) error {
		// DB 등록
		err = txDB.InsertOpenstackBenchmarks(benchmarksTable)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return response.ErrorfReqRes(c, benchmarksTable, common.CodeFailedDatabase, err)
	}

	return response.WriteWithCode(c, nil, common.BenchmarksExecuing, nil)
}

// GetBenchmarksListHandler - 클러스터 Benchmarks 결과 리스트 (Openstack)
// @Tags        Openstack-Cluster-Benchmarks
// @Summary     GetBenchmarksList
// @Description 클러스터 Benchmarks 결과 리스트 (Openstack)
// @ID          GetBenchmarksList
// @Produce     json
// @Param       cloudId 	path     string true "Cloud ID"
// @Param       clusterId 	path     string true "Cluster ID"
// @Success     200     {object} response.ReturnData
// @Router      /clouds/{cloudId}/clusters/{clusterId}/benchmarks [get]
func (a *API) GetBenchmarksListHandler(c echo.Context) error {
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}
	clusterId := c.Param("clusterId")
	if clusterId == "" {
		return response.ErrorfReqRes(c, clusterId, common.CodeInvalidParm, nil)
	}

	result, err := a.Db.GetOpenstackBenchmarksList(cloudId, clusterId)
	if err != nil {
		return response.Errorf(c, common.CodeFailedDatabase, err)
	}

	return response.Write(c, nil, result)
}

// GetBenchmarksHandler - 클러스터 Benchmarks 결과 상세 조회 (Openstack)
// @Tags        Openstack-Cluster-Benchmarks
// @Summary     GetBenchmarks
// @Description 클러스터 Benchmarks 결과 상세 조회 (Openstack)
// @ID          GetBenchmarks
// @Produce     json
// @Param       cloudId   		path     string true "Cloud ID"
// @Param       clusterId 		path     string true "Cluster ID"
// @Param       benchmarksId	path     string true "Benchmarks ID"
// @Success     200       {object} response.ReturnData
// @Router      /clouds/{cloudId}/clusters/{clusterId}/benchmarks/{benchmarksId} [get]
func (a *API) GetBenchmarksHandler(c echo.Context) error {
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}

	clusterId := c.Param("clusterId")
	if clusterId == "" {
		return response.ErrorfReqRes(c, clusterId, common.CodeInvalidParm, nil)
	}

	benchmarksId := c.Param("benchmarksId")
	if benchmarksId == "" {
		return response.ErrorfReqRes(c, clusterId, common.CodeInvalidParm, nil)
	}

	// Cluster Benchmarks 정보 조회
	benchmarksTable, err := a.Db.GetOpenstackBenchmarks(cloudId, clusterId, benchmarksId)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	}
	if benchmarksTable == nil {
		return response.ErrorfReqRes(c, benchmarksTable, common.ClusterNotFound, err)
	}

	return response.Write(c, nil, benchmarksTable)
}
