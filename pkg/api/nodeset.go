/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package api

/*******************************
 ** NodeSet for Openstack
 *******************************/

// TODO: Get NodeSet List
// TODO: Get NodeSet
// TODO: Add NodeSet
// TODO: Update NodeSet
// TODO: Delete NodeSet

// // GetClusterListHandler - 전체 클러스터 리스트 (Openstack)
// // @Tags        Openstack-Cluster
// // @Summary     GetClusterList
// // @Description 전체 클러스터 리스트 (Openstack)
// // @ID          GetClusterList
// // @Produce     json
// // @Param       cloudId path     string true "Cloud ID"
// // @Success     200     {object} response.ReturnData
// // @Router      /clouds/{cloudId}/clusters [get]
// func (a *API) GetClusterListHandler(c echo.Context) error {
// 	cloudId := c.Param("cloudId")
// 	if cloudId == "" {
// 		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
// 	}

// 	result, err := a.Db.GetOpenstackClusters(cloudId)
// 	if err != nil {
// 		return response.Errorf(c, common.CodeFailedDatabase, err)
// 	}

// 	return response.Write(c, nil, result)
// }

// // GetClusterHandler - 클러스터 상세 조회 (Openstack)
// // @Tags        Openstack-Cluster
// // @Summary     GetCluster
// // @Description 클러스터 상세 조회 (Openstack)
// // @ID          GetCluster
// // @Produce     json
// // @Param       cloudId   path     string true "Cloud ID"
// // @Param       clusterId path     string true "Cluster ID"
// // @Success     200       {object} response.ReturnData
// // @Router      /clouds/{cloudId}/clusters/{clusterId} [get]
// func (a *API) GetClusterHandler(c echo.Context) error {
// 	cloudId := c.Param("cloudId")
// 	if cloudId == "" {
// 		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
// 	}

// 	clusterId := c.Param("clusterId")
// 	if cloudId == "" {
// 		return response.ErrorfReqRes(c, clusterId, common.CodeInvalidParm, nil)
// 	}

// 	openstackClusterSet := &model.OpenstackClusterSet{}

// 	// Cluster정보 조회
// 	clusterTable, err := a.Db.GetOpenstackCluster(cloudId, clusterId)
// 	if err != nil {
// 		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
// 	}
// 	if clusterTable == nil {
// 		return response.ErrorfReqRes(c, clusterTable, common.ClusterNotFound, err)
// 	}

// 	// Node 정보 조회
// 	nodeSets, err := a.Db.GetNodeSets(clusterId)
// 	if err != nil {
// 		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
// 	}
// 	if len(nodeSets) == 0 {
// 		return response.ErrorfReqRes(c, clusterTable, common.NodeSetNotFound, err)
// 	}

// 	openstackClusterSet.FromTable(clusterTable, nodeSets)

// 	// Provisioned 상태면 Kubernetes Node 정보 조회 및 설정
// 	if *clusterTable.Status == common.StatusProvisioned {
// 		nodeList, err := kubemethod.GetNodeList(*clusterTable.Name)
// 		if err != nil {
// 			return response.ErrorfReqRes(c, nil, common.CodeFailedK8SAPI, err)
// 		}

// 		// Add node info
// 		for _, node := range nodeList {
// 			find := false
// 			for _, nodeSet := range openstackClusterSet.Nodes.MasterSets {
// 				if strings.Contains(node.Name, "-"+nodeSet.Name+"-") {
// 					nodeSet.Nodes = append(nodeSet.Nodes, node)
// 					find = true
// 					break
// 				}
// 			}

// 			if !find {
// 				for _, nodeSet := range openstackClusterSet.Nodes.WorkerSets {
// 					if strings.Contains(node.Name, "-"+nodeSet.Name+"-") {
// 						nodeSet.Nodes = append(nodeSet.Nodes, node)
// 						break
// 					}
// 				}
// 			}
// 		}
// 	}

// 	return response.Write(c, nil, openstackClusterSet)
// }

