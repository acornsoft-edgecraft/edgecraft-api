package api

import (
	"fmt"
	"reflect"
	"strings"
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
	} else if resCloud == nil {
		return response.ErrorfReqRes(c, resCloud, common.DatabaseFalseData, err)
	}

	// resCloudNode, err := a.Db.GetCloudNode(cloudUid, *resCloudCluster.CloudClusterUid)
	// if err != nil {
	// 	return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	// } else if resCloud == nil {
	// 	return response.ErrorfReqRes(c, resCloud, common.DatabaseFalseData, err)
	// }

	resClusterK8s, err := a.Db.SelectK8sCloudCluster(cloudUid)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if resCloud == nil {
		return response.ErrorfReqRes(c, resCloud, common.DatabaseFalseData, err)
	}

	resClusterBaremetal, err := a.Db.SelectBaremetalCloudCluster(cloudUid)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if resCloud == nil {
		return response.ErrorfReqRes(c, resCloud, common.DatabaseFalseData, err)
	}

	resCloudNode, err := a.Db.SelectNodeCloudCluster(cloudUid)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if resCloud == nil {
		return response.ErrorfReqRes(c, resCloud, common.DatabaseFalseData, err)
	}

	resNodes, err := a.Db.SelectCloudNode(cloudUid, *resClusterK8s.CloudClusterUid)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if resCloud == nil {
		return response.ErrorfReqRes(c, resCloud, common.DatabaseFalseData, err)
	}

	var masterNodes []model.MasterNode
	var workerNodes []model.WorkerNode

	for _, nodes := range resNodes {
		if *nodes.CloudNodeType == "master" {
			masterBarematal := model.NodeBaremetal{
				CloudNodeHostName:             nodes.CloudNodeHostName,
				CloudNodeBmcAddress:           nodes.CloudNodeBmcAddress,
				CloudNodeMacAddress:           nodes.CloudNodeMacAddress,
				CloudNodeBootMode:             nodes.CloudNodeBootMode,
				CloudNodeOnlinePower:          nodes.CloudNodeOnlinePower,
				CloudNodeExternalProvisioning: nodes.CloudNodeExternalProvisioning,
			}
			nodes := model.Nodes{
				CloudNodeName:  nodes.CloudNodeName,
				CloudNodeIp:    nodes.CloudNodeIp,
				CloudNodeLabel: nodes.CloudNodeLabel,
			}

			masterNode := model.MasterNode{
				Baremetal: masterBarematal,
				Node:      nodes,
			}
			masterNodes = append(masterNodes, masterNode)
		} else if *nodes.CloudNodeType == "worker" {
			workerBarematal := model.NodeBaremetal{
				CloudNodeHostName:             nodes.CloudNodeHostName,
				CloudNodeBmcAddress:           nodes.CloudNodeBmcAddress,
				CloudNodeMacAddress:           nodes.CloudNodeMacAddress,
				CloudNodeBootMode:             nodes.CloudNodeBootMode,
				CloudNodeOnlinePower:          nodes.CloudNodeOnlinePower,
				CloudNodeExternalProvisioning: nodes.CloudNodeExternalProvisioning,
			}
			nodes := model.Nodes{
				CloudNodeName:  nodes.CloudNodeName,
				CloudNodeIp:    nodes.CloudNodeIp,
				CloudNodeLabel: nodes.CloudNodeLabel,
			}

			workerNode := model.WorkerNode{
				Baremetal: workerBarematal,
				Node:      nodes,
			}
			workerNodes = append(workerNodes, workerNode)
		}

	}

	// var etcd model.Etcd
	// var storage model.StorageClass
	ft := reflect.TypeOf(resCloudCluster)

	var res mr.RegisterCloud
	fType := reflect.TypeOf(res)
	fmt.Println("fType: ", fType.Kind())

	examiner(ft, fType)

	for i := 0; i < fType.NumField(); i++ {
		f := fType.Field(i)
		if f.Name == "adasf" {
			fmt.Println("end")
		}
	}

	fmt.Println("------------------")
	// b, _ := json.MarshalIndent(res, "", " ")
	// fmt.Printf("Before : %+v\n", string(b))
	// SetField(&res, "CloudName", "asdfasdf")
	// fmt.Printf("After : %+v\n", string(b))

	fmt.Println("----- TypeOfStruct -------------")
	type Wham struct {
		Username string   `json:"username,omitempty"`
		Password string   `json:"password"`
		ID       int64    `json:"_id"`
		Homebase []string `json:"homebase"`
	}
	// w := Wham{
	// 	Username: "maria",
	// 	Password: "hunter2",
	// 	ID:       42,
	// 	Homebase: "2434 Main St",
	// }
	// var aa Wham
	// fmt.Printf("%+v\n", aa)
	// SetField(&aa, "username", "larry")
	// SetField(&aa, "_id", 44)
	// SetField(&aa, "homebase", 44)
	// fmt.Printf("%+v\n", aa)

	// text := "asdfasdf"
	// aaa := model.Cloud{
	// 	CloudName: &text,
	// }
	// var aaa model.Cloud
	// bbb, _ := json.MarshalIndent(aaa, "", " ")
	// fmt.Println("Before: ", string(bbb))
	// SetField(&aaa, "name", "1111111")
	// fmt.Println("After: ", string(bbb))

	res.Cloud = *resCloud
	res.Cluster.K8s = *resClusterK8s
	res.Cluster.Baremetal = *resClusterBaremetal
	res.Nodes.CloudClusterLoadbalancerUse = resCloudNode.CloudClusterLoadbalancerUse
	res.Nodes.CloudClusterLoadbalancerAddress = resCloudNode.CloudClusterLoadbalancerAddress
	res.Nodes.CloudClusterLoadbalancerPort = resCloudNode.CloudClusterLoadbalancerPort
	res.Nodes.MasterNode = masterNodes
	res.Nodes.WorkerNode = workerNodes
	res.EtcdStorage.Etcd.CloudClusterExternalEtcdUse = resCloudCluster.CloudClusterExternalEtcdUse
	res.EtcdStorage.Etcd.ExternalEtcdEndPoints = resCloudCluster.ExternalEtcdEndPoints
	res.EtcdStorage.Etcd.ExternalEtcdCertificateCa = resCloudCluster.ExternalEtcdCertificateCa
	res.EtcdStorage.Etcd.ExternalEtcdCertificateCert = resCloudCluster.ExternalEtcdCertificateCert
	res.EtcdStorage.Etcd.ExternalEtcdCertificateKey = resCloudCluster.ExternalEtcdCertificateKey
	res.EtcdStorage.CloudClusterStorageClass = resCloudCluster.CloudClusterStorageClass

	return response.Write(c, nil, &res)
}

