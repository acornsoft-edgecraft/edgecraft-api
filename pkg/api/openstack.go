/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package api

import (
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/response"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	"github.com/labstack/echo/v4"
)

/*******************************
 ** Cluster for Openstack
 *******************************/

// GetClusterListHandler - 전체 클러스터 리스트 (Openstack)
// @Tags        Openstack-Cluster
// @Summary     GetClusterList
// @Description 전체 클러스터 리스트 (Openstack)
// @ID          GetClusterList
// @Produce     json
// @Param       cloudId path     string true "Cloud ID"
// @Success     200     {object} response.ReturnData
// @Router      /clouds/{cloudId}/clusters [get]
func (a *API) GetClusterListHandler(c echo.Context) error {
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}

	result, err := a.Db.GetOpenstackClusters(cloudId)
	if err != nil {
		return response.Errorf(c, common.CodeFailedDatabase, err)
	}

	return response.Write(c, nil, result)
}

// GetClusterHandler - 클러스터 상세 조회 (Openstack)
// @Tags        Openstack-Cluster
// @Summary     GetCluster
// @Description 클러스터 상세 조회 (Openstack)
// @ID          GetCluster
// @Produce     json
// @Param       cloudId   path     string true "Cloud ID"
// @Param       clusterId path     string true "Cluster ID"
// @Success     200       {object} response.ReturnData
// @Router      /clouds/{cloudId}/clusters/{clusterId} [get]
func (a *API) GetClusterHandler(c echo.Context) error {
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}

	clusterId := c.Param("clusterId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, clusterId, common.CodeInvalidParm, nil)
	}

	cloudSet := &model.CloudSet{}

	// Cloud 조회
	cloudTable, err := a.Db.GetCloud(cloudId)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if cloudTable == nil {
		return response.ErrorfReqRes(c, cloudTable, common.DatabaseFalseData, err)
	}
	cloudTable.ToSet(cloudSet)

	// Cluster 조회
	clusters, err := a.Db.GetClusters(cloudId)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if len(clusters) == 0 {
		return response.ErrorfReqRes(c, clusters, common.DatabaseFalseData, err)
	}

	clusters[0].ToSet(cloudSet)

	// Node 조회
	nodes, err := a.Db.GetNodes(cloudId, *clusters[0].ClusterUid)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if nodes == nil {
		return response.ErrorfReqRes(c, nodes, common.DatabaseFalseData, err)
	}

	cloudSet.Nodes = &model.NodesInfo{}
	cloudSet.Nodes.FromTable(clusters[0], nodes)

	return response.Write(c, nil, cloudSet)

}

// SetClusterHandler - 클러스터 추가 (Openstack)
// @Tags        Openstack-Cluster
// @Summary     SetCluster
// @Description 클러스터 추가 (Openstack)
// @ID          SetCluster
// @Produce     json
// @Param       cloudId             path     string                    true "Cloud ID"
// @Param       OpenstackClusterSet body     model.OpenstackClusterSet true "Openstack Cluster Info"
// @Success     200                 {object} response.ReturnData
// @Router      /clouds/{cloudId}/clusters [post]
func (a *API) SetClusterHandler(c echo.Context) error {
	// TODO: 로그인 사용자 정보 활용 방법은?
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}

	var clusterSet model.OpenstackClusterSet
	err := getRequestData(c.Request(), &clusterSet)
	if err != nil {
		return response.ErrorfReqRes(c, clusterSet, common.CodeInvalidData, err)
	}

	clusterTable, nodeSetTables := clusterSet.ToTable(cloudId, false, "system", time.Now())

	// Openstack Cluster 정보 저장

	// Start. Transaction 얻어옴
	txdb, err := a.Db.BeginTransaction()
	if err != nil {
		return response.ErrorfReqRes(c, clusterSet, common.CodeFailedDatabase, err)
	}

	// Cluster 등록
	err = txdb.InsertOpenstackCluster(clusterTable)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, clusterTable, common.CodeFailedDatabase, err)
	}

	// NodeSet 등록
	for _, nodeSetTable := range nodeSetTables {
		err = txdb.InsertNodeSet(nodeSetTable)
		if err != nil {
			txErr := txdb.Rollback()
			if txErr != nil {
				logger.Info("DB rollback Failed.", txErr)
			}
			return response.ErrorfReqRes(c, nodeSetTable, common.CodeFailedDatabase, err)
		}
	}

	txErr := txdb.Commit()
	if txErr != nil {
		logger.Info("DB commit Failed.", txErr)
	}

	// Openstack Cluster 생성 (CAPI with Openstack provider)

	return response.WriteWithCode(c, clusterSet, common.OpenstackClusterRegisteredAndProvisioning, nil)
}

