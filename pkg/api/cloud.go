package api

import (
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/response"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"

	"github.com/labstack/echo/v4"
)

/*******************************
 ** Cloud
 *******************************/

// GetCloudListHandler - 전체 클라우드 리스트
// @Tags Cloud
// @Summary GetCloudList
// @Description Get all cloud list
// @ID GetCloudList
// @Produce json
// @Success 200 {object} response.ReturnData
// @Router /clouds [get]
func (a *API) GetCloudListHandler(c echo.Context) error {
	result, err := a.Db.GetClouds()
	if err != nil {
		return response.Errorf(c, common.CodeFailedDatabase, err)
	}
	return response.Write(c, nil, result)
}

// GetCloudHandler - 클라우드 상세 정보
// @Tags Cloud
// @Summary GetCloud
// @Description Get specific cloud
// @ID GetCloud
// @Produce json
// @Param cloudId path string true "cloudId"
// @Success 200 {object} response.ReturnData
// @Router /clouds/{cloudId} [get]
func (a *API) GetCloudHandler(c echo.Context) error {
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
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
// @Param cloudId path string true "cloudId"
// @Param cloudSet body model.CloudSet true "Cloud Set"
// @Success 200 {object} response.ReturnData
// @Router /clouds/{cloudId} [put]
func (a *API) UpdateCloudHandler(c echo.Context) error {
	// TODO: 로그인 사용자 정보 활용 방법은?
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}

	var cloudSet model.CloudSet
	var at time.Time = time.Now()

	err := getRequestData(c.Request(), &cloudSet)
	if err != nil {
		return response.ErrorfReqRes(c, cloudSet, common.CodeInvalidData, err)
	}

	// 공통 코드 조회
	codeTable, err := a.Db.GetCode("CloudStatus", 1)
	if err != nil || codeTable == nil {
		return response.ErrorfReqRes(c, cloudSet, common.CodeFailedDatabase, err)
	}

	// Cloud 조회
	cloudTable, err := a.Db.GetCloud(cloudId)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if cloudTable == nil {
		return response.ErrorfReqRes(c, cloudTable, common.DatabaseFalseData, err)
	} else if cloudTable.Status == codeTable.Code {
		// 저장 상태만 수정 가능
		return response.ErrorfReqRes(c, cloudTable, common.CloudCreated_CantUpdate, nil)
	}

	// Cluster 조회
	clusterTables, err := a.Db.GetClusters(cloudId)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if len(clusterTables) == 0 {
		return response.ErrorfReqRes(c, clusterTables, common.DatabaseCodeFalseData, err)
	}

	// Node 조회
	nodeTables, err := a.Db.GetNodes(cloudId, *clusterTables[0].ClusterUid)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if nodeTables == nil {
		return response.ErrorfReqRes(c, nodeTables, common.DatabaseFalseData, err)
	}

	newCloudTable, newClusterTable, newNodeTables := cloudSet.ToTable(true, "system", at)

	// Start. Transaction 얻어옴
	txdb, err := a.Db.BeginTransaction()
	if err != nil {
		return response.ErrorfReqRes(c, cloudSet, common.CodeFailedDatabase, err)
	}

	// Cloud 갱신
	newCloudTable.Creator = cloudTable.Creator
	newCloudTable.Created = cloudTable.Created
	cnt, err := txdb.UpdateCloud(newCloudTable)
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
	newClusterTable.Creator = clusterTables[0].Creator
	newClusterTable.Created = clusterTables[0].Created
	cnt, err = txdb.UpdateCluster(newClusterTable)
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
	cnt, err = txdb.DeleteNodes(cloudId, *clusterTables[0].ClusterUid)
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
	for _, nodeTable := range newNodeTables {
		for _, oldNodeTable := range nodeTables {
			if oldNodeTable.NodeUid == nodeTable.NodeUid {
				nodeTable.Creator = oldNodeTable.Creator
				nodeTable.Created = oldNodeTable.Created
				break
			}
		}

		if nodeTable.Creator == nil {
			nodeTable.Creator = utils.StringPtr("system")
			nodeTable.Created = utils.TimePtr(at)
		}
		nodeTable.Updater = utils.StringPtr("system")
		nodeTable.Updated = utils.TimePtr(at)

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
// @Param cloudId path string true "cloudId"
// @Success 200 {object} response.ReturnData
// @Router /clouds/{cloudId} [delete]
func (a *API) DeleteCloudHandler(c echo.Context) error {
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
 ** Cloud - Node
 *******************************/

// GetCloudNodeListHandler - 클라우드에 속한 노드 리스트 조회
// @Tags CloudNode
// @Summary GetCloudNodeList
// @Description 클라우드에 속한 노드 리스트 조회
// @ID GetCloudNodeList
// @Produce json
// @Param cloudId path string true "Cloud ID"
// @Success 200 {object} response.ReturnData
// @Router /clouds/{cloudId}/nodes [get]
func (a *API) GetCloudNodeListHandler(c echo.Context) error {
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}

	// Node 정보 조회
	list, err := a.Db.GetNodesByCloud(cloudId)
	if err != nil {
		return response.Errorf(c, common.CodeFailedDatabase, err)
	}
	if list == nil {
		return response.ErrorfReqRes(c, list, common.DatabaseFalseData, err)
	}

	// NodesInfo 구조체 작성
	var nodes *model.NodesInfo = &model.NodesInfo{}
	for _, node := range list {
		var nodeSpecInfo *model.NodeSpecificInfo = &model.NodeSpecificInfo{}
		nodeSpecInfo.FromTable(node)

		if *node.Type == 1 {
			nodes.MasterNodes = append(nodes.MasterNodes, nodeSpecInfo)
		} else {
			nodes.WorkerNodes = append(nodes.WorkerNodes, nodeSpecInfo)
		}
	}

	return response.Write(c, nil, nodes)
}

// GetCloudNodeHandler - 클라우드에 속한 노드 상세정보 조회
// @Tags CloudNode
// @Summary GetCloudNode
// @Description 클라우드에 속한 노드 상세정보 조회
// @ID GetCloudNode
// @Produce json
// @Param cloudId path string true "Cloud ID"
// @Param nodeId path string true "Node ID"
// @Success 200 {object} response.ReturnData
// @Router /clouds/{cloudId}/nodes/{nodeId} [get]
func (a *API) GetCloudNodeHandler(c echo.Context) error {
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}

	nodeId := c.Param("nodeId")
	if nodeId == "" {
		return response.ErrorfReqRes(c, nodeId, common.CodeInvalidParm, nil)
	}

	// Node 조회
	nodeTable, err := a.Db.GetNodeByCloud(cloudId, nodeId)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if nodeTable == nil {
		return response.ErrorfReqRes(c, nodeTable, common.DatabaseFalseData, err)
	}

	var node *model.NodeSpecificInfo = &model.NodeSpecificInfo{}
	node.FromTable(nodeTable)

	return response.Write(c, nil, node)
}

// SetCloudNodeHandler - 클라우드에 노드 등록
// @Tags CloudNode
// @Summary SetCloudNode
// @Description 클라우드에 노드 등록
// @ID SetCloudNode
// @Produce json
// @Param cloudId path string true "Cloud ID"
// @Param node body model.NodeSpecificInfo true "Node Specific Info"
// @Success 200 {object} response.ReturnData
// @Router /clouds/{cloudId}/nodes [post]
func (a *API) SetCloudNodeHandler(c echo.Context) error {
	// TODO: 로그인 사용자 정보 활용 방법은?

	// Cloud Id 수신
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}

	// Node 정보 수신
	var node *model.NodeSpecificInfo
	err := getRequestData(c.Request(), &node)
	if err != nil {
		return response.ErrorfReqRes(c, node, common.CodeInvalidData, err)
	}

	// Cluster 정보 조회
	clusters, err := a.Db.GetClusters(cloudId)
	if err != nil {
		return response.ErrorfReqRes(c, cloudId, common.CodeFailedDatabase, err)
	}
	if clusters == nil {
		return response.ErrorfReqRes(c, cloudId, common.DatabaseFalseData, err)
	}

	// Node 정보 설정
	var nodeTable *model.NodeTable = &model.NodeTable{}
	node.ToTable(nodeTable, false, "system", time.Now())
	nodeTable.CloudUid = &cloudId
	nodeTable.ClusterUid = clusters[0].ClusterUid
	nodeTable.Status = utils.IntPrt(1)

	// Start. Transaction 얻어옴
	txdb, err := a.Db.BeginTransaction()
	if err != nil {
		return response.ErrorfReqRes(c, node, common.CodeFailedDatabase, err)
	}

	// Node 등록
	err = txdb.InsertNode(nodeTable)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, node, common.CodeFailedDatabase, err)
	}

	txErr := txdb.Commit()
	if txErr != nil {
		logger.Info("DB commit Failed.", txErr)
	}

	return response.Write(c, node, nil)
}

