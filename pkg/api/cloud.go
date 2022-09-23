package api

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/response"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"

	//mr "github.com/acornsoft-edgecraft/edgecraft-api/pkg/model/response"

	"github.com/labstack/echo/v4"
)

// AllCloudListHandler - 전체 클라우드 리스트
// @Tags Cloud
// @Summary AllClooudList
// @Description Get all cloud list
// @ID AllCloudList
// @Produce json
// @Success 200 {object} response.ReturnData
// @Router /clouds [get]
func (a *API) AllCloudListHandler(c echo.Context) error {
	res, err := a.Db.GetAllCloud()
	if err != nil {
		return response.Errorf(c, common.CodeFailedDatabase, err)
	}
	return response.Write(c, nil, res)
}

// SelectCloudHandler - 클라우드 상세 정보
// @Tags Cloud
// @Summary SelectCloud
// @Description Get specific cloud
// @ID SelectCloud
// @Produce json
// @Param cloudUid path string true "cloudUid"
// @Success 200 {object} response.ReturnData
// @Router /clouds/{cloudUid} [get]
func (a *API) SelectCloudHandler(c echo.Context) error {
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
	nodes, err := a.Db.SelectNodes(cloudUid, clusters[0].ClusterUid)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if nodes == nil {
		return response.ErrorfReqRes(c, nodes, common.DatabaseFalseData, err)
	}

	cloudSet.Nodes = &model.NodesInfo{}
	cloudSet.Nodes.FromTable(clusters[0], nodes)

	// 	resCloud, err := a.Db.GetCloud(cloudUid)

	// 	resCloudCluster, err := a.Db.GetCloudCluster(cloudUid)
	// 	if err != nil {
	// 		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	// 	} else if resCloudCluster == nil {
	// 		return response.ErrorfReqRes(c, resCloudCluster, common.DatabaseFalseData, err)
	// 	}

	// 	// resCloudNode, err := a.Db.GetCloudNode(cloudUid, *resCloudCluster.CloudClusterUid)
	// 	// if err != nil {
	// 	// 	return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	// 	// } else if resCloud == nil {
	// 	// 	return response.ErrorfReqRes(c, resCloud, common.DatabaseFalseData, err)
	// 	// }

	// 	// resClusterK8s, err := a.Db.SelectK8sCloudCluster(cloudUid)
	// 	// if err != nil {
	// 	// 	return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	// 	// } else if resCloud == nil {
	// 	// 	return response.ErrorfReqRes(c, resCloud, common.DatabaseFalseData, err)
	// 	// }

	// 	// resClusterBaremetal, err := a.Db.SelectBaremetalCloudCluster(cloudUid)
	// 	// if err != nil {
	// 	// 	return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	// 	// } else if resCloud == nil {
	// 	// 	return response.ErrorfReqRes(c, resCloud, common.DatabaseFalseData, err)
	// 	// }

	// 	resCloudNode, err := a.Db.SelectNodeCloudCluster(cloudUid)
	// 	if err != nil {
	// 		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	// 	} else if resCloudNode == nil {
	// 		return response.ErrorfReqRes(c, resCloudNode, common.DatabaseFalseData, err)
	// 	}

	// 	resNodes, err := a.Db.SelectCloudNode(cloudUid, *resCloudCluster.CloudClusterUid)
	// 	if err != nil {
	// 		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	// 	} else if resNodes == nil {
	// 		return response.ErrorfReqRes(c, resNodes, common.DatabaseFalseData, err)
	// 	}

	// 	var masterNodes []model.MasterNode
	// 	var workerNodes []model.WorkerNode

	// 	for _, nodes := range resNodes {
	// 		if *nodes.CloudNodeType == "master" {
	// 			// masterBarematal := model.NodeBaremetal{
	// 			// 	CloudNodeHostName:             nodes.CloudNodeHostName,
	// 			// 	CloudNodeBmcAddress:           nodes.CloudNodeBmcAddress,
	// 			// 	CloudNodeMacAddress:           nodes.CloudNodeMacAddress,
	// 			// 	CloudNodeBootMode:             nodes.CloudNodeBootMode,
	// 			// 	CloudNodeOnlinePower:          nodes.CloudNodeOnlinePower,
	// 			// 	CloudNodeExternalProvisioning: nodes.CloudNodeExternalProvisioning,
	// 			// }

	// 			// nodes := model.Nodes{
	// 			// 	CloudNodeName:  nodes.CloudNodeName,
	// 			// 	CloudNodeIp:    nodes.CloudNodeIp,
	// 			// 	CloudNodeLabel: nodes.CloudNodeLabel,
	// 			// }

	// 			// masterNode := model.MasterNode{
	// 			// 	Baremetal: masterBarematal,
	// 			// 	Node:      nodes,
	// 			// }

	// 			var masterNode model.MasterNode
	// 			if err := utils.SetFieldsInStruct(&masterNode, &nodes); err != nil {
	// 				logger.Errorf("SetFields In Struct ERROR : %s", err)
	// 			}
	// 			masterNodes = append(masterNodes, masterNode)
	// 		} else if *nodes.CloudNodeType == "worker" {
	// 			// workerBarematal := model.NodeBaremetal{
	// 			// 	CloudNodeHostName:             nodes.CloudNodeHostName,
	// 			// 	CloudNodeBmcAddress:           nodes.CloudNodeBmcAddress,
	// 			// 	CloudNodeMacAddress:           nodes.CloudNodeMacAddress,
	// 			// 	CloudNodeBootMode:             nodes.CloudNodeBootMode,
	// 			// 	CloudNodeOnlinePower:          nodes.CloudNodeOnlinePower,
	// 			// 	CloudNodeExternalProvisioning: nodes.CloudNodeExternalProvisioning,
	// 			// }

	// 			// nodes := model.Nodes{
	// 			// 	CloudNodeName:  nodes.CloudNodeName,
	// 			// 	CloudNodeIp:    nodes.CloudNodeIp,
	// 			// 	CloudNodeLabel: nodes.CloudNodeLabel,
	// 			// }

	// 			// workerNode := model.WorkerNode{
	// 			// 	Baremetal: workerBarematal,
	// 			// 	Node:      nodes,
	// 			// }

	// 			var workerNode model.WorkerNode
	// 			if err := utils.SetFieldsInStruct(&workerNode, &nodes); err != nil {
	// 				logger.Errorf("SetFields In Struct ERROR : %s", err)
	// 			}
	// 			workerNodes = append(workerNodes, workerNode)
	// 		}

	//}

	// 	var res mr.RegisterCloud

	// 	res.Cloud = *resCloud
	// 	// res.Cluster.K8s = *resClusterK8s
	// 	// res.Cluster.Baremetal = *resClusterBaremetal

	// 	// Set Cluster fields
	// 	if err := utils.SetFieldsInStruct(&res.Cluster, resCloudCluster); err != nil {
	// 		logger.Errorf("SetFields In Struct ERROR : %s", err)
	// 	}

	// 	// Set Nodes fields
	// 	if err := utils.SetFieldsInStruct(&res.Nodes, resCloudNode); err != nil {
	// 		logger.Errorf("SetFields In Struct ERROR : %s", err)
	// 	}
	// 	// res.Nodes.CloudClusterLoadbalancerUse = resCloudNode.CloudClusterLoadbalancerUse
	// 	// res.Nodes.CloudClusterLoadbalancerAddress = resCloudNode.CloudClusterLoadbalancerAddress
	// 	// res.Nodes.CloudClusterLoadbalancerPort = resCloudNode.CloudClusterLoadbalancerPort
	// 	res.Nodes.MasterNode = masterNodes
	// 	res.Nodes.WorkerNode = workerNodes

	// 	// Set ETCD fields
	// 	if err := utils.SetFieldsInStruct(&res.EtcdStorage, resCloudCluster); err != nil {
	// 		logger.Errorf("SetFields In Struct ERROR : %s", err)
	// 	}
	// 	// res.EtcdStorage.Etcd.CloudClusterExternalEtcdUse = resCloudCluster.CloudClusterExternalEtcdUse
	// 	// res.EtcdStorage.Etcd.ExternalEtcdEndPoints = resCloudCluster.ExternalEtcdEndPoints
	// 	// res.EtcdStorage.Etcd.ExternalEtcdCertificateCa = resCloudCluster.ExternalEtcdCertificateCa
	// 	// res.EtcdStorage.Etcd.ExternalEtcdCertificateCert = resCloudCluster.ExternalEtcdCertificateCert
	// 	// res.EtcdStorage.Etcd.ExternalEtcdCertificateKey = resCloudCluster.ExternalEtcdCertificateKey
	// 	// res.EtcdStorage.CloudClusterStorageClass = resCloudCluster.CloudClusterStorageClass

	//return response.Write(c, nil, &res)
	return response.Write(c, nil, cloudSet)
}

