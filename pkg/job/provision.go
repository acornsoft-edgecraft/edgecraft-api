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
	BootstrapProvider common.BootstrapProvider
	CustomData        interface{}
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
			if data.BootstrapProvider == common.MicroK8s {
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
func InvokeProvisionCheck(worker *IWorker, db db.DB, cloudId, clusterId, clusterName, namespace string, bootstrapProvider common.BootstrapProvider) error {
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

// applyRbacRoles - Apply RBAC Roles for microk8s
func applyRbacRoles(task string, taskData interface{}) {
	data := taskData.(*TaskData)
	retryTicker := time.NewTicker(30 * time.Second)
	count := 0
	for range retryTicker.C {
		count += 1
		if count > 100 {
			retryTicker.Stop()
			logger.WithField("task", task).Info("Close apply RBAC Roles. retry count (100 times in hour) over.")
			return
		}

		clusterTable, err := data.Database.GetOpenstackCluster(data.CloudId, data.ClusterId)
		if err != nil {
			logger.WithField("task", task).WithError(err).Infof("Get Openstack cluster info failed. (cause: %s)", err.Error())
			retryTicker.Stop()
			return
		}
		if clusterTable == nil {
			logger.WithField("task", task).WithError(err).Infof("Get Openstack cluster info failed. (cause: %s)", err.Error())
			retryTicker.Stop()
			return
		}
		// 클러스터 상태 조회
		if *clusterTable.Status != common.StatusProvisioned {
			logger.WithField("task", task).Warnf("Openstack Cluster Status is not provisioned (status: %s)", *clusterTable.Status)
			continue
		}

		// ClusterRole 확인
		clusterRoleName := "system:kube-apiserver-to-kubelet"
		roleExists, err := kubemethod.ExistsClusterRole(data.ClusterName, clusterRoleName)
		if err != nil {
			logger.WithField("task", task).Warnf("Exists ClusterRole failed. (cause: %s)", err.Error())
			continue
		}

		// ClusterRoleBinding 확인
		clusterRBName := "system:kube-apiserver"
		rbExists, err := kubemethod.ExistsClusterRoleBinding(data.ClusterName, clusterRBName)
		if err != nil {
			logger.WithField("task", task).Warnf("Exists ClusterRoleBinding failed. (cause: %s)", err.Error())
			continue
		}

		if roleExists && rbExists {
			logger.WithField("task", task).Infof("ClusterRole & ClusterRoleBinding already exists.")
			retryTicker.Stop()
			return
		} else {
			logger.WithField("task", task).Info("Apply RBAC Roles manifest.")
			// Apply rbac roles
			err = kubemethod.Apply(data.ClusterName, data.CustomData.(string))
			if err != nil {
				logger.WithField("task", task).Warnf("Apply RBAC Roles manifest failed. (cause: %s)", err.Error())
				continue
			}
			logger.WithField("task", task).Info("Apply RBAC Roles manifest completed.")
		}
	}
}

// InvokeProvisioned - Providion 완료 후 작업
func InvokeProvisioned(worker *IWorker, db db.DB, cloudId, clusterId, clusterName, manifest string) error {
	taskData := &TaskData{
		Database:    db,
		CloudId:     cloudId,
		ClusterId:   clusterId,
		ClusterName: clusterName,
		CustomData:  manifest,
	}

	taskInfo := TaskInfo{
		TaskData: taskData,
		TaskFunc: applyRbacRoles,
	}

	// Apply RBAC Roles for microk8s
	err := (*worker).QueueTask("apply-rbac-roles", zeroDuration, taskInfo)
	if err != nil {
		logger.Warnf("Failed apply-rbac-roles, err: %v", err)
		return err
	}

	return nil
}
