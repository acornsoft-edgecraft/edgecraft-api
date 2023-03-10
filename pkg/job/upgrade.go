/*
Copyright 2023 Acornsoft Authors. All right reserved.
*/
package job

import (
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/kubemethod"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/db"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
)

// // checkProvisionCluster - Provision 처리 확인 및 후처리
// func checkProvisionCluster(task string, taskData interface{}) {
// 	data := taskData.(*TaskData)
// 	retryTicker := time.NewTicker(30 * time.Second)
// 	count := 0
// 	for range retryTicker.C {
// 		// Get phase for workload cluster
// 		phase, err := kubemethod.GetProvisionPhase(data.Namespace, data.ClusterName)
// 		if err != nil {
// 			logger.WithField("task", task).WithError(err).Infof("Retrieve provision status for (%s) failed.", data.ClusterName)
// 		} else {
// 			var state int = common.StatusSaved
// 			var provisionState = strings.ToLower(phase)
// 			if provisionState == "provisioned" {
// 				state = common.StatusProvisioned
// 			} else if provisionState == "failed" || provisionState == "pending" || provisionState == "Unknown" {
// 				state = common.StatusFailed
// 			} else if provisionState == "provisioning" {
// 				state = common.StatusProvisioning
// 			}

// 			logger.WithField("task", task).Infof("Checked state (%d), phase (%s), cluster (%s)", state, phase, data.ClusterName)

// 			if provisionState != "" && provisionState != "provisioning" {
// 				// update database. provisioned
// 				affected, err := data.Database.UpdateOpenstackClusterStatus(data.CloudId, data.ClusterId, state)
// 				if err != nil {
// 					logger.WithField("task", task).WithError(err).Infof("Update provision state (%d) for (%s) failed.", state, data.ClusterName)
// 					retryTicker.Stop()
// 					return
// 				} else if affected != 1 {
// 					logger.WithField("task", task).WithError(err).Infof("Close checking for provision state (%d) for (%s) failed. (data not found, check cloud/cluster id)", state, data.ClusterName)
// 					retryTicker.Stop()
// 					return
// 				}

// 				logger.WithField("task", task).Infof("Checking provisioned and update to database [state: %d, cluster: %s]", state, data.ClusterName)
// 				retryTicker.Stop()
// 				return
// 			}
// 		}

// 		count += 1
// 		if count > 120 {
// 			retryTicker.Stop()
// 			logger.WithField("task", task).Info("Close checking provisioned. retry count (120 times in hour) over.")
// 			return
// 		}
// 	}
// }

// checkControlPlaneUpgraded - 클러스터의 컨트롤 플레인 업그레이드 완료 여부 검증
func checkControlPlaneUpgraded(task string, taskData interface{}) {
	data := taskData.(*TaskData)
	var err error

	retryTicker := time.NewTicker(30 * time.Second)
	count := 0

	vals := utils.SplitString("~||~", data.CustomData.(string))

	// Check control plane upgraged
	for range retryTicker.C {
		// check control plane status
		ing, err := kubemethod.GetControlPlaneUpdatePhase(data.Namespace, vals[0])
		if err != nil {
			logger.WithField("task", task).WithError(err).Infof("Control Plane version update for (%s) failed.", data.ClusterName)
			retryTicker.Stop()
			return
		}

		if !ing {
			logger.WithField("task", task).Infof("Control Plane version update for (%s) completed.", data.ClusterName)
			retryTicker.Stop()
			return
		}

		count += 1
		if count > 120 {
			retryTicker.Stop()
			logger.WithField("task", task).Info("Close checking control plane version upgraded. retry count (120 times in hour) over.")
			return
		}
	}

	// Upgrade worker version
	if err == nil {
		err = kubemethod.Apply(data.ClusterName, vals[1])
		if err != nil {
			logger.Errorf("Upgrading Control Plane Kubernetes version failed. (cause: %s)", err.Error())
		}
	}
}