func examiner(t reflect.Type, ftype reflect.Type) {
	switch t.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		examiner(t.Elem(), ftype)
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			fmt.Println("asdfasdf:: ", f.Name)
		}
	}
}

func SetField(item interface{}, fieldName string, value interface{}) error {
	v := reflect.ValueOf(item).Elem()
	if !v.CanAddr() {
		return fmt.Errorf("cannot assign to the item passed, item must be a pointer in order to assign")
	}
	// It's possible we can cache this, which is why precompute all these ahead of time.
	findJsonName := func(t reflect.StructTag) (string, error) {
		if jt, ok := t.Lookup("json"); ok {
			return strings.Split(jt, ",")[0], nil
		}
		return "", fmt.Errorf("field %s tag provided does not define a json tag", fieldName)
	}
	fieldNames := map[string]int{}
	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		tag := typeField.Tag
		jname, _ := findJsonName(tag)
		fieldNames[jname] = i
	}

	fieldNum, ok := fieldNames[fieldName]
	if !ok {
		return fmt.Errorf("field %s does not exist within the provided item", fieldName)
	}
	fieldVal := v.Field(fieldNum)
	// fieldVal.SetString(value)
	fieldVal.Set(reflect.ValueOf(value))
	return nil
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
	cluster.CloudClusterStorageClass = req.EtcdStorage.CloudClusterStorageClass

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
		node.CloudNodeMacAddress = req.Nodes.MasterNode[i].Baremetal.CloudNodeMacAddress
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
		node.CloudNodeMacAddress = req.Nodes.WorkerNode[i].Baremetal.CloudNodeMacAddress
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
