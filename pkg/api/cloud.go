package api

import (
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/response"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	mr "github.com/acornsoft-edgecraft/edgecraft-api/pkg/model/response"
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

	// resCloudCluster, err := a.Db.GetCloudCluster(cloudUid, "cloud")
	// if err != nil {
	// 	return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	// } else if resCloudCluster == nil {
	// 	return response.ErrorfReqRes(c, resCloudCluster, common.DatabaseFalseData, err)
	// }

	// resCloudNode, err := a.Db.GetCloudNode(cloudUid)
	// if err != nil {
	// 	return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	// } else if resCloudNode == nil {
	// 	return response.ErrorfReqRes(c, resCloudNode, common.DatabaseFalseData, err)
	// }

	return response.Write(c, nil, nil)
}

func (a *API) UpdateCloudHandler(c echo.Context) error {
	// check param UID
	cloudUid, err := uuid.FromString(c.Param("cloudUid"))
	if err != nil {
		return response.ErrorfReqRes(c, cloudUid, common.CodeInvalidParm, err)
	}

	now := time.Now()

	var req model.Cloud
	err = getRequestData(c.Request(), &req)
	if err != nil {
		return response.ErrorfReqRes(c, req, common.CodeInvalidData, err)
	}
	req.CloudUID = &cloudUid
	req.UpdatedAt = &now

	// -- Service Logic
	// Start. Transaction 얻어옴
	txdb, err := a.Db.BeginTransaction()
	if err != nil {
		return response.ErrorfReqRes(c, req, common.CodeFailedDatabase, err)
	}

	// Clud 등록 업데이트
	count, err := txdb.UpdateCloud(&req)
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
	cluster.CreatedAt = &now
	cluster.CloudUid = cloud.CloudUID
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
	cluster.CloudClusterStorageClass = &req.EtcdStorage.StorageClass

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
		node.CloudNodeHostName = req.Nodes.MasterNode[i].Baremetal.HostName
		node.CloudNodeBmcAddress = req.Nodes.MasterNode[i].Baremetal.BmcAddress
		node.CloudNodeMacAddress = req.Nodes.MasterNode[i].Baremetal.BootMacAddress
		node.CloudNodeBootMode = req.Nodes.MasterNode[i].Baremetal.BootMode
		node.CloudNodeOnlinePower = req.Nodes.MasterNode[i].Baremetal.OonlinePower
		node.CloudNodeMacAddress = req.Nodes.MasterNode[i].Baremetal.BootMacAddress
		node.CloudNodeExternalProvisioning = req.Nodes.MasterNode[i].Baremetal.ExternalProvisioning
		node.CloudNodeName = req.Nodes.MasterNode[i].Node.NodeName
		node.CloudNodeIp = req.Nodes.MasterNode[i].Node.IpAddress
		node.CloudNodeLabel = &req.Nodes.MasterNode[i].Node.Labels

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
		node.CloudNodeHostName = req.Nodes.WorkerNode[i].Baremetal.HostName
		node.CloudNodeBmcAddress = req.Nodes.WorkerNode[i].Baremetal.BmcAddress
		node.CloudNodeMacAddress = req.Nodes.WorkerNode[i].Baremetal.BootMacAddress
		node.CloudNodeBootMode = req.Nodes.WorkerNode[i].Baremetal.BootMode
		node.CloudNodeOnlinePower = req.Nodes.WorkerNode[i].Baremetal.OonlinePower
		node.CloudNodeMacAddress = req.Nodes.WorkerNode[i].Baremetal.BootMacAddress
		node.CloudNodeExternalProvisioning = req.Nodes.WorkerNode[i].Baremetal.ExternalProvisioning
		node.CloudNodeName = req.Nodes.WorkerNode[i].Node.NodeName
		node.CloudNodeIp = req.Nodes.WorkerNode[i].Node.IpAddress
		node.CloudNodeLabel = &req.Nodes.WorkerNode[i].Node.Labels

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

	// 1. cloud 삭제
	count, err := txdb.DeleteCloud(cloudUid)
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