// // SetClusterHandler - 클러스터 추가 (Openstack)
// // @Tags        Openstack-Cluster
// // @Summary     SetCluster
// // @Description 클러스터 추가 (Openstack)
// // @ID          SetCluster
// // @Produce     json
// // @Param       cloudId             path     string                    true "Cloud ID"
// // @Param       OpenstackClusterSet body     model.OpenstackClusterSet true "Openstack Cluster Info"
// // @Success     200                 {object} response.ReturnData
// // @Router      /clouds/{cloudId}/clusters [post]
// func (a *API) SetClusterHandler(c echo.Context) error {
// 	// TODO: 로그인 사용자 정보 활용 방법은?
// 	cloudId := c.Param("cloudId")
// 	if cloudId == "" {
// 		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
// 	}

// 	var clusterSet model.OpenstackClusterSet
// 	err := getRequestData(c.Request(), &clusterSet)
// 	if err != nil {
// 		return response.ErrorfReqRes(c, clusterSet, common.CodeInvalidData, err)
// 	}

// 	clusterTable, nodeSetTables := clusterSet.ToTable(cloudId, false, "system", time.Now())

// 	// Start. Transaction 얻어옴
// 	txdb, err := a.Db.BeginTransaction()
// 	if err != nil {
// 		return response.ErrorfReqRes(c, clusterSet, common.CodeFailedDatabase, err)
// 	}

// 	// Cluster 등록
// 	err = txdb.InsertOpenstackCluster(clusterTable)
// 	if err != nil {
// 		txErr := txdb.Rollback()
// 		if txErr != nil {
// 			logger.Info("DB rollback Failed.", txErr)
// 		}
// 		return response.ErrorfReqRes(c, clusterTable, common.CodeFailedDatabase, err)
// 	}

// 	// NodeSet 등록
// 	for _, nodeSetTable := range nodeSetTables {
// 		err = txdb.InsertNodeSet(nodeSetTable)
// 		if err != nil {
// 			txErr := txdb.Rollback()
// 			if txErr != nil {
// 				logger.Info("DB rollback Failed.", txErr)
// 			}
// 			return response.ErrorfReqRes(c, nodeSetTable, common.CodeFailedDatabase, err)
// 		}
// 	}

// 	txErr := txdb.Commit()
// 	if txErr != nil {
// 		logger.Info("DB commit Failed.", txErr)
// 	}

// 	if !clusterSet.SaveOnly {
// 		// Provisioning (background)
// 		err = ProvisioningOpenstackCluster(a.Worker, a.Db, clusterTable, nodeSetTables, a.getCodeNameByKey("K8sVersions", *clusterTable.Version))
// 		if err != nil {
// 			return response.ErrorfReqRes(c, nil, common.ProvisioningCheckJobFailed, err)
// 		}

// 		return response.WriteWithCode(c, clusterSet, common.OpenstackClusterProvisioning, nil)
// 	} else {
// 		// Saved
// 		return response.WriteWithCode(c, clusterSet, common.OpenstackClusterRegistered, nil)
// 	}
// }

// // UpdateClusterHandler - 클러스터 수정 (Openstack)
// // @Tags        Openstack-Cluster
// // @Summary     UpdateCluster
// // @Description 클러스터 수정 (Openstack)
// // @ID          UpdateCluster
// // @Produce     json
// // @Param       cloudId             path     string                    true "Cloud ID"
// // @Param       clusterId           path     string                    true "Cluster ID"
// // @Param       OpenstackClusterSet body     model.OpenstackClusterSet true "Openstack Cluster Info"
// // @Success     200                 {object} response.ReturnData
// // @Router      /clouds/{cloudId}/clusters/{clusterId} [put]
// func (a *API) UpdateClusterHandler(c echo.Context) error {
// 	// TODO: 로그인 사용자 정보 활용 방법은?
// 	cloudId := c.Param("cloudId")
// 	if cloudId == "" {
// 		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
// 	}

