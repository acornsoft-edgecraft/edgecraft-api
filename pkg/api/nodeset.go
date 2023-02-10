/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package api

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/kubemethod"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/response"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/db"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
	"github.com/labstack/echo/v4"
)

/*******************************
 ** NodeSet for Openstack
 *******************************/

// GetNodeSetListHandler - 클래스터의 NodeSet 리스트 조회 (Openstack)
// @Tags        Openstack-Cluster-NodeSet
// @Summary     GetNodeSetList
// @Description 클러스터의 NodeSet 리스트 (Openstack)
// @ID          GetNodeSetList
// @Produce     json
// @Param       cloudId 	path     string true "Cloud ID"
// @Param       clusterId 	path     string true "Cluster ID"
// @Success     200     {object} response.ReturnData
// @Router      /clouds/{cloudId}/clusters/{clusterId}/nodesets [get]
func (a *API) GetNodeSetListHandler(c echo.Context) error {
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}

	clusterId := c.Param("clusterId")
	if clusterId == "" {
		return response.ErrorfReqRes(c, clusterId, common.CodeInvalidParm, nil)
	}

	// 클러스터 검증
	clusterTable, err := a.Db.GetOpenstackCluster(cloudId, clusterId)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	}
	if clusterTable == nil {
		return response.ErrorfReqRes(c, clusterTable, common.ClusterNotFound, err)
	}

	// NodeSet 리스트 조회
	nodeSetTables, err := a.Db.GetNodeSets(clusterId)
	if err != nil {
		return response.Errorf(c, common.CodeFailedDatabase, err)
	}
	if len(nodeSetTables) == 0 {
		return response.Errorf(c, common.NodeSetNotFound, err)
	}

	// 클러스터 DB정보를 NodeSetInfo 정보로 설정
	nodeSetInfo := model.OpenstackNodeSetInfo{}
	nodeSetInfo.FromTable(clusterTable, nodeSetTables)

	// Provisioned 상태인 경우는 Node정보 구성
	k8sFailed := false
	if *clusterTable.Status == common.StatusProvisioned {
		k8sFailed = kubemethod.ArrangeK8SNodesToNodeSetInfo(*clusterTable.Name, nodeSetInfo)
	}

	if k8sFailed {
		return response.WriteWithCode(c, nil, common.KubernetesNotYet, nodeSetInfo)
	} else {
		return response.Write(c, nil, nodeSetInfo)
	}
}

// GetNodeSetHandler - 클래스터의 NodeSet 상세 조회 (Openstack)
// @Tags        Openstack-Cluster-NodeSet
// @Summary     GetNodeSet
// @Description 클러스터의 NodeSet 상세 (Openstack)
// @ID          GetNodeSet
// @Produce     json
// @Param       cloudId 	path     string true "Cloud ID"
// @Param       clusterId 	path     string true "Cluster ID"
// @Param       nodeSetId 	path     string true "NodeSet ID"
// @Success     200     {object} response.ReturnData
// @Router      /clouds/{cloudId}/clusters/{clusterId}/nodesets/{nodeSetId} [get]
func (a *API) GetNodeSetHandler(c echo.Context) error {
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}

	clusterId := c.Param("clusterId")
	if clusterId == "" {
		return response.ErrorfReqRes(c, clusterId, common.CodeInvalidParm, nil)
	}

	nodeSetId := c.Param("nodeSetId")
	if nodeSetId == "" {
		return response.ErrorfReqRes(c, nodeSetId, common.CodeInvalidParm, nil)
	}

	// 클러스터 검증
	clusterTable, err := a.Db.GetOpenstackCluster(cloudId, clusterId)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	}
	if clusterTable == nil {
		return response.ErrorfReqRes(c, clusterTable, common.ClusterNotFound, err)
	}

	// NodeSet 조회
	nodeSetTable, err := a.Db.GetNodeSet(clusterId, nodeSetId)
	if err != nil {
		return response.Errorf(c, common.CodeFailedDatabase, err)
	}
	if nodeSetTable == nil {
		return response.Errorf(c, common.NodeSetNotFound, err)
	}

	// 클러스터 DB정보를 NodeSetInfo 정보로 설정
	nodeSetInfo := model.NodeSetInfo{}
	nodeSetInfo.FromTable(nodeSetTable)

	// Provisioned 상태인 경우는 Node정보 구성
	k8sFailed := false
	if *clusterTable.Status == common.StatusProvisioned {
		nodeList, err := kubemethod.GetNodeList(*clusterTable.Name)
		if err != nil {
			k8sFailed = true
		} else {
			for _, node := range nodeList {
				if strings.Contains(node.Name, "-"+nodeSetInfo.Name+"-") {
					nodeSetInfo.Nodes = append(nodeSetInfo.Nodes, node)
				}
			}
		}
	}

	if k8sFailed {
		return response.WriteWithCode(c, nil, common.KubernetesNotYet, nodeSetInfo)
	} else {
		return response.Write(c, nil, nodeSetInfo)
	}
}

