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
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
)

// checkBackResProcessed - 클러스터의 Backup / Restore 처리 완료 여부 검증
func checkBackResProcessed(task string, taskData interface{}) {
	data := taskData.(*TaskData)
	var isBackup bool = false

	retryTicker := time.NewTicker(30 * time.Second)
	count := 0

	backresInfo := data.CustomData.(model.BackResInfo)
	if backresInfo.Type == "B" {
		isBackup = true
	}

	// Check backup / restore completed.
	for range retryTicker.C {
		// check backup / restore status
		phase, err := kubemethod.GetBackResStatusPhase(data.ClusterName, data.Namespace, backresInfo.Name, isBackup)
		if err != nil {
			logger.WithField("task", task).WithError(err).Infof("Backup/Restore check for (%s) failed.", backresInfo.Name)
			retryTicker.Stop()
			return
		}

		if phase != "" {
			logger.WithField("task", task).Infof("Close checking Backup/Restore for (%s) checked. (%s) state", backresInfo.Name, phase)
			if phase == "C" || phase == "F" || phase == "P" {
				// update database. backup/restore
				affected, err := data.Database.UpdateBackResStatus(backresInfo.CloudUid, backresInfo.ClusterUid, backresInfo.BackResUid, phase)
				if err != nil {
					logger.WithField("task", task).WithError(err).Infof("Update backup/restore state (%s) for (%s) failed.", phase, backresInfo.Name)
				} else if affected != 1 {
					logger.WithField("task", task).WithError(err).Infof("Close checking for backup/restore state (%s) for (%s) failed. (data not found, check cloud/cluster/backres ids)", phase, backresInfo.Name)
				}

				retryTicker.Stop()
				return
			}
		} else {
			logger.WithField("task", task).Infof("Backup/Restore checking for (%s) processing... wait", backresInfo.Name)
		}

		count += 1
		if count > 120 {
			retryTicker.Stop()
			logger.WithField("task", task).Info("Close checking control plane version upgraded. retry count (120 times in hour) over.")
			return
		}
	}
}

// InvokeBackRes - 특정 클러스터에 대한 Backup / Restore 처리
func InvokeBackRes(worker *IWorker, db db.DB, clusterName, namespace string, backresInfo *model.BackResInfo, backresManifest string) error {
	// Backup / Restore Manifest 적용 (Kubernetes로 전송)
	err := kubemethod.Apply(clusterName, backresManifest)
	if err != nil {
		logger.Errorf("Apply backup/restore manifest failed. (cause: %s)", err.Error())
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
		CustomData:        *backresInfo,
	}

	taskInfo := TaskInfo{
		TaskData: taskData,
		TaskFunc: checkBackResProcessed,
	}

	// Running check and worker for Backup/Restore (Background)
	err = (*worker).QueueTask("check-backup-restore-processing", zeroDuration, taskInfo)
	if err != nil {
		return err
	}

	return nil
}