// UpdateCloudNodeHandler - 클라우드의 노드 수정
// @Tags CloudNode
// @Summary UpdateCloudNode
// @Description 클라우드의 노드 수정
// @ID UpdateCloudNode
// @Produce json
// @Param cloudId path string true "Cloud ID"
// @Param nodeId path string true "Node ID"
// @Param node body model.NodeSpecificInfo true "Nodes Info"
// @Success 200 {object} response.ReturnData
// @Router /clouds/{cloudId}/nodes/{nodeId} [put]
func (a *API) UpdateCloudNodeHandler(c echo.Context) error {
	// TODO: 로그인 사용자 정보 활용 방법은?

	// Cloud Id 수신
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}

	// Node Id 수신
	nodeId := c.Param("nodeId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}

	// Node 정보 수신
	var node *model.NodeSpecificInfo
	err := getRequestData(c.Request(), &node)
	if err != nil {
		return response.ErrorfReqRes(c, node, common.CodeInvalidData, err)
	}

	// Cluster 정보 조회
	clusterTables, err := a.Db.GetClusters(cloudId)
	if err != nil {
		return response.ErrorfReqRes(c, cloudId, common.CodeFailedDatabase, err)
	}
	if clusterTables == nil {
		return response.ErrorfReqRes(c, cloudId, common.DatabaseFalseData, err)
	}

	// Node 정보 조회
	nodeTable, err := a.Db.GetNode(cloudId, *clusterTables[0].ClusterUid, nodeId)
	if err != nil {
		return response.ErrorfReqRes(c, nodeId, common.CodeFailedDatabase, err)
	}
	if nodeTable == nil {
		return response.ErrorfReqRes(c, nodeId, common.DatabaseFalseData, err)
	}

	// Node 정보 설정
	node.ToTable(nodeTable, true, "system", time.Now())

	// Start. Transaction 얻어옴
	txdb, err := a.Db.BeginTransaction()
	if err != nil {
		return response.ErrorfReqRes(c, node, common.CodeFailedDatabase, err)
	}

	// Node 수정
	count, err := txdb.UpdateNode(nodeTable)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, node, common.CodeFailedDatabase, err)
	}
	if count == 0 {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, node, common.DatabaseFalseData, err)
	}

	txErr := txdb.Commit()
	if txErr != nil {
		logger.Info("DB commit Failed.", txErr)
	}

	return response.Write(c, node, nil)
}

