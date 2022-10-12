/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package api

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/response"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	"github.com/labstack/echo/v4"
)

/*******************************
 ** Cluster for Openstack
 *******************************/

// GetClusterListHandler - 전체 클러스터 리스트 (Openstack)
// @Tags Openstack-Cluster
// @Summary GetClusterList
// @Description 전체 클러스터 리스트 (Openstack)
// @ID GetClusterList
// @Produce json
// @Param cloudId path string true "Cloud ID"
// @Success 200 {object} response.ReturnData
// @Router /clouds/{cloudId}/clusters [get]
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
// @Tags Openstack-Cluster
// @Summary GetCluster
// @Description 클러스터 상세 조회 (Openstack)
// @ID GetCluster
// @Produce json
// @Param cloudId path string true "Cloud ID"
// @Param clusterId path string true "Cluster ID"
// @Success 200 {object} response.ReturnData
// @Router /clouds/{cloudId}/clusters/{clusterId} [get]
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
// @Tags Openstack-Cluster
// @Summary SetCluster
// @Description 클러스터 추가 (Openstack)
// @ID SetCluster
// @Produce json
// @Param cloudId path string true "Cloud ID"
// @Param clusterId path string true "Cluster ID"
// @Param OSClusterInfo body model.OSClusterInfo true "Openstack Cluster Info"
// @Success 200 {object} response.ReturnData
// @Router /clouds/{cloudId}/clusters [post]
func (a *API) SetClusterHandler(c echo.Context) error {
	return nil

}

// UpdateClusterHandler - 클러스터 수정 (Openstack)
// @Tags Openstack-Cluster
// @Summary UpdateCluster
// @Description 클러스터 수정 (Openstack)
// @ID UpdateCluster
// @Produce json
// @Param cloudId path string true "Cloud ID"
// @Param clusterId path string true "Cluster ID"
// @Param OSClusterInfo body model.OSClusterInfo true "Openstack Cluster Info"
// @Success 200 {object} response.ReturnData
// @Router /clouds/{cloudId}/clusters/{clusterId} [put]
func (a *API) UpdateClusterHandler(c echo.Context) error {
	return nil

}

// DeleteClusterHandler - 클러스터 삭제 (Openstack)
// @Tags Openstack-Cluster
// @Summary DeleteCluster
// @Description 클러스터 삭제 (Openstack)
// @ID DeleteCluster
// @Produce json
// @Param cloudId path string true "Cloud ID"
// @Success 200 {object} response.ReturnData
// @Router /clouds/{cloudId}/clusters/{clusterId} [delete]
func (a *API) DeleteClusterHandler(c echo.Context) error {
	return nil
}