// InvokeK8sVersionUpgrade - 클러스터의 Kubernetes Version Upgrade 처리
func InvokeK8sVersionUpgrade(worker *IWorker, db db.DB, clusterName, namespace, masterSetName, controlPlanesManifest, workersManifest string) error {
	// 컨트롤 플레인 버전 업그레이드 적용 (Kubernetes로 전송)
	err := kubemethod.Apply(clusterName, controlPlanesManifest)
	if err != nil {
		logger.Errorf("Upgrading Control Plane Kubernetes version failed. (cause: %s)", err.Error())
		return err
	}

	// Setting task info
	taskData := &TaskData{
		Database:          db,
		CloudId:           "",
		ClusterId:         "",
		ClusterName:       clusterName,
		BootstrapProvider: common.Kubeadm,
		Namespace:         namespace,
		CustomData:        utils.JoinStrings("~||~", masterSetName, workersManifest),
	}

	taskInfo := TaskInfo{
		TaskData: taskData,
		TaskFunc: checkControlPlaneUpgraded,
	}

	// Running check and worker upgrade (Background)
	err = (*worker).QueueTask("check-control-plane-upgrade", zeroDuration, taskInfo)
	if err != nil {
		return err
	}

	// 	taskInfo.TaskFunc = checkProvisionCluster

	// 	// check provisioned
	// 	err = (*worker).QueueTask("check-provisioned", zeroDuration, taskInfo)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	return nil

	// // 시나리오
	// //	1. Background Job 실행
	// //	2. Control Plane Version Upgrade
	// //	3. Check Control Plane Upgraded
	// //	4. Worker Version Upgrade

	// // TODO: Invoke Background
	// // 1. Control Plane
	// // 2. Worker (Control Plane 완성 후)

	// // // Provision 검증 작업 등록 (Backgroud)
	// // err = job.InvokeProvisionCheck(worker, db, *cluster.CloudUid, *cluster.ClusterUid, *cluster.Name, *cluster.Namespace, *cluster.BootstrapProvider)
	// // if err != nil {
	// // 	logger.WithError(err).Infof("Openstack Cluster [%s] provision check job failed.", cluster.Name)
	// // 	return err
	// // }

	// // logger.Infof("Openstack Cluster [%s] provision submitted.", cluster.Name)

	// // Processing template
	// temp := getFunctionalTemplate(getTemplatePath(cluster.BootstrapProvider, "upgrade"))

	// // // Processing template
	// // fm := template.FuncMap{"replace": replace}
	// // tPath := getTemplatePath(cluster.BootstrapProvider, "upgrade")
	// // temp := template.Must(template.New(path.Base(tPath)).Funcs(fm).ParseFiles(tPath))
	// // temp, err := template.ParseFiles(getTemplatePath(cluster.BootstrapProvider, "upgrade"))
	// // temp = temp.Funcs(template.FuncMap{"replace": replace})

	// // if err != nil {
	// // 	logger.Errorf("Template has errors. cause(%s)", err.Error())
	// // 	return err
	// // }

	// // TODO: 진행상황을 어떻게 클라이언트에 보여줄 것인가?
	// var buff bytes.Buffer
	// err := temp.Execute(&buff, data)
	// if err != nil {
	// 	logger.Errorf("Template execution failed. cause(%s)", err.Error())
	// 	return err
	// }

	// logger.Infof("processed cluster update templating yaml (%s)", buff.String())

	// // 템플릿 적용 (Kubernetes로 전송)
	// err = kubemethod.Apply(*cluster.Name, buff.String())
	// if err != nil {
	// 	logger.Errorf("Kubernetes version upgrade failed. (cause: %s)", err.Error())
	// 	return err
	// }

	// // // 데이터 갱신 (트랜잭션 구간)
	// // err = database.TransactionScope(func(txDB db.DB) error {
	// // 	// Version 정보 갱신
	// // 	cluster.Version = &upgradeInfo.Version
	// // 	cluster.OpenstackInfo.ImageName = upgradeInfo.Image

	// // 	affectedRows, err := txDB.UpdateOpenstackCluster(cluster)
	// // 	if err != nil {
	// // 		return err
	// // 	}
	// // 	if affectedRows == 0 {
	// // 		return errors.New("no data found (update)")
	// // 	}

	// // 	return nil
	// // })

	// if err != nil {
	// 	return err
	// }

	return nil
}

// // InvokeProvisionCheck - Providion 처리 중인 클러스터에 대한 진행 검증 작업
// func InvokeProvisionCheck(worker *IWorker, db db.DB, cloudId, clusterId, clusterName, namespace string, bootstrapProvider common.BootstrapProvider) error {
// 	taskData := &TaskData{
// 		Database:          db,
// 		CloudId:           cloudId,
// 		ClusterId:         clusterId,
// 		ClusterName:       clusterName,
// 		BootstrapProvider: bootstrapProvider,
// 		Namespace:         namespace,
// 	}

// 	taskInfo := TaskInfo{
// 		TaskData: taskData,
// 		TaskFunc: checkProvisionKubeConfig,
// 	}

// 	// check kubeconfig
// 	err := (*worker).QueueTask("check-kubeconfig", zeroDuration, taskInfo)
// 	if err != nil {
// 		return err
// 	}

// 	taskInfo.TaskFunc = checkProvisionCluster

// 	// check provisioned
// 	err = (*worker).QueueTask("check-provisioned", zeroDuration, taskInfo)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