// SetNodeSetHandler - 클래스터의 NodeSet 추가 (Openstack)
// @Tags        Openstack-Cluster-NodeSet
// @Summary     SetNodeSet
// @Description 클러스터에 NodeSet 추가 (Openstack)
// @ID          SetNodeSet
// @Produce     json
// @Param       cloudId 	path     string true "Cloud ID"
// @Param       clusterId 	path     string true "Cluster ID"
// @Param       NodeSet		body     model.NodeSetInfo true "NodeSet Info"
// @Success     200     {object} response.ReturnData
// @Router      /clouds/{cloudId}/clusters/{clusterId}/nodesets [post]
func (a *API) SetNodeSetHandler(c echo.Context) error {
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}

	clusterId := c.Param("clusterId")
	if clusterId == "" {
		return response.ErrorfReqRes(c, clusterId, common.CodeInvalidParm, nil)
	}

	nodeSet := model.NodeSetInfo{}
	err := getRequestData(c.Request(), &nodeSet)
	if err != nil {
		return response.ErrorfReqRes(c, nodeSet, common.CodeInvalidData, err)
	}

	// 클러스터 검증
	clusterTable, err := a.Db.GetOpenstackCluster(cloudId, clusterId)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	}
	if clusterTable == nil {
		return response.ErrorfReqRes(c, clusterTable, common.ClusterNotFound, err)
	}

	// NodeSet 정보 구성
	nodeSetTable := &model.NodeSetTable{}
	nodeSet.ToTable(nodeSetTable, false, "system", time.Now())
	nodeSetTable.ClusterUid = clusterTable.ClusterUid
	nodeSetTable.Type = utils.IntPrt(common.NodeTypeWorker)

	// 트랜잭션 구간 처리
	provisioningFailed := false
	err = a.Db.TransactionScope(func(txDB db.DB) error {
		// NodeSet 추가
		err := txDB.InsertNodeSet(nodeSetTable)
		if err != nil {
			return err
		}

		// NodeSet Provisioning (apply)
		err = ProvisioningOpenstackNodeSet(a.Worker, a.Db, clusterTable, []*model.NodeSetTable{nodeSetTable}, a.getCodeNameByKey("K8sVersions", *clusterTable.Version))
		if err != nil {
			provisioningFailed = true
			logger.WithError(err).Info("NodeSet provisioning failed")
			return err
		}

		return nil
	})
	if err != nil {
		if provisioningFailed {
			return response.ErrorfReqRes(c, nodeSetTable, common.ProvisioningNodeSetFailed, err)
		} else {
			return response.ErrorfReqRes(c, nodeSetTable, common.CodeFailedDatabase, err)
		}
	}

	return response.WriteWithCode(c, nil, common.OpenstackClusterNodeSetProvisioning, nil)
}

// // TODO: Update NodeSet

// // UpdateNodeSetHandler - 클래스터의 NodeSet 갱신 (Openstack)
// // @Tags        Openstack-Cluster-NodeSet
// // @Summary     UpdateNodeSet
// // @Description 클러스터에 NodeSet 갱신 (Openstack)
// // @ID          UpdateNodeSet
// // @Produce     json
// // @Param       cloudId 	path     string true "Cloud ID"
// // @Param       clusterId 	path     string true "Cluster ID"
// // @Param       NodeSet		body     model.NodeSetInfo true "NodeSet Info"
// // @Success     200     {object} response.ReturnData
// // @Router      /clouds/{cloudId}/clusters/{clusterId}/nodesets [put]
// func (a *API) UpdateNodeSetHandler(c echo.Context) error {
// 	return nil
// }

// TODO: Delete NodeSet