// UpdateClusterHandler - 클러스터 수정 (Openstack)
// @Tags        Openstack-Cluster
// @Summary     UpdateCluster
// @Description 클러스터 수정 (Openstack)
// @ID          UpdateCluster
// @Produce     json
// @Param       cloudId             path     string                    true "Cloud ID"
// @Param       clusterId           path     string                    true "Cluster ID"
// @Param       OpenstackClusterSet body     model.OpenstackClusterSet true "Openstack Cluster Info"
// @Success     200                 {object} response.ReturnData
// @Router      /clouds/{cloudId}/clusters/{clusterId} [put]
func (a *API) UpdateClusterHandler(c echo.Context) error {
	return nil

}

// DeleteClusterHandler - 클러스터 삭제 (Openstack)
// @Tags        Openstack-Cluster
// @Summary     DeleteCluster
// @Description 클러스터 삭제 (Openstack)
// @ID          DeleteCluster
// @Produce     json
// @Param       cloudId path     string true "Cloud ID"
// @Success     200     {object} response.ReturnData
// @Router      /clouds/{cloudId}/clusters/{clusterId} [delete]
func (a *API) DeleteClusterHandler(c echo.Context) error {
	// TODO: Delete Cluster and resources
	return nil
}

// ProvisioningClusterHandler - 클러스터 Provisioning (Openstack)
// @Tags        Openstack-Cluster
// @Summary     ProvisioningCluster
// @Description 저장된 클러스터 정보를 이용해서 Provision 처리 (Openstack)
// @ID          ProvisioningCluster
// @Produce     json
// @Param       cloudId   path     string true "Cloud ID"
// @Param       clusterId path     string true "Cluster ID"
// @Success     200       {object} response.ReturnData
// @Router      /clouds/{cloudId}/clusters/{clusterId} [post]
func (a *API) ProvisioningClusterHandler(c echo.Context) error {
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}

	clusterId := c.Param("clusterId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, clusterId, common.CodeInvalidParm, nil)
	}

	// Cluster 정보 조회
	clusterTable, err := a.Db.GetOpenstackCluster(cloudId, clusterId)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if clusterTable == nil {
		return response.ErrorfReqRes(c, clusterTable, common.DatabaseFalseData, err)
	}

	// Cluster 상태 검증
	if *clusterTable.Status != 1 {
		return response.ErrorfReqRes(c, clusterTable, common.ProvisioningOnlySaved, err)
	}

	// Kubernetes 버전 조회
	codeTable, err := a.Db.GetCode("K8sVersions", *clusterTable.Version)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if codeTable == nil {
		return response.ErrorfReqRes(c, clusterTable, common.DatabaseFalseData, err)
	}

	// NodeSets 정보 조회
	nodeSetTables, err := a.Db.GetNodeSets(clusterId)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if len(nodeSetTables) == 0 {
		return response.ErrorfReqRes(c, nodeSetTables, common.DatabaseFalseData, err)
	}

	// Provisioning (background)
	err = ProvisioningOpenstackCluster(clusterTable, nodeSetTables, *codeTable.Name)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.ProvisioningFailed, err)
	}

	// TODO: Workload Cluster 관련 후처리 작업
	//       - kubeconfig 설정
	//       - check cluster provisioned or failed
	//       - machine crated (controlplane, workers)
	// TODO: Provisioning 종료를 확인하는 방법은? Webhook, readyness??

	// // Get Pod List
	// podList, err := kubemethod.GetPodList("os-cluster", "")
	// if err != nil {
	// 	return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	// }

	// // Get Node List
	// nodeList, err := kubemethod.GetNodeList("os-cluster")
	// if err != nil {
	// 	return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	// }

	// // Get Kubeconfig
	// data, err := kubemethod.GetKubeconfig("default", "os-cluster-kubeconfig", "value")
	// if err != nil {
	// 	return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	// }

	// // Add cluster's kubeconfig
	// err = config.HostCluster.Add([]byte(data))
	// if err != nil {
	// 	return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	// }

	// // Remove cluster's kubeconfig
	// err := config.HostCluster.Remove("os-cluster")
	// if err != nil {
	// 	return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	// }

	//return response.WriteWithCode(c, nil, common.OpenstackClusterProvisioning, data)
	//return response.WriteWithCode(c, nil, common.OpenstackClusterProvisioning, podList)
	return response.WriteWithCode(c, nil, common.OpenstackClusterProvisioning, nil)
}