// 	clusterId := c.Param("clusterId")
// 	if cloudId == "" {
// 		return response.ErrorfReqRes(c, clusterId, common.CodeInvalidParm, nil)
// 	}

// 	var clusterSet model.OpenstackClusterSet
// 	err := getRequestData(c.Request(), &clusterSet)
// 	if err != nil {
// 		return response.ErrorfReqRes(c, clusterSet, common.CodeInvalidData, err)
// 	}

// 	// 클러스터 정보 조회
// 	clusterTable, err := a.Db.GetOpenstackCluster(cloudId, clusterId)
// 	if err != nil {
// 		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
// 	}
// 	if clusterTable == nil {
// 		return response.ErrorfReqRes(c, nil, common.ClusterNotFound, err)
// 	}

// 	// 클러스터 상태 조회
// 	if *clusterTable.Status != common.StatusSaved && *clusterTable.Status != common.StatusDeleted && *clusterTable.Status != common.StatusFailed {
// 		return response.ErrorfReqRes(c, nil, common.CreatedCloudNoUpdatable, err)
// 	}

// 	// 수신된 변경 정보 구성
// 	clusterTable, nodeSetTables := clusterSet.ToTable(cloudId, false, "system", time.Now())

// 	// 트랜잭션 구간 처리
// 	err = a.Db.TransactionScope(func(txDB db.DB) error {
// 		// Cluster 등록
// 		affectedRows, err := txDB.UpdateOpenstackCluster(clusterTable)
// 		if err != nil {
// 			return err
// 		}
// 		if affectedRows == 0 {
// 			return errors.New("no data found (update)")
// 		}

// 		// 기존 NodeSet 삭제
// 		_, err = txDB.DeleteNodeSets(clusterId)
// 		if err != nil {
// 			return err
// 		}

// 		// NodeSet 등록
// 		for _, nodeSetTable := range nodeSetTables {
// 			err = txDB.InsertNodeSet(nodeSetTable)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		return nil
// 	})
// 	if err != nil {
// 		return response.ErrorfReqRes(c, clusterTable, common.CodeFailedDatabase, err)
// 	}

// 	if !clusterSet.SaveOnly {
// 		// Provisioning (background)
// 		err = ProvisioningOpenstackCluster(a.Worker, a.Db, clusterTable, nodeSetTables, a.getCodeNameByKey("K8sVersions", *clusterTable.Version))
// 		if err != nil {
// 			return response.ErrorfReqRes(c, nil, common.ProvisioningCheckJobFailed, err)
// 		}

// 		return response.WriteWithCode(c, clusterSet, common.OpenstackClusterProvisioning, nil)
// 	} else {
// 		// Saved
// 		return response.WriteWithCode(c, clusterSet, common.OpenstackClusterRegistered, nil)
// 	}
// }

// // DeleteClusterHandler - 클러스터 삭제 (Openstack)
// // @Tags        Openstack-Cluster
// // @Summary     DeleteCluster
// // @Description 클러스터 삭제 (Openstack)
// // @ID          DeleteCluster
// // @Produce     json
// // @Param       cloudId   path     string true "Cloud ID"
// // @Param       clusterId path   string true "Cluster ID"
// // @Success     200     {object} response.ReturnData
// // @Router      /clouds/{cloudId}/clusters/{clusterId} [delete]
// func (a *API) DeleteClusterHandler(c echo.Context) error {
// 	cloudId := c.Param("cloudId")
// 	if cloudId == "" {
// 		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
// 	}

// 	clusterId := c.Param("clusterId")
// 	if cloudId == "" {
// 		return response.ErrorfReqRes(c, clusterId, common.CodeInvalidParm, nil)
// 	}

// 	// 클러스터 정보 조회
// 	clusterTable, err := a.Db.GetOpenstackCluster(cloudId, clusterId)
// 	if err != nil {
// 		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
// 	}
// 	if clusterTable == nil {
// 		return response.ErrorfReqRes(c, nil, common.ClusterNotFound, err)
// 	}