// DeleteNodeSetHandler - 클래스터의 NodeSet 삭제 (Openstack)
// @Tags        Openstack-Cluster-NodeSet
// @Summary     DeleteNodeSet
// @Description 클러스터의 NodeSet 삭제 (Openstack)
// @ID          DeleteNodeSet
// @Produce     json
// @Param       cloudId 	path     string true "Cloud ID"
// @Param       clusterId 	path     string true "Cluster ID"
// @Param       nodeSetId 	path     string true "NodeSet ID"
// @Success     200     {object} response.ReturnData
// @Router      /clouds/{cloudId}/clusters/{clusterId}/nodesets/{nodeSetId} [delete]
func (a *API) DeleteNodeSetHandler(c echo.Context) error {
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}

	clusterId := c.Param("clusterId")
	if clusterId == "" {
		return response.ErrorfReqRes(c, clusterId, common.CodeInvalidParm, nil)
	}

	nodeSetId := c.Param("nodeSetId")
	if nodeSetId == "" {
		return response.ErrorfReqRes(c, nodeSetId, common.CodeInvalidParm, nil)
	}

	// 클러스터 검증
	clusterTable, err := a.Db.GetOpenstackCluster(cloudId, clusterId)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	}
	if clusterTable == nil {
		return response.ErrorfReqRes(c, clusterTable, common.ClusterNotFound, err)
	}

	// NodeSet 조회
	nodeSetTable, err := a.Db.GetNodeSet(clusterId, nodeSetId)
	if err != nil {
		return response.Errorf(c, common.CodeFailedDatabase, err)
	}
	if nodeSetTable == nil {
		return response.Errorf(c, common.NodeSetNotFound, err)
	}

	// 트랜잭션 구간 처리
	k8sFailed := false
	err = a.Db.TransactionScope(func(txDB db.DB) error {
		// NodeSet 삭제
		affectedRows, err := txDB.DeleteNodeSet(clusterId, nodeSetId)
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New("cannot find noddeset for deleting")
		}

		err = kubemethod.RemoveNodeSet(*clusterTable.Name, *nodeSetTable.Name, *clusterTable.Namespace)
		if err != nil {
			k8sFailed = true
			logger.WithError(err).Info("Provioned NodeSet delete failed")
			return err
		}

		return nil
	})
	if err != nil {
		if k8sFailed {
			return response.ErrorfReqRes(c, nodeSetTable, common.ProvisionedNodeSetDeleteFailed, err)
		} else {
			return response.ErrorfReqRes(c, nodeSetTable, common.CodeFailedDatabase, err)
		}
	}

	return response.WriteWithCode(c, nil, common.OpenstackClusterNodeSetDeleting, nil)
}

// UpdateNodeCountHandler - 클래스터의 NodeSet의 NodeCount 갱신 (Openstack)
// @Tags        Openstack-Cluster-NodeSet
// @Summary     UpdateNodeCount
// @Description 클러스터의 NodeSet에 NodeCount 갱신 (Openstack)
// @ID          UpdateNodeCount
// @Produce     json
// @Param       cloudId 	path     string true "Cloud ID"
// @Param       clusterId 	path     string true "Cluster ID"
// @Param       nodeSetId 	path     string true "NodeSet ID"
// @Param       nodeCount 	path     string true "NodeCount"
// @Success     200     {object} response.ReturnData
// @Router      /clouds/{cloudId}/clusters/{clusterId}/nodesets/{nodeSetId}/{nodeCount} [get]
func (a *API) UpdateNodeCountHandler(c echo.Context) error {
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}

	clusterId := c.Param("clusterId")
	if clusterId == "" {
		return response.ErrorfReqRes(c, clusterId, common.CodeInvalidParm, nil)
	}

	nodeSetId := c.Param("nodeSetId")
	if nodeSetId == "" {
		return response.ErrorfReqRes(c, nodeSetId, common.CodeInvalidParm, nil)
	}

	nodeCountParam := c.Param("nodeCount")
	if nodeCountParam == "" {
		return response.ErrorfReqRes(c, nodeCountParam, common.CodeInvalidParm, nil)
	}

	nodeCount, err := strconv.Atoi(nodeCountParam)
	if err != nil {
		return response.ErrorfReqRes(c, nodeCountParam, common.CodeInvalidParm, err)
	} else if nodeCount <= 0 {
		return response.ErrorfReqRes(c, nodeCountParam, common.CodeInvalidParm, errors.New("nodeCount must be great than zero"))
	}

	// 클러스터 검증
	clusterTable, err := a.Db.GetOpenstackCluster(cloudId, clusterId)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	}
	if clusterTable == nil {
		return response.ErrorfReqRes(c, clusterTable, common.ClusterNotFound, err)
	}

	// NodeSet 조회
	nodeSetTable, err := a.Db.GetNodeSet(clusterId, nodeSetId)
	if err != nil {
		return response.Errorf(c, common.CodeFailedDatabase, err)
	}
	if nodeSetTable == nil {
		return response.Errorf(c, common.NodeSetNotFound, err)
	}

	// NodeSet에 NodeCount 반영
	nodeSetTable.NodeCount = utils.IntPrt(nodeCount)

	// 트랜잭션 구간 처리
	k8sFailed := false
	err = a.Db.TransactionScope(func(txDB db.DB) error {
		// NodeSet 갱신
		affectedRows, err := txDB.UpdateNodeSet(nodeSetTable)
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New("cannot find noddeset for nodecount updating")
		}

		// Node count 변경
		err = kubemethod.UpdateNodeCount(*clusterTable.Name, *nodeSetTable.Name, *clusterTable.Namespace, *clusterTable.BootstrapProvider, *nodeSetTable.Type, nodeCount)
		if err != nil {
			k8sFailed = true
			return err
		}

		return nil
	})
	if err != nil {
		if k8sFailed {
			return response.ErrorfReqRes(c, nodeSetTable, common.ProvisionNodeCountChangeFailed, err)
		} else {
			return response.ErrorfReqRes(c, nodeSetTable, common.CodeFailedDatabase, err)
		}
	}

	return response.WriteWithCode(c, nil, common.NodeCountUpdated, nil)
}
