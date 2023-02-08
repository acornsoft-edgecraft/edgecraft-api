/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package job

import (
	"strings"
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/kubemethod"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/config"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/db"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
)

const (
	zeroDuration time.Duration = 0
)

// TaskData - Provision에 사용할 클러스터 식별 데이터
type TaskData struct {
	Database          db.DB
	CloudId           string
	ClusterId         string
	ClusterName       string
	Namespace         string
	BootstrapProvider int
}

// adjustKubeconfigForMicroK8s - MicroK8s인 경우에 Kubeconfig 내용을 Cluster Name 기준으로 조정
func adjustKubeconfigForMicroK8s(kubeconfig string, clusterName string) string {
	// Replace cluster name microk8s-cluster to clusterName
	kubeconfig = strings.Replace(kubeconfig, "microk8s-cluster", clusterName, -1)
	// Replace user name admin to clusterName-admin
	kubeconfig = strings.Replace(kubeconfig, "admin", clusterName+"-admin", -1)
	// Replace context name microk8s to clusterName-admin@clusterName
	kubeconfig = strings.Replace(kubeconfig, "microk8s", clusterName+"-admin@"+clusterName, -1)

	return kubeconfig
}

// checkProvisionKubeConfig - Provision 처리 중인 클러스터의 kubeconfig 정보 처리
func checkProvisionKubeConfig(task string, taskData interface{}) {
	data := taskData.(*TaskData)
	retryTicker := time.NewTicker(30 * time.Second)
	count := 0
	for range retryTicker.C {
		// Get kubeconfig for workload cluster
		kubeconfig, err := kubemethod.GetKubeconfig(data.Namespace, data.ClusterName, "value")
		if err != nil {
			logger.WithField("task", task).WithError(err).Infof("Retrieve kubeconfig for (%s) failed.", data.ClusterName)
		} else {
			// FIX: MicroK8s Kubeconfig 조정
			if data.BootstrapProvider == 2 {
				kubeconfig = adjustKubeconfigForMicroK8s(kubeconfig, data.ClusterName)
				logger.WithField("test", task).Info("Adjust microk8s kubeconfig for (%s-%d)", data.ClusterName, data.BootstrapProvider)
			}

			// Add cluster's kubeconfig
			err = config.HostCluster.Add([]byte(kubeconfig))
			if err != nil {
				logger.WithField("task", task).WithError(err).Infof("Add kubeconfig for (%s) to configmap failed.", data.ClusterName)
			} else {
				logger.WithField("task", task).Info("Checking kubeconfig and added")
				retryTicker.Stop()
				return
			}
		}

		count += 1
		if count > 120 {
			retryTicker.Stop()
			logger.WithField("task", task).Info("End the check of kubeconfig. exceeding the number of retries (120 times in 1 hour).")
			return
		}
	}
}

// checkDeletedKubeConfig - 삭제되는 클러스터의 kubeconfig 삭제 처리
func checkDeletedKubeConfig(task string, taskData interface{}) {
	data := taskData.(*TaskData)
	retryTicker := time.NewTicker(30 * time.Second)
	count := 0

	for range retryTicker.C {
		// Get kubeconfig for workload cluster
		_, err := kubemethod.GetKubeconfig(data.Namespace, data.ClusterName, "value")
		if err != nil {
			if utils.CheckK8sNotFound(err) {
				// Remove cluster's kubeconfig
				err = config.HostCluster.Remove(data.ClusterName)
				if err != nil {
					logger.WithField("task", task).WithError(err).Infof("Remove kubeconfig for (%s) from configmap or file failed.", data.ClusterName)
				} else {
					logger.WithField("task", task).Info("Checking kubeconfig and removed")
					retryTicker.Stop()
					return
				}
			}
			logger.WithField("task", task).WithError(err).Infof("Retrieve kubeconfig for (%s) failed.", data.ClusterName)
		}

		count += 1
		if count > 120 {
			retryTicker.Stop()
			logger.WithField("task", task).Info("End the delete check of kubeconfig. exceeding the number of retries (120 times in 1 hour).")
			return
		}
	}
}