// 	// 삭제 작업
// 	if *clusterTable.Status == common.StatusDeleting {
// 		// 삭제 중이면 종료
// 		return response.WriteWithCode(c, nil, common.OpenstsackClusterAlreadyDeleting, nil)
// 	} else if *clusterTable.Status == common.StatusProvisioning || *clusterTable.Status == common.StatusProvisioned || *clusterTable.Status == common.StatusFailed {
// 		// 상태 provioning(2), provisioned(3), failed(4)인 경우는 클러스터 삭제 진행
// 		// 프로비젼 상태인 클러스터 삭제
// 		err := kubemethod.RemoveOpenstackProvisioned(clusterId, *clusterTable.Name, *clusterTable.Namespace)
// 		if err != nil {
// 			return response.ErrorfReqRes(c, nil, common.DeleteProvisionedClusterJobFailed, err)
// 		}

// 		// Start. Transaction 얻어옴
// 		txdb, err := a.Db.BeginTransaction()
// 		if err != nil {
// 			return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
// 		}

// 		// 클러스터 상태 변경 (deleting)
// 		affectedRows, err := a.Db.UpdateOpenstackClusterStatus(cloudId, clusterId, 5)
// 		if err != nil {
// 			txErr := txdb.Rollback()
// 			if txErr != nil {
// 				logger.Info("DB rollback Failed.", txErr)
// 			}
// 			return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
// 		}
// 		if affectedRows == 0 {
// 			txErr := txdb.Rollback()
// 			if txErr != nil {
// 				logger.Info("DB rollback Failed.", txErr)
// 			}
// 			return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, errors.New("cannot found cluster"))
// 		}

// 		txErr := txdb.Commit()
// 		if txErr != nil {
// 			logger.Info("DB commit Failed.", txErr)
// 		}

// 		// TDelete checking job (remove cluster, remove kubeconfig, update cluster status)
// 		err = job.InvokeDeleteCheck(a.Worker, a.Db, cloudId, clusterId, *clusterTable.Name, *clusterTable.Namespace)
// 		if err != nil {
// 			logger.WithError(err).Infof("Openstack Cluster [%s] provision check job failed.", *clusterTable.Name)
// 			return response.ErrorfReqRes(c, nil, common.DeleteProvisionedClusterJobFailed, err)
// 		}

// 		return response.WriteWithCode(c, nil, common.OpenstackClusterDeleting, nil)
// 	} else {
// 		// Start. Transaction 얻어옴
// 		txdb, err := a.Db.BeginTransaction()
// 		if err != nil {
// 			return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
// 		}

// 		// NodeSet 삭제
// 		affectedRows, err := a.Db.DeleteNodeSets(*clusterTable.ClusterUid)
// 		if err != nil {
// 			txErr := txdb.Rollback()
// 			if txErr != nil {
// 				logger.Info("DB rollback Failed.", txErr)
// 			}
// 			return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
// 		}
// 		if affectedRows == 0 {
// 			txErr := txdb.Rollback()
// 			if txErr != nil {
// 				logger.Info("DB rollback Failed.", txErr)
// 			}
// 			return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, errors.New("cannot found nodesets"))
// 		}

// 		// Cluster Data 삭제
// 		affectedRows, err = a.Db.DeleteOpenstackCluster(cloudId, clusterId)
// 		if err != nil {
// 			txErr := txdb.Rollback()
// 			if txErr != nil {
// 				logger.Info("DB rollback Failed.", txErr)
// 			}
// 			return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
// 		}
// 		if affectedRows == 0 {
// 			txErr := txdb.Rollback()
// 			if txErr != nil {
// 				logger.Info("DB rollback Failed.", txErr)
// 			}
// 			return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, errors.New("cannot found cluster"))
// 		}

// 		txErr := txdb.Commit()
// 		if txErr != nil {
// 			logger.Info("DB commit Failed.", txErr)
// 		}

