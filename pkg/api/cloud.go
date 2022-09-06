package api

import (
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/response"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	mr "github.com/acornsoft-edgecraft/edgecraft-api/pkg/model/response"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

func (a *API) AllCloudListHandler(c echo.Context) error {
	res, err := a.Db.GetAllCloud()
	if err != nil {
		return response.Errorf(c, common.CodeFailedDatabase, err)
	}
	return response.Write(c, nil, res)
}

func (a *API) GetCloudHandler(c echo.Context) error {
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

func (a *API) SelectCloudHandler(c echo.Context) error {
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

	resCloudCluster, err := a.Db.GetCloudCluster(cloudUid)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if resCloudCluster == nil {
		return response.ErrorfReqRes(c, resCloudCluster, common.DatabaseFalseData, err)
	}

	// resCloudNode, err := a.Db.GetCloudNode(cloudUid, *resCloudCluster.CloudClusterUid)
	// if err != nil {
	// 	return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	// } else if resCloud == nil {
	// 	return response.ErrorfReqRes(c, resCloud, common.DatabaseFalseData, err)
	// }

	// resClusterK8s, err := a.Db.SelectK8sCloudCluster(cloudUid)
	// if err != nil {
	// 	return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	// } else if resCloud == nil {
	// 	return response.ErrorfReqRes(c, resCloud, common.DatabaseFalseData, err)
	// }

	// resClusterBaremetal, err := a.Db.SelectBaremetalCloudCluster(cloudUid)
	// if err != nil {
	// 	return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	// } else if resCloud == nil {
	// 	return response.ErrorfReqRes(c, resCloud, common.DatabaseFalseData, err)
	// }

	resCloudNode, err := a.Db.SelectNodeCloudCluster(cloudUid)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if resCloudNode == nil {
		return response.ErrorfReqRes(c, resCloudNode, common.DatabaseFalseData, err)
	}

	resNodes, err := a.Db.SelectCloudNode(cloudUid, *resCloudCluster.CloudClusterUid)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if resNodes == nil {
		return response.ErrorfReqRes(c, resNodes, common.DatabaseFalseData, err)
	}

	var masterNodes []model.MasterNode
	var workerNodes []model.WorkerNode

	for _, nodes := range resNodes {
		if *nodes.CloudNodeType == "master" {
			// masterBarematal := model.NodeBaremetal{
			// 	CloudNodeHostName:             nodes.CloudNodeHostName,
			// 	CloudNodeBmcAddress:           nodes.CloudNodeBmcAddress,
			// 	CloudNodeMacAddress:           nodes.CloudNodeMacAddress,
			// 	CloudNodeBootMode:             nodes.CloudNodeBootMode,
			// 	CloudNodeOnlinePower:          nodes.CloudNodeOnlinePower,
			// 	CloudNodeExternalProvisioning: nodes.CloudNodeExternalProvisioning,
			// }

			// nodes := model.Nodes{
			// 	CloudNodeName:  nodes.CloudNodeName,
			// 	CloudNodeIp:    nodes.CloudNodeIp,
			// 	CloudNodeLabel: nodes.CloudNodeLabel,
			// }

			// masterNode := model.MasterNode{
			// 	Baremetal: masterBarematal,
			// 	Node:      nodes,
			// }

			var masterNode model.MasterNode
			if err := utils.SetFieldsInStruct(&masterNode, &nodes); err != nil {
				logger.Errorf("SetFields In Struct ERROR : %s", err)
			}
			masterNodes = append(masterNodes, masterNode)
		} else if *nodes.CloudNodeType == "worker" {
			// workerBarematal := model.NodeBaremetal{
			// 	CloudNodeHostName:             nodes.CloudNodeHostName,
			// 	CloudNodeBmcAddress:           nodes.CloudNodeBmcAddress,
			// 	CloudNodeMacAddress:           nodes.CloudNodeMacAddress,
			// 	CloudNodeBootMode:             nodes.CloudNodeBootMode,
			// 	CloudNodeOnlinePower:          nodes.CloudNodeOnlinePower,
			// 	CloudNodeExternalProvisioning: nodes.CloudNodeExternalProvisioning,
			// }

			// nodes := model.Nodes{
			// 	CloudNodeName:  nodes.CloudNodeName,
			// 	CloudNodeIp:    nodes.CloudNodeIp,
			// 	CloudNodeLabel: nodes.CloudNodeLabel,
			// }

			// workerNode := model.WorkerNode{
			// 	Baremetal: workerBarematal,
			// 	Node:      nodes,
			// }

			var workerNode model.WorkerNode
			if err := utils.SetFieldsInStruct(&workerNode, &nodes); err != nil {
				logger.Errorf("SetFields In Struct ERROR : %s", err)
			}
			workerNodes = append(workerNodes, workerNode)
		}

	}

	var res mr.RegisterCloud

	res.Cloud = *resCloud
	// res.Cluster.K8s = *resClusterK8s
	// res.Cluster.Baremetal = *resClusterBaremetal

	// Set Cluster fields
	if err := utils.SetFieldsInStruct(&res.Cluster, resCloudCluster); err != nil {
		logger.Errorf("SetFields In Struct ERROR : %s", err)
	}

	// Set Nodes fields
	if err := utils.SetFieldsInStruct(&res.Nodes, resCloudNode); err != nil {
		logger.Errorf("SetFields In Struct ERROR : %s", err)
	}
	// res.Nodes.CloudClusterLoadbalancerUse = resCloudNode.CloudClusterLoadbalancerUse
	// res.Nodes.CloudClusterLoadbalancerAddress = resCloudNode.CloudClusterLoadbalancerAddress
	// res.Nodes.CloudClusterLoadbalancerPort = resCloudNode.CloudClusterLoadbalancerPort
	res.Nodes.MasterNode = masterNodes
	res.Nodes.WorkerNode = workerNodes

	// Set ETCD fields
	if err := utils.SetFieldsInStruct(&res.EtcdStorage, resCloudCluster); err != nil {
		logger.Errorf("SetFields In Struct ERROR : %s", err)
	}
	// res.EtcdStorage.Etcd.CloudClusterExternalEtcdUse = resCloudCluster.CloudClusterExternalEtcdUse
	// res.EtcdStorage.Etcd.ExternalEtcdEndPoints = resCloudCluster.ExternalEtcdEndPoints
	// res.EtcdStorage.Etcd.ExternalEtcdCertificateCa = resCloudCluster.ExternalEtcdCertificateCa
	// res.EtcdStorage.Etcd.ExternalEtcdCertificateCert = resCloudCluster.ExternalEtcdCertificateCert
	// res.EtcdStorage.Etcd.ExternalEtcdCertificateKey = resCloudCluster.ExternalEtcdCertificateKey
	// res.EtcdStorage.CloudClusterStorageClass = resCloudCluster.CloudClusterStorageClass

	return response.Write(c, nil, &res)
}

func (a *API) UpdateCloudHandler(c echo.Context) error {
	// check param UID
	cloudUid, err := uuid.FromString(c.Param("cloudUid"))
	if err != nil {
		return response.ErrorfReqRes(c, cloudUid, common.CodeInvalidParm, err)
	}

	now := time.Now()

	var req mr.RegisterCloud
	err = getRequestData(c.Request(), &req)
	if err != nil {
		return response.ErrorfReqRes(c, req, common.CodeInvalidData, err)
	}

	// -- Service Logic
	// Start. Transaction 얻어옴
	txdb, err := a.Db.BeginTransaction()
	if err != nil {
		return response.ErrorfReqRes(c, req, common.CodeFailedDatabase, err)
	}

	// 1. Cloud 등록 업데이트
	var cloud model.Cloud
	getCloud, err := txdb.GetCloud(cloudUid)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB Rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
	} else if getCloud == nil {
		return response.ErrorfReqRes(c, getCloud, common.DatabaseFalseData, err)
	}

	// cloud.CloudName = getCloud.CloudName
	// Set fields in struct
	// if err := utils.SetFieldsInStruct(&cloud, &req.Cloud); err != nil {
	// 	logger.Errorf("SetFields In Struct ERROR : %s", err)
	// }

	user := "user"
	cloud = *getCloud
	cloud = req.Cloud

	cloud.UpdatedAt = &now
	cloud.Updater = &user
	count, err := txdb.UpdateCloud(&cloud)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB Rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
	}

	// 2. Cluster 등록 업데이트
	var cluster model.CloudCluster
	resCloudCluster, err := a.Db.GetCloudCluster(cloudUid)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if resCloudCluster == nil {
		return response.ErrorfReqRes(c, resCloudCluster, common.DatabaseFalseData, err)
	}

	cluster = *resCloudCluster

	cluster.CloudK8sVersion = req.Cluster.K8s.CloudK8sVersion
	cluster.CloudClusterPodCidr = req.Cluster.K8s.CloudClusterPodCidr
	cluster.CloudClusterServiceCidr = req.Cluster.K8s.CloudClusterServiceCidr
	cluster.CloudClusterBmcCredentialSecret = req.Cluster.Baremetal.CloudClusterBmcCredentialSecret
	cluster.CloudClusterBmcCredentialUser = req.Cluster.Baremetal.CloudClusterBmcCredentialUser
	cluster.CloudClusterBmcCredentialPassword = req.Cluster.Baremetal.CloudClusterBmcCredentialPassword
	cluster.CloudClusterImageUrl = req.Cluster.Baremetal.CloudClusterImageUrl
	cluster.CloudClusterImageChecksum = req.Cluster.Baremetal.CloudClusterImageChecksum
	cluster.CloudClusterImageChecksumType = req.Cluster.Baremetal.CloudClusterImageChecksumType
	cluster.CloudClusterImageFormat = req.Cluster.Baremetal.CloudClusterImageFormat
	cluster.CloudClusterMasterExtraConfig = req.Cluster.Baremetal.CloudClusterMasterExtraConfig
	cluster.CloudClusterWorkerExtraConfig = req.Cluster.Baremetal.CloudClusterWorkerExtraConfig
	cluster.CloudClusterLoadbalancerUse = req.Nodes.CloudClusterLoadbalancerUse
	cluster.CloudClusterLoadbalancerAddress = req.Nodes.CloudClusterLoadbalancerAddress
	cluster.CloudClusterLoadbalancerPort = req.Nodes.CloudClusterLoadbalancerPort
	cluster.CloudClusterExternalEtcdUse = req.EtcdStorage.Etcd.CloudClusterExternalEtcdUse
	cluster.ExternalEtcdEndPoints = req.EtcdStorage.Etcd.ExternalEtcdEndPoints
	cluster.ExternalEtcdCertificateCa = req.EtcdStorage.Etcd.ExternalEtcdCertificateCa
	cluster.ExternalEtcdCertificateCert = req.EtcdStorage.Etcd.ExternalEtcdCertificateCert
	cluster.ExternalEtcdCertificateKey = req.EtcdStorage.Etcd.ExternalEtcdCertificateKey
	cluster.CloudClusterStorageClass = req.EtcdStorage.CloudClusterStorageClass

	cluster.Updater = &user
	cluster.UpdatedAt = &now
	cluster.CloudUid = cloud.CloudUID

	count, err = txdb.UpdateCloudCluster(&cluster)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB Rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
	}

	// 3. Node 등록 업데이트
	var node model.CloudNode
	resCloudNode, err := a.Db.SelectNodeCloudCluster(cloudUid)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if resCloudNode == nil {
		return response.ErrorfReqRes(c, resCloudNode, common.DatabaseFalseData, err)
	}

	resNodes, err := a.Db.SelectCloudNode(cloudUid, *resCloudCluster.CloudClusterUid)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if resNodes == nil {
		return response.ErrorfReqRes(c, resNodes, common.DatabaseFalseData, err)
	}

	nodeType := ""

	var resNodeUid []interface{}
	for _, nodes := range resNodes {
		resNodeUid = append(resNodeUid, nodes.CloudNodeUid)
	}

	resRemoveNodes := hashset.New()
	resRemoveNodes.Add(resNodeUid...)

	setRemoveNodes := hashset.New()
	var updateNodes []*model.CloudNode
	// - MasterNode 등록 업데이트
	for _, master := range req.Nodes.MasterNode {
		for _, field := range resNodes {
			if master.CloudNodeUid != nil && *master.CloudNodeUid == *field.CloudNodeUid {
				// nodeType = "master"
				// field.CloudNodeType = &nodeType
				field.CloudNodeHostName = master.Baremetal.CloudNodeHostName
				field.CloudNodeBmcAddress = master.Baremetal.CloudNodeBmcAddress
				field.CloudNodeMacAddress = master.Baremetal.CloudNodeMacAddress
				field.CloudNodeBootMode = master.Baremetal.CloudNodeBootMode
				field.CloudNodeOnlinePower = master.Baremetal.CloudNodeOnlinePower
				field.CloudNodeExternalProvisioning = master.Baremetal.CloudNodeExternalProvisioning
				field.CloudNodeName = master.Node.CloudNodeName
				field.CloudNodeIp = master.Node.CloudNodeIp
				field.CloudNodeLabel = master.Node.CloudNodeLabel

				field.Updater = &user
				field.UpdatedAt = &now

				updateNodes = append(updateNodes, &field)
				// count, err = txdb.UpdateCloudNode(&field)
				// if err != nil {
				// 	txErr := txdb.Rollback()
				// 	if txErr != nil {
				// 		logger.Info("DB Rollback Failed.", txErr)
				// 	}
				// 	return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
				// }
				setRemoveNodes.Add(field.CloudNodeUid)
			}
		}
		// New Node
		if master.CloudNodeUid == nil {
			nodeType = "master"
			node.CloudNodeType = &nodeType
			node.CloudNodeHostName = master.Baremetal.CloudNodeHostName
			node.CloudNodeBmcAddress = master.Baremetal.CloudNodeBmcAddress
			node.CloudNodeMacAddress = master.Baremetal.CloudNodeMacAddress
			node.CloudNodeBootMode = master.Baremetal.CloudNodeBootMode
			node.CloudNodeOnlinePower = master.Baremetal.CloudNodeOnlinePower
			node.CloudNodeExternalProvisioning = master.Baremetal.CloudNodeExternalProvisioning
			node.CloudNodeName = master.Node.CloudNodeName
			node.CloudNodeIp = master.Node.CloudNodeIp
			node.CloudNodeLabel = master.Node.CloudNodeLabel

			node.Creator = &user
			node.CreatedAt = &now
			node.CloudUid = cloud.CloudUID
			node.CloudClusterUid = cluster.CloudClusterUid

			err = txdb.CreateCloudNode(&node)
			if err != nil {
				txErr := txdb.Rollback()
				if txErr != nil {
					logger.Info("DB Rollback Failed.", txErr)
				}
				return response.ErrorfReqRes(c, req, common.CodeFailedDatabase, err)
			}
			setRemoveNodes.Add(node.CloudNodeUid)
		}
	}

	// - WorkerrNode 등록 업데이트
	for _, worker := range req.Nodes.WorkerNode {
		for _, field := range resNodes {
			if worker.CloudNodeUid != nil && *worker.CloudNodeUid == *field.CloudNodeUid {
				// nodeType = "worker"
				// field.CloudNodeType = &nodeType
				field.CloudNodeHostName = worker.Baremetal.CloudNodeHostName
				field.CloudNodeBmcAddress = worker.Baremetal.CloudNodeBmcAddress
				field.CloudNodeMacAddress = worker.Baremetal.CloudNodeMacAddress
				field.CloudNodeBootMode = worker.Baremetal.CloudNodeBootMode
				field.CloudNodeOnlinePower = worker.Baremetal.CloudNodeOnlinePower
				field.CloudNodeExternalProvisioning = worker.Baremetal.CloudNodeExternalProvisioning
				field.CloudNodeName = worker.Node.CloudNodeName
				field.CloudNodeIp = worker.Node.CloudNodeIp
				field.CloudNodeLabel = worker.Node.CloudNodeLabel

				field.Updater = &user
				field.UpdatedAt = &now

				updateNodes = append(updateNodes, &field)
				// count, err = txdb.UpdateCloudNode(&field)
				// if err != nil {
				// 	txErr := txdb.Rollback()
				// 	if txErr != nil {
				// 		logger.Info("DB Rollback Failed.", txErr)
				// 	}
				// 	return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
				// }
				setRemoveNodes.Add(field.CloudNodeUid)
			}
		}
		// New Node
		if worker.CloudNodeUid == nil {
			nodeType = "worker"
			node.CloudNodeType = &nodeType
			node.CloudNodeHostName = worker.Baremetal.CloudNodeHostName
			node.CloudNodeBmcAddress = worker.Baremetal.CloudNodeBmcAddress
			node.CloudNodeMacAddress = worker.Baremetal.CloudNodeMacAddress
			node.CloudNodeBootMode = worker.Baremetal.CloudNodeBootMode
			node.CloudNodeOnlinePower = worker.Baremetal.CloudNodeOnlinePower
			node.CloudNodeExternalProvisioning = worker.Baremetal.CloudNodeExternalProvisioning
			node.CloudNodeName = worker.Node.CloudNodeName
			node.CloudNodeIp = worker.Node.CloudNodeIp
			node.CloudNodeLabel = worker.Node.CloudNodeLabel

			node.Creator = &user
			node.CreatedAt = &now
			node.CloudUid = cloud.CloudUID
			node.CloudClusterUid = cluster.CloudClusterUid

			err = txdb.CreateCloudNode(&node)
			if err != nil {
				txErr := txdb.Rollback()
				if txErr != nil {
					logger.Info("DB Rollback Failed.", txErr)
				}
				return response.ErrorfReqRes(c, req, common.CodeFailedDatabase, err)
			}
			setRemoveNodes.Add(node.CloudNodeUid)
		}
	}

	// Update Nodes
	count, err = txdb.UpdateCloudNodes(updateNodes)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB Rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
	}

	// Remove Node Updated
	difference := resRemoveNodes.Difference(setRemoveNodes)
	for _, r := range difference.Values() {
		for _, d := range resNodes {
			if r == d.CloudNodeUid {
				count, err = txdb.DeleteCloudNode(*d.CloudNodeUid)
				if err != nil {
					txErr := txdb.Rollback()
					if txErr != nil {
						logger.Info("DB Rollback Failed.", txErr)
					}
					return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
				}
			}
		}
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

func (a *API) RegisterCloudHandler(c echo.Context) error {
	now := time.Now()

	var req mr.RegisterCloud
	err := getRequestData(c.Request(), &req)
	if err != nil {
		return response.ErrorfReqRes(c, req, common.CodeInvalidData, err)
	}

	// -- Service Logic
	// Start. Transaction 얻어옴
	txdb, err := a.Db.BeginTransaction()
	if err != nil {
		return response.ErrorfReqRes(c, req, common.CodeFailedDatabase, err)
	}

	// 1. Clud 등록
	var cloud model.Cloud = req.Cloud
	cloud.CreatedAt = &now
	err = txdb.CreateCloud(&cloud)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB Rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, req, common.CodeFailedDatabase, err)
	}
	req.Cloud = cloud

	// 2. Cluster 등록
	var cluster model.CloudCluster
	if err := utils.SetFieldsInStruct(&cluster, &req.Cluster.K8s); err != nil {
		logger.Errorf("SetFields In Struct ERROR : %s", err)
	}
	if err := utils.SetFieldsInStruct(&cluster, &req.Cluster.Baremetal); err != nil {
		logger.Errorf("SetFields In Struct ERROR : %s", err)
	}
	if err := utils.SetFieldsInStruct(&cluster, &req.Nodes); err != nil {
		logger.Errorf("SetFields In Struct ERROR : %s", err)
	}
	if err := utils.SetFieldsInStruct(&cluster, &req.EtcdStorage.Etcd); err != nil {
		logger.Errorf("SetFields In Struct ERROR : %s", err)
	}
	if err := utils.SetFieldsInStruct(&cluster, &req.EtcdStorage); err != nil {
		logger.Errorf("SetFields In Struct ERROR : %s", err)
	}
	// if err := utils.SetFieldsInStruct(&cluster, &req); err != nil {
	// 	logger.Errorf("SetFields In Struct ERROR : %s", err)
	// }
	// cluster.CloudK8sVersion = req.Cluster.K8s.CloudK8sVersion
	// cluster.CloudClusterPodCidr = req.Cluster.K8s.CloudClusterPodCidr
	// cluster.CloudClusterServiceCidr = req.Cluster.K8s.CloudClusterServiceCidr
	// cluster.CloudClusterBmcCredentialSecret = req.Cluster.Baremetal.CloudClusterBmcCredentialSecret
	// cluster.CloudClusterBmcCredentialUser = req.Cluster.Baremetal.CloudClusterBmcCredentialUser
	// cluster.CloudClusterBmcCredentialPassword = req.Cluster.Baremetal.CloudClusterBmcCredentialPassword
	// cluster.CloudClusterImageUrl = req.Cluster.Baremetal.CloudClusterImageUrl
	// cluster.CloudClusterImageChecksum = req.Cluster.Baremetal.CloudClusterImageChecksum
	// cluster.CloudClusterImageChecksumType = req.Cluster.Baremetal.CloudClusterImageChecksumType
	// cluster.CloudClusterImageFormat = req.Cluster.Baremetal.CloudClusterImageFormat
	// cluster.CloudClusterMasterExtraConfig = req.Cluster.Baremetal.CloudClusterMasterExtraConfig
	// cluster.CloudClusterWorkerExtraConfig = req.Cluster.Baremetal.CloudClusterWorkerExtraConfig
	// cluster.CloudClusterLoadbalancerUse = req.Nodes.CloudClusterLoadbalancerUse
	// cluster.CloudClusterLoadbalancerAddress = req.Nodes.CloudClusterLoadbalancerAddress
	// cluster.CloudClusterLoadbalancerPort = req.Nodes.CloudClusterLoadbalancerPort
	// cluster.CloudClusterExternalEtcdUse = req.EtcdStorage.Etcd.CloudClusterExternalEtcdUse
	// cluster.ExternalEtcdEndPoints = req.EtcdStorage.Etcd.ExternalEtcdEndPoints
	// cluster.ExternalEtcdCertificateCa = req.EtcdStorage.Etcd.ExternalEtcdCertificateCa
	// cluster.ExternalEtcdCertificateCert = req.EtcdStorage.Etcd.ExternalEtcdCertificateCert
	// cluster.ExternalEtcdCertificateKey = req.EtcdStorage.Etcd.ExternalEtcdCertificateKey
	// cluster.CloudClusterStorageClass = req.EtcdStorage.CloudClusterStorageClass

	cluster.CreatedAt = &now
	cluster.CloudUid = cloud.CloudUID

	err = txdb.CreateCloudCluster(&cluster)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB Rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, req, common.CodeFailedDatabase, err)
	}

	// 3. Node 등록
	var node model.CloudNode
	node.CreatedAt = &now
	node.CloudUid = cloud.CloudUID
	node.CloudClusterUid = cluster.CloudClusterUid
	nodeType := ""

	// - MasterNode 등록
	for i := 0; i < len(req.Nodes.MasterNode); i++ {
		nodeType = "master"
		node.CloudNodeType = &nodeType
		node.CloudNodeHostName = req.Nodes.MasterNode[i].Baremetal.CloudNodeHostName
		node.CloudNodeBmcAddress = req.Nodes.MasterNode[i].Baremetal.CloudNodeBmcAddress
		node.CloudNodeMacAddress = req.Nodes.MasterNode[i].Baremetal.CloudNodeMacAddress
		node.CloudNodeBootMode = req.Nodes.MasterNode[i].Baremetal.CloudNodeBootMode
		node.CloudNodeOnlinePower = req.Nodes.MasterNode[i].Baremetal.CloudNodeOnlinePower
		node.CloudNodeExternalProvisioning = req.Nodes.MasterNode[i].Baremetal.CloudNodeExternalProvisioning
		node.CloudNodeName = req.Nodes.MasterNode[i].Node.CloudNodeName
		node.CloudNodeIp = req.Nodes.MasterNode[i].Node.CloudNodeIp
		node.CloudNodeLabel = req.Nodes.MasterNode[i].Node.CloudNodeLabel

		err = txdb.CreateCloudNode(&node)
		if err != nil {
			txErr := txdb.Rollback()
			if txErr != nil {
				logger.Info("DB Rollback Failed.", txErr)
			}
			return response.ErrorfReqRes(c, req, common.CodeFailedDatabase, err)
		}

	}

	// - WorkerrNode 등록
	for i := 0; i < len(req.Nodes.WorkerNode); i++ {
		nodeType = "worker"
		node.CloudNodeType = &nodeType
		node.CloudNodeHostName = req.Nodes.WorkerNode[i].Baremetal.CloudNodeHostName
		node.CloudNodeBmcAddress = req.Nodes.WorkerNode[i].Baremetal.CloudNodeBmcAddress
		node.CloudNodeMacAddress = req.Nodes.WorkerNode[i].Baremetal.CloudNodeMacAddress
		node.CloudNodeBootMode = req.Nodes.WorkerNode[i].Baremetal.CloudNodeBootMode
		node.CloudNodeOnlinePower = req.Nodes.WorkerNode[i].Baremetal.CloudNodeOnlinePower
		node.CloudNodeExternalProvisioning = req.Nodes.WorkerNode[i].Baremetal.CloudNodeExternalProvisioning
		node.CloudNodeName = req.Nodes.WorkerNode[i].Node.CloudNodeName
		node.CloudNodeIp = req.Nodes.WorkerNode[i].Node.CloudNodeIp
		node.CloudNodeLabel = req.Nodes.WorkerNode[i].Node.CloudNodeLabel

		err = txdb.CreateCloudNode(&node)
		if err != nil {
			txErr := txdb.Rollback()
			if txErr != nil {
				logger.Info("DB Rollback Failed.", txErr)
			}
			return response.ErrorfReqRes(c, req, common.CodeFailedDatabase, err)
		}

	}

	// End. Transaction Commit
	txErr := txdb.Commit()
	if txErr != nil {
		logger.Info("DB commit Failed.", txErr)
	}

	return response.Write(c, req, req)
}

func (a *API) DeleteCloudHandler(c echo.Context) error {
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

	utils.Print("--DeleteAllCloudNode--")
	utils.Print(cloudUid)
	// 1. Cloud - Nodes 삭제
	count, err := txdb.DeleteAllCloudNode(cloudUid)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB Rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
	}

	// 2. Cloud - Cluster 삭제
	count, err = txdb.DeleteAllCloudCluster(cloudUid)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB Rollback Failed.", txErr)
		}
		return response.ErrorfReqRes(c, cloudUid, common.CodeFailedDatabase, err)
	}

	// 3. Cloud - Cloud 삭제
	count, err = txdb.DeleteCloud(cloudUid)
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
