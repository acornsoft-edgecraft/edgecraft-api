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

		if ing {
			logger.WithField("task", task).Infof("Control Plane version update for (%s) completed.", data.ClusterName)
			retryTicker.Stop()
			break
		} else {
			logger.WithField("task", task).Infof("Control Plane version update for (%s) processing... wait", data.ClusterName)
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
		//err = kubemethod.Apply(data.ClusterName, vals[1])
		err = kubemethod.Apply("", vals[1])
		if err != nil {
			logger.Errorf("Upgrading Control Plane Kubernetes version failed. (cause: %s)", err.Error())
		} else {
			logger.Infof("Update worker machines for (%s) started", data.ClusterName)
		}
	}
}

// InvokeK8sVersionUpgrade - 클러스터의 Kubernetes Version Upgrade 처리
func InvokeK8sVersionUpgrade(worker *IWorker, db db.DB, clusterName, namespace, masterSetName, controlPlanesManifest, workersManifest string) error {
	// 컨트롤 플레인 버전 업그레이드 적용 (Kubernetes로 전송)
	//err := kubemethod.Apply(clusterName, controlPlanesManifest)
	err := kubemethod.Apply("", controlPlanesManifest)
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

	return nil
}