// 		return response.WriteWithCode(c, nil, common.OpenstackClusterInfoDeleted, nil)
// 	}
// }

// // ProvisioningClusterHandler - 클러스터 Provisioning (Openstack)
// // @Tags        Openstack-Cluster
// // @Summary     ProvisioningCluster
// // @Description 저장된 클러스터 정보를 이용해서 Provision 처리 (Openstack)
// // @ID          ProvisioningCluster
// // @Produce     json
// // @Param       cloudId   path     string true "Cloud ID"
// // @Param       clusterId path     string true "Cluster ID"
// // @Success     200       {object} response.ReturnData
// // @Router      /clouds/{cloudId}/clusters/{clusterId} [post]
// func (a *API) ProvisioningClusterHandler(c echo.Context) error {
// 	cloudId := c.Param("cloudId")
// 	if cloudId == "" {
// 		return response.ErrorfReqRes(c, cloudId, common.CodeInvalidParm, nil)
// 	}

// 	clusterId := c.Param("clusterId")
// 	if cloudId == "" {
// 		return response.ErrorfReqRes(c, clusterId, common.CodeInvalidParm, nil)
// 	}

// 	// Cluster 정보 조회
// 	clusterTable, err := a.Db.GetOpenstackCluster(cloudId, clusterId)
// 	if err != nil {
// 		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
// 	}
// 	if clusterTable == nil {
// 		return response.ErrorfReqRes(c, clusterTable, common.ClusterNotFound, nil)
// 	}

// 	// Cluster 상태 검증
// 	if *clusterTable.Status != common.StatusSaved && *clusterTable.Status != common.StatusDeleted {
// 		return response.ErrorfReqRes(c, clusterTable, common.ProvisioningOnlySavedOrDeleted, err)
// 	}

// 	// NodeSets 정보 조회
// 	nodeSetTables, err := a.Db.GetNodeSets(clusterId)
// 	if err != nil {
// 		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
// 	}
// 	if len(nodeSetTables) == 0 {
// 		return response.ErrorfReqRes(c, nodeSetTables, common.NodeSetNotFound, err)
// 	}

// 	// Provisioning (background)
// 	err = ProvisioningOpenstackCluster(a.Worker, a.Db, clusterTable, nodeSetTables, a.getCodeNameByKey("K8sVersions", *clusterTable.Version))
// 	if err != nil {
// 		return response.ErrorfReqRes(c, nil, common.ProvisioningCheckJobFailed, err)
// 	}

// 	// // Get Pod List
// 	// podList, err := kubemethod.GetPodList("os-cluster", "")
// 	// if err != nil {
// 	// 	return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
// 	// }

// 	// // Get Node List
// 	// nodeList, err := kubemethod.GetNodeList("os-cluster")
// 	// if err != nil {
// 	// 	return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
// 	// }

// 	// // Get Kubeconfig
// 	// data, err := kubemethod.GetKubeconfig(*clusterTable.Namespace, *clusterTable.Name, "value")
// 	// if err != nil {
// 	// 	return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
// 	// }

// 	// // Add cluster's kubeconfig
// 	// err = config.HostCluster.Add([]byte(data))
// 	// if err != nil {
// 	// 	return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
// 	// }

// 	// // Remove cluster's kubeconfig
// 	// err := config.HostCluster.Remove("os-cluster")
// 	// if err != nil {
// 	// 	return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
// 	// }

// 	// Check workload cluster provisioning complete
// 	// data, err := kubemethod.GetProvisioned(*clusterTable.Namespace, *clusterTable.Name,
// 	// 	"infrastructure.cluster.x-k8s.io", "v1alpha3", "openstackclusters")
// 	// if err != nil {
// 	// 	return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
// 	// }

// 	//return response.WriteWithCode(c, nil, common.OpenstackClusterProvisioning, data)
// 	//return response.WriteWithCode(c, nil, common.OpenstackClusterProvisioning, podList)
// 	return response.WriteWithCode(c, nil, common.OpenstackClusterProvisioning, nil)
// }