/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package job

import (
	"strings"
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/kubemethod"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/config"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/db"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
)

const (
	zeroDuration time.Duration = 0
)

// TaskData - Provision에 사용할 클러스터 식별 데이터
type TaskData struct {
	Database    db.DB
	CloudId     string
	ClusterId   string
	ClusterName string
	Namespace   string
}

// checkKubeConfig - Provision 처리 중인 클러스터의 kubeconfig 정보 처리
func checkKubeConfig(task string, taskData interface{}) {
	data := taskData.(*TaskData)
	retryTicker := time.NewTicker(10 * time.Second)
	count := 0
	for range retryTicker.C {
		// Get kubeconfig for workload cluster
		kubeconfig, err := kubemethod.GetKubeconfig(data.Namespace, data.ClusterName, "value")
		if err != nil {
			logger.WithField("task", task).WithError(err).Infof("Retrieve kubeconfig for (%s) failed.", data.ClusterName)
		} else {
			// Add cluster's kubeconfig
			err = config.HostCluster.Add([]byte(kubeconfig))
			if err != nil {
				logger.WithField("task", task).WithError(err).Infof("Add kubeconfig for (%s) to configmap failed.", data.ClusterName)
			} else {
				logger.WithField("task", task).Info("Checking kubeconfig and applied")
				retryTicker.Stop()
				return
			}
		}

		count += 1
		if count > 30 {
			retryTicker.Stop()
			logger.WithField("task", task).Info("Close checking for kubeconfig. retry count (30 times in 5 minues) over.")
			return
		}
	}
}

// checkProvisioned - Provision 처리 종료 여부
func checkProvisioned(task string, taskData interface{}) {
	data := taskData.(*TaskData)
	retryTicker := time.NewTicker(30 * time.Second)
	count := 0
	for range retryTicker.C {
		// Get phase for workload cluster
		phase, err := kubemethod.GetProvisionPhase(data.Namespace, data.ClusterName)
		if err != nil {
			logger.WithField("task", task).WithError(err).Infof("Retrieve provision status for (%s) failed.", data.ClusterName)
		} else {
			var state int = 1
			if strings.ToLower(phase) == "Provisioned" {
				state = 3
			} else if phase == "Failed" || phase == "Pending" || phase == "Unknown" {
				state = 4
			} else if phase == "provisioning" {
				state = 2
			}

			if phase != "" && phase != "Provisioning" {
				// update database. provisioned
				affected, err := data.Database.UpdateOpenstackClusterStatus(data.CloudId, data.ClusterId, state)
				if err != nil {
					logger.WithField("task", task).WithError(err).Infof("Update provision state (%d) for (%s) failed.", state, data.ClusterName)
					retryTicker.Stop()
					return
				} else if affected != 1 {
					logger.WithField("task", task).WithError(err).Infof("Close checking for provision state (%d) for (%s) failed. (data not found, check cloud/cluster id)", state, data.ClusterName)
					retryTicker.Stop()
					return
				}

				logger.WithField("task", task).Info("Checking provisioned and update to database")
				retryTicker.Stop()
				return
			}
		}

		count += 1
		if count > 120 {
			retryTicker.Stop()
			logger.WithField("task", task).Info("Close checking provisioned. retry count (120 times in hour) over.")
			return
		}
	}
}

// InvokeProvisionCheck - Providion 처리 중인 클러스터에 대한 진행 검증 작업
func InvokeProvisionCheck(worker *IWorker, db db.DB, cloudId, clusterId, clusterName, namespace string) {
	taskData := &TaskData{
		Database:    db,
		CloudId:     cloudId,
		ClusterId:   clusterId,
		ClusterName: clusterName,
		Namespace:   namespace,
	}

	taskInfo := TaskInfo{
		TaskData: taskData,
		TaskFunc: checkKubeConfig,
	}

	// check kubeconfig
	(*worker).QueueTask("check-kubeconfig", zeroDuration, taskInfo)

	taskInfo.TaskFunc = checkProvisioned

	// check provisioned
	(*worker).QueueTask("check-provisioned", zeroDuration, taskInfo)
}

// // InvokeProvisionCheck - Provision 정보를 확인하기 위한 작업 구동
// func (w *worker) InvokeWorkerInvokeProvisionCheck(clusterId, clusterName string) error {
// 	var provisionInfo = ProvisionInfo{
// 		ClusterId:   clusterId,
// 		ClusterName: clusterName,
// 	}

// 	time.After(0)
// 	return w.QueueTask(clusterId, zeroDuration, provisionInfo)

// 	// var input queueTaskInput
// 	// if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
// 	// 	logger.WithError(err).Info("failed to read POST body")
// 	// 	renderResponse(w, http.StatusBadRequest, `{"error": "failed to read POST body"}`)
// 	// 	return
// 	// }
// 	// defer req.Body.Close()

// 	// // parse the work duration from the request body.
// 	// workDuration, errParse := time.ParseDuration(input.WorkDuration)
// 	// if errParse != nil {
// 	// 	logger.WithError(errParse).Info("faile to parse work duration in request")
// 	// 	renderResponse(w, http.StatusBadRequest, `{"error": "failed to parse work duration in request"}`)
// 	// 	return
// 	// }

// 	// // queue the task in background task manager
// 	// if err := h.worker.QueueTask(input.TaskID, workDuration); err != nil {
// 	// 	logger.WithError(err).Info("failed to queue task")
// 	// 	if err == job.ErrWorkerBusy {
// 	// 		w.Header().Set("Retry-After", "60")
// 	// 		renderResponse(w, http.StatusServiceUnavailable, `{"error": "workers are busy, try again later"}`)
// 	// 		return
// 	// 	}
// 	// 	renderResponse(w, http.StatusInternalServerError, `{"error": "failed to queue task"}`)
// 	// 	return
// 	// }

// 	// renderResponse(w, http.StatusAccepted, `{"status": "task queued successfully"}`)
// }