// RegisterCloudHandler - 클라우드 등록
// @Tags Cloud
// @Summary RegisterCloud
// @Description Register cloud
// @ID RegisterCloud
// @Produce json
// @Param cloudSet body model.CloudSet true "Cloud Set"
// @Success 200 {object} response.ReturnData
// @Router /clouds [post]
func (a *API) RegisterCloudHandler(c echo.Context) error {
	// TODO: 로그인 사용자 정보 활용 방법은?
	var cloudSet model.CloudSet

	err := getRequestData(c.Request(), &cloudSet)
	if err != nil {
		return response.ErrorfReqRes(c, cloudSet, common.CodeInvalidData, err)
	}

	cloudTable, clusterTable, nodeTables := cloudSet.ToTable()

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
// @Param cloudInfo body model.RegisterCloud true "Cloud Info"
// @Success 200 {object} response.ReturnData
// @Router /clouds/{cloudUid} [put]
// func (a *API) UpdateCloudHandler(c echo.Context) error {
// 	// check param UID
// 	cloudUid, err := uuid.FromString(c.Param("cloudUid"))
// 	if err != nil {
// 		return response.ErrorfReqRes(c, cloudUid, common.CodeInvalidParm, err)
// 	}

// 	now := time.Now()

// 	var req mr.RegisterCloud
// 	err = getRequestData(c.Request(), &req)
// 	if err != nil {
// 		return response.ErrorfReqRes(c, req, common.CodeInvalidData, err)
// 	}

// 	// -- Service Logic
// 	// Start. Transaction 얻어옴
// 	txdb, err := a.Db.BeginTransaction()
// 	if err != nil {
// 		return response.ErrorfReqRes(c, req, common.CodeFailedDatabase, err)
// 	}

// 	// 1. Cloud 등록 업데이트
// 	var cloud model.Cloud
// 	getCloud, err := txdb.GetCloud(cloudUid)
// 	if err != nil {
// 		txErr := txdb.Rollback()
// 		if txErr != nil {
// 			logger.Info("DB Rollback Failed.", txErr)
// 		}
// 		return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
// 	} else if getCloud == nil {
// 		return response.ErrorfReqRes(c, getCloud, common.DatabaseFalseData, err)
// 	}

// 	// cloud.CloudName = getCloud.CloudName
// 	// Set fields in struct
// 	// if err := utils.SetFieldsInStruct(&cloud, &req.Cloud); err != nil {
// 	// 	logger.Errorf("SetFields In Struct ERROR : %s", err)
// 	// }

// 	user := "user"
// 	cloud = *getCloud
// 	cloud = req.Cloud

// 	cloud.UpdatedAt = &now
// 	cloud.Updater = &user
// 	count, err := txdb.UpdateCloud(&cloud)
// 	if err != nil {
// 		txErr := txdb.Rollback()
// 		if txErr != nil {
// 			logger.Info("DB Rollback Failed.", txErr)
// 		}
// 		return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
// 	}

// 	// 2. Cluster 등록 업데이트
// 	var cluster model.CloudCluster
// 	resCloudCluster, err := a.Db.GetCloudCluster(cloudUid)
// 	if err != nil {
// 		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
// 	} else if resCloudCluster == nil {
// 		return response.ErrorfReqRes(c, resCloudCluster, common.DatabaseFalseData, err)
// 	}

// 	cluster = *resCloudCluster

// 	cluster.CloudK8sVersion = req.Cluster.K8s.CloudK8sVersion
// 	cluster.CloudClusterPodCidr = req.Cluster.K8s.CloudClusterPodCidr
// 	cluster.CloudClusterServiceCidr = req.Cluster.K8s.CloudClusterServiceCidr
// 	cluster.CloudClusterBmcCredentialSecret = req.Cluster.Baremetal.CloudClusterBmcCredentialSecret
// 	cluster.CloudClusterBmcCredentialUser = req.Cluster.Baremetal.CloudClusterBmcCredentialUser
// 	cluster.CloudClusterBmcCredentialPassword = req.Cluster.Baremetal.CloudClusterBmcCredentialPassword
// 	cluster.CloudClusterImageUrl = req.Cluster.Baremetal.CloudClusterImageUrl
// 	cluster.CloudClusterImageChecksum = req.Cluster.Baremetal.CloudClusterImageChecksum
// 	cluster.CloudClusterImageChecksumType = req.Cluster.Baremetal.CloudClusterImageChecksumType
// 	cluster.CloudClusterImageFormat = req.Cluster.Baremetal.CloudClusterImageFormat
// 	cluster.CloudClusterMasterExtraConfig = req.Cluster.Baremetal.CloudClusterMasterExtraConfig
// 	cluster.CloudClusterWorkerExtraConfig = req.Cluster.Baremetal.CloudClusterWorkerExtraConfig
// 	cluster.CloudClusterLoadbalancerUse = req.Nodes.CloudClusterLoadbalancerUse
// 	cluster.CloudClusterLoadbalancerAddress = req.Nodes.CloudClusterLoadbalancerAddress
// 	cluster.CloudClusterLoadbalancerPort = req.Nodes.CloudClusterLoadbalancerPort
// 	cluster.CloudClusterExternalEtcdUse = req.EtcdStorage.Etcd.CloudClusterExternalEtcdUse
// 	cluster.ExternalEtcdEndPoints = req.EtcdStorage.Etcd.ExternalEtcdEndPoints
// 	cluster.ExternalEtcdCertificateCa = req.EtcdStorage.Etcd.ExternalEtcdCertificateCa
// 	cluster.ExternalEtcdCertificateCert = req.EtcdStorage.Etcd.ExternalEtcdCertificateCert
// 	cluster.ExternalEtcdCertificateKey = req.EtcdStorage.Etcd.ExternalEtcdCertificateKey
// 	cluster.CloudClusterStorageClass = req.EtcdStorage.CloudClusterStorageClass

// 	cluster.Updater = &user
// 	cluster.UpdatedAt = &now
// 	cluster.CloudUid = cloud.CloudUID

// 	count, err = txdb.UpdateCloudCluster(&cluster)
// 	if err != nil {
// 		txErr := txdb.Rollback()
// 		if txErr != nil {
// 			logger.Info("DB Rollback Failed.", txErr)
// 		}
// 		return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
// 	}

// 	// 3. Node 등록 업데이트
// 	var node model.CloudNode
// 	resCloudNode, err := a.Db.SelectNodeCloudCluster(cloudUid)
// 	if err != nil {
// 		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
// 	} else if resCloudNode == nil {
// 		return response.ErrorfReqRes(c, resCloudNode, common.DatabaseFalseData, err)
// 	}

// 	resNodes, err := a.Db.SelectCloudNode(cloudUid, *resCloudCluster.CloudClusterUid)
// 	if err != nil {
// 		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
// 	} else if resNodes == nil {
// 		return response.ErrorfReqRes(c, resNodes, common.DatabaseFalseData, err)
// 	}

// 	nodeType := ""

// 	var resNodeUid []interface{}
// 	for _, nodes := range resNodes {
// 		resNodeUid = append(resNodeUid, nodes.CloudNodeUid)
// 	}

// 	resRemoveNodes := hashset.New()
// 	resRemoveNodes.Add(resNodeUid...)

// 	setRemoveNodes := hashset.New()
// 	var updateNodes []*model.CloudNode
// 	// - MasterNode 등록 업데이트
// 	for _, master := range req.Nodes.MasterNode {
// 		for _, field := range resNodes {
// 			if master.CloudNodeUid != nil && *master.CloudNodeUid == *field.CloudNodeUid {
// 				// nodeType = "master"
// 				// field.CloudNodeType = &nodeType
// 				field.CloudNodeHostName = master.Baremetal.CloudNodeHostName
// 				field.CloudNodeBmcAddress = master.Baremetal.CloudNodeBmcAddress
// 				field.CloudNodeMacAddress = master.Baremetal.CloudNodeMacAddress
// 				field.CloudNodeBootMode = master.Baremetal.CloudNodeBootMode
// 				field.CloudNodeOnlinePower = master.Baremetal.CloudNodeOnlinePower
// 				field.CloudNodeExternalProvisioning = master.Baremetal.CloudNodeExternalProvisioning
// 				field.CloudNodeName = master.Node.CloudNodeName
// 				field.CloudNodeIp = master.Node.CloudNodeIp
// 				field.CloudNodeLabel = master.Node.CloudNodeLabel

// 				field.Updater = &user
// 				field.UpdatedAt = &now

// 				updateNodes = append(updateNodes, &field)
// 				// count, err = txdb.UpdateCloudNode(&field)
// 				// if err != nil {
// 				// 	txErr := txdb.Rollback()
// 				// 	if txErr != nil {
// 				// 		logger.Info("DB Rollback Failed.", txErr)
// 				// 	}
// 				// 	return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
// 				// }
// 				setRemoveNodes.Add(field.CloudNodeUid)
// 			}
// 		}
// 		// New Node
// 		if master.CloudNodeUid == nil {
// 			nodeType = "master"
// 			node.CloudNodeType = &nodeType
// 			node.CloudNodeHostName = master.Baremetal.CloudNodeHostName
// 			node.CloudNodeBmcAddress = master.Baremetal.CloudNodeBmcAddress
// 			node.CloudNodeMacAddress = master.Baremetal.CloudNodeMacAddress
// 			node.CloudNodeBootMode = master.Baremetal.CloudNodeBootMode
// 			node.CloudNodeOnlinePower = master.Baremetal.CloudNodeOnlinePower
// 			node.CloudNodeExternalProvisioning = master.Baremetal.CloudNodeExternalProvisioning
// 			node.CloudNodeName = master.Node.CloudNodeName
// 			node.CloudNodeIp = master.Node.CloudNodeIp
// 			node.CloudNodeLabel = master.Node.CloudNodeLabel

// 			node.Creator = &user
// 			node.CreatedAt = &now
// 			node.CloudUid = cloud.CloudUID
// 			node.CloudClusterUid = cluster.CloudClusterUid

// 			err = txdb.CreateCloudNode(&node)
// 			if err != nil {
// 				txErr := txdb.Rollback()
// 				if txErr != nil {
// 					logger.Info("DB Rollback Failed.", txErr)
// 				}
// 				return response.ErrorfReqRes(c, req, common.CodeFailedDatabase, err)
// 			}
// 			setRemoveNodes.Add(node.CloudNodeUid)
// 		}
// 	}

// 	// - WorkerrNode 등록 업데이트
// 	for _, worker := range req.Nodes.WorkerNode {
// 		for _, field := range resNodes {
// 			if worker.CloudNodeUid != nil && *worker.CloudNodeUid == *field.CloudNodeUid {
// 				// nodeType = "worker"
// 				// field.CloudNodeType = &nodeType
// 				field.CloudNodeHostName = worker.Baremetal.CloudNodeHostName
// 				field.CloudNodeBmcAddress = worker.Baremetal.CloudNodeBmcAddress
// 				field.CloudNodeMacAddress = worker.Baremetal.CloudNodeMacAddress
// 				field.CloudNodeBootMode = worker.Baremetal.CloudNodeBootMode
// 				field.CloudNodeOnlinePower = worker.Baremetal.CloudNodeOnlinePower
// 				field.CloudNodeExternalProvisioning = worker.Baremetal.CloudNodeExternalProvisioning
// 				field.CloudNodeName = worker.Node.CloudNodeName
// 				field.CloudNodeIp = worker.Node.CloudNodeIp
// 				field.CloudNodeLabel = worker.Node.CloudNodeLabel

// 				field.Updater = &user
// 				field.UpdatedAt = &now

// 				updateNodes = append(updateNodes, &field)
// 				// count, err = txdb.UpdateCloudNode(&field)
// 				// if err != nil {
// 				// 	txErr := txdb.Rollback()
// 				// 	if txErr != nil {
// 				// 		logger.Info("DB Rollback Failed.", txErr)
// 				// 	}
// 				// 	return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
// 				// }
// 				setRemoveNodes.Add(field.CloudNodeUid)
// 			}
// 		}
// 		// New Node
// 		if worker.CloudNodeUid == nil {
// 			nodeType = "worker"
// 			node.CloudNodeType = &nodeType
// 			node.CloudNodeHostName = worker.Baremetal.CloudNodeHostName
// 			node.CloudNodeBmcAddress = worker.Baremetal.CloudNodeBmcAddress
// 			node.CloudNodeMacAddress = worker.Baremetal.CloudNodeMacAddress
// 			node.CloudNodeBootMode = worker.Baremetal.CloudNodeBootMode
// 			node.CloudNodeOnlinePower = worker.Baremetal.CloudNodeOnlinePower
// 			node.CloudNodeExternalProvisioning = worker.Baremetal.CloudNodeExternalProvisioning
// 			node.CloudNodeName = worker.Node.CloudNodeName
// 			node.CloudNodeIp = worker.Node.CloudNodeIp
// 			node.CloudNodeLabel = worker.Node.CloudNodeLabel

// 			node.Creator = &user
// 			node.CreatedAt = &now
// 			node.CloudUid = cloud.CloudUID
// 			node.CloudClusterUid = cluster.CloudClusterUid

// 			err = txdb.CreateCloudNode(&node)
// 			if err != nil {
// 				txErr := txdb.Rollback()
// 				if txErr != nil {
// 					logger.Info("DB Rollback Failed.", txErr)
// 				}
// 				return response.ErrorfReqRes(c, req, common.CodeFailedDatabase, err)
// 			}
// 			setRemoveNodes.Add(node.CloudNodeUid)
// 		}
// 	}

// 	// Update Nodes
// 	count, err = txdb.UpdateCloudNodes(updateNodes)
// 	if err != nil {
// 		txErr := txdb.Rollback()
// 		if txErr != nil {
// 			logger.Info("DB Rollback Failed.", txErr)
// 		}
// 		return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
// 	}

// 	// Remove Node Updated
// 	difference := resRemoveNodes.Difference(setRemoveNodes)
// 	for _, r := range difference.Values() {
// 		for _, d := range resNodes {
// 			if r == d.CloudNodeUid {
// 				count, err = txdb.DeleteCloudNode(*d.CloudNodeUid)
// 				if err != nil {
// 					txErr := txdb.Rollback()
// 					if txErr != nil {
// 						logger.Info("DB Rollback Failed.", txErr)
// 					}
// 					return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
// 				}
// 			}
// 		}
// 	}

// 	// End. Transaction Commit
// 	txErr := txdb.Commit()
// 	if txErr != nil {
// 		logger.Info("DB commit Failed.", txErr)
// 	}

// 	if count == 0 {
// 		return response.ErrorfReqRes(c, cloudUid, common.DatabaseEmptyData, nil)
// 	}

// 	return response.Write(c, nil, count)
// }

// DeleteCloudHandler - 클라우드 삭제
// @Tags Cloud
// @Summary DeleteCloud
// @Description Delete cloud
// @ID DeleteCloud
// @Produce json
// @Param cloudUId path string true "CloudUid"
// @Success 200 {object} response.ReturnData
// @Router /clouds/{cloudUid} [delete]
// func (a *API) DeleteCloudHandler(c echo.Context) error {
// 	// check param UID
// 	cloudUid, err := uuid.FromString(c.Param("cloudUid"))
// 	if err != nil {
// 		return response.ErrorfReqRes(c, cloudUid, common.CodeInvalidParm, err)
// 	}

// 	// -- Service Logic
// 	// Start. Transaction 얻어옴
// 	txdb, err := a.Db.BeginTransaction()
// 	if err != nil {
// 		return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
// 	}

// 	utils.Print("--DeleteAllCloudNode--")
// 	utils.Print(cloudUid)
// 	// 1. Cloud - Nodes 삭제
// 	count, err := txdb.DeleteAllCloudNode(cloudUid)
// 	if err != nil {
// 		txErr := txdb.Rollback()
// 		if txErr != nil {
// 			logger.Info("DB Rollback Failed.", txErr)
// 		}
// 		return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
// 	}

// 	// 2. Cloud - Cluster 삭제
// 	count, err = txdb.DeleteAllCloudCluster(cloudUid)
// 	if err != nil {
// 		txErr := txdb.Rollback()
// 		if txErr != nil {
// 			logger.Info("DB Rollback Failed.", txErr)
// 		}
// 		return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
// 	}

// 	// 3. Cloud - Cloud 삭제
// 	count, err = txdb.DeleteCloud(cloudUid)
// 	if err != nil {
// 		txErr := txdb.Rollback()
// 		if txErr != nil {
// 			logger.Info("DB Rollback Failed.", txErr)
// 		}
// 		return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
// 	}

// 	// End. Transaction Commit
// 	txErr := txdb.Commit()
// 	if txErr != nil {
// 		logger.Info("DB commit Failed.", txErr)
// 	}

// 	if count == 0 {
// 		return response.ErrorfReqRes(c, cloudUid, common.DatabaseEmptyData, nil)
// 	}

// 	return response.Write(c, nil, count)
// }