// DeleteCloudNodeHandler - 클라우드의 노드 삭제
// @Tags CloudNode
// @Summary DeleteCloudNode
// @Description 클라우드의 노드 삭제
// @ID DeleteCloudNode
// @Produce json
// @Param cloudId path string true "Cloud ID"
// @Param nodeId path string true "Node ID"
// @Success 200 {object} response.ReturnData
// @Router /clouds/{cloudId}/nodes/{nodeId} [delete]
func (a *API) DeleteCloudNodeHandler(c echo.Context) error {
	// Cloud Id 수신
	cloudId := c.Param("cloudId")
	if cloudId == "" {
		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
	}

	// Node Id 수신
	nodeId := c.Param("nodeId")
	if nodeId == "" {
		return response.ErrorfReqRes(c, nodeId, common.CodeInvalidParm, nil)
	}

	// Start. Transaction 얻어옴
	txdb, err := a.Db.BeginTransaction()
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	}

	// Node 삭제
	count, err := txdb.DeleteNode(nodeId)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, nodeId, common.CodeFailedDatabase, err)
	}
	if count == 0 {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, nodeId, common.DatabaseFalseData, err)
	}

	txErr := txdb.Commit()
	if txErr != nil {
		logger.Info("DB commit Failed.", txErr)
	}

	return response.Write(c, nodeId, nil)
}