// checkProvisionCluster - Provision 처리 확인 및 후처리
func checkProvisionCluster(task string, taskData interface{}) {
	data := taskData.(*TaskData)
	retryTicker := time.NewTicker(30 * time.Second)
	count := 0
	for range retryTicker.C {
		// Get phase for workload cluster
		phase, err := kubemethod.GetProvisionPhase(data.Namespace, data.ClusterName)
		if err != nil {
			logger.WithField("task", task).WithError(err).Infof("Retrieve provision status for (%s) failed.", data.ClusterName)
		} else {
			var state int = common.StatusSaved
			var provisionState = strings.ToLower(phase)
			if provisionState == "provisioned" {
				state = common.StatusProvisioned
			} else if provisionState == "failed" || provisionState == "pending" || provisionState == "Unknown" {
				state = common.StatusFailed
			} else if provisionState == "provisioning" {
				state = common.StatusProvisioning
			}

			logger.WithField("task", task).Infof("Checked state (%d), phase (%s), cluster (%s)", state, phase, data.ClusterName)

			if provisionState != "" && provisionState != "provisioning" {
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

				logger.WithField("task", task).Infof("Checking provisioned and update to database [state: %d, cluster: %s]", state, data.ClusterName)
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

// checkDeleteCluster - 클러스터의 삭제 여부 검증 및 후처리
func checkDeleteCluster(task string, taskData interface{}) {
	data := taskData.(*TaskData)
	retryTicker := time.NewTicker(30 * time.Second)
	count := 0
	for range retryTicker.C {
		// Get phase for workload cluster
		_, err := kubemethod.GetProvisionPhase(data.Namespace, data.ClusterName)
		if err != nil {
			if utils.CheckK8sNotFound(err) {
				// Deleted, update database. deleted
				affected, err := data.Database.UpdateOpenstackClusterStatus(data.CloudId, data.ClusterId, common.StatusDeleted)
				if err != nil {
					logger.WithField("task", task).WithError(err).Infof("Update deleted state (%d) for (%s) failed.", common.StatusDeleted, data.ClusterName)
					retryTicker.Stop()
					return
				} else if affected != 1 {
					logger.WithField("task", task).WithError(err).Infof("Close checking for delete state (%d) for (%s) failed. (data not found, check cloud/cluster id)", common.StatusDeleted, data.ClusterName)
					retryTicker.Stop()
					return
				}

				logger.WithField("task", task).Infof("Checking deleted and update to database [state: %d, cluster: %s]", common.StatusDeleted, data.ClusterName)
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
func InvokeProvisionCheck(worker *IWorker, db db.DB, cloudId, clusterId, clusterName, namespace string, bootstrapProvider int) error {
	taskData := &TaskData{
		Database:          db,
		CloudId:           cloudId,
		ClusterId:         clusterId,
		ClusterName:       clusterName,
		BootstrapProvider: bootstrapProvider,
		Namespace:         namespace,
	}

	taskInfo := TaskInfo{
		TaskData: taskData,
		TaskFunc: checkProvisionKubeConfig,
	}

	// check kubeconfig
	err := (*worker).QueueTask("check-kubeconfig", zeroDuration, taskInfo)
	if err != nil {
		return err
	}

	taskInfo.TaskFunc = checkProvisionCluster

	// check provisioned
	err = (*worker).QueueTask("check-provisioned", zeroDuration, taskInfo)
	if err != nil {
		return err
	}

	return nil
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

// InvokeDeleteCheck - 프로비전된 클러스터의 삭제에 대한 진행 검증 작업
func InvokeDeleteCheck(worker *IWorker, db db.DB, cloudId, clusterId, clusterName, namespace string) error {
	taskData := &TaskData{
		Database:    db,
		CloudId:     cloudId,
		ClusterId:   clusterId,
		ClusterName: clusterName,
		Namespace:   namespace,
	}

	taskInfo := TaskInfo{
		TaskData: taskData,
		TaskFunc: checkDeletedKubeConfig,
	}

	// check kubeconfig
	err := (*worker).QueueTask("check-deleted-kubeconfig", zeroDuration, taskInfo)
	if err != nil {
		return err
	}

	taskInfo.TaskFunc = checkDeleteCluster

	// check provisioned
	err = (*worker).QueueTask("check-deleted-cluster", zeroDuration, taskInfo)
	if err != nil {
		return err
	}

	return nil
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
