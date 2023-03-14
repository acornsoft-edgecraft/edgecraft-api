/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package api

import (
	"bytes"
	"errors"
	"path"
	"strings"
	"text/template"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/kubemethod"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/db"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/job"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
)

// replace - 지정한 문자열에서 지정한 문자를 치환한다.
func replace(input, from, to string) string {
	return strings.Replace(input, from, to, -1)
}

// getFunctionalTemplate - 템플릿에 사용할 함수가 설정된 템플릿 반환
func getFunctionalTemplate(filePath string) *template.Template {
	return template.Must(template.New(path.Base(filePath)).Funcs(template.FuncMap{"replace": replace}).ParseFiles(filePath))
}

// getTemplatePath - Bootstrap Provider 정보에 따라 처리할 템플릿 파일을 결정한다.
func getTemplatePath(providerType *common.BootstrapProvider, templateType string) string {
	switch *providerType {
	case common.MicroK8s:
		if templateType == "cluster" {
			return "./conf/templates/capi/openstack_mk8s_cluster.yaml"
		} else if templateType == "upgrade_controlplanes" {
			return "./conf/templates/capi/openstack_mk8s_upgrade_controlplane.yaml"
		} else if templateType == "upgrade_workers" {
			return "./conf/templates/capi/openstack_mk8s_upgrade_worker.yaml"
		} else {
			return "./conf/templates/capi/openstack_mk8s_nodeset.yaml"
		}
	case common.K3s:
		if templateType == "cluster" {
			return "./conf/templates/capi/openstack_k3s_cluster.yaml"
		} else if templateType == "upgrade_controlplanes" {
			return "./conf/templates/capi/openstack_k3s_upgrade_controlplane.yaml"
		} else if templateType == "upgrade_workers" {
			return "./conf/templates/capi/openstack_k3s_upgrade_worker.yaml"
		} else {
			return "./conf/templates/capi/openstack_k3s_nodeset.yaml"
		}
	default:
		if templateType == "cluster" {
			return "./conf/templates/capi/openstack_cluster.yaml"
		} else if templateType == "upgrade_controlplanes" {
			return "./conf/templates/capi/openstack_upgrade_controlplane.yaml"
		} else if templateType == "upgrade_workers" {
			return "./conf/templates/capi/openstack_upgrade_worker.yaml"
		} else {
			return "./conf/templates/capi/openstack_nodeset.yaml"
		}
	}
}

// ProvisioningOpenstackCluster - 오픈스택 클러스터 Provisioning
func ProvisioningOpenstackCluster(worker *job.IWorker, db db.DB, cluster *model.OpenstackClusterTable, nodeSets []*model.NodeSetTable, k8sVersion string) error {
	// Make provision data
	data := model.OpenstackClusterSet{}
	data.FromTable(cluster, nodeSets)
	data.K8s.VersionName = k8sVersion

	// Processing template
	temp, err := template.ParseFiles(getTemplatePath(cluster.BootstrapProvider, "cluster"))
	if err != nil {
		logger.Errorf("Template has errors. cause(%s)", err.Error())
		return err
	}

	// TODO: 진행상황을 어떻게 클라이언트에 보여줄 것인가?
	var buff bytes.Buffer
	err = temp.Execute(&buff, data)
	if err != nil {
		logger.Errorf("Template execution failed. cause(%s)", err.Error())
		return err
	}

	logger.Infof("processed cluster templating yaml (%s)", buff.String())

	// 처리된 템플릿을 Kubernetes로 전송
	err = kubemethod.Apply(*cluster.Name, buff.String())
	if err != nil {
		logger.Errorf("Workload cluster creation failed. (cause: %s)", err.Error())
		return err
	}

	// Update to provisioning status
	// Start. Transaction 얻어옴
	txdb, err := db.BeginTransaction()
	if err != nil {
		return err
	}

	affected, err := db.UpdateOpenstackClusterStatus(*cluster.CloudUid, *cluster.ClusterUid, 2)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			logger.Info("DB rollback Failed.", txErr)
		}
		return err
	} else if affected != 1 {
		logger.Errorf("Cannot update databse to provisioning state (cloud: %s, cluster: %s)", *cluster.CloudUid, *cluster.ClusterUid)
		return errors.New("cannot update database")
	}

	txErr := txdb.Commit()
	if txErr != nil {
		logger.Info("DB commit Failed.", txErr)
	}

	// Provision 검증 작업 등록 (Backgroud)
	err = job.InvokeProvisionCheck(worker, db, *cluster.CloudUid, *cluster.ClusterUid, *cluster.Name, *cluster.Namespace, *cluster.BootstrapProvider)
	if err != nil {
		logger.WithError(err).Infof("Openstack Cluster [%s] provision check job failed.", cluster.Name)
		return err
	}

	logger.Infof("Openstack Cluster [%s] provision submitted.", cluster.Name)
	return nil
}

// ProvisioningOpenstackNodeSet - 오픈스택 NodeSet Provisioning
func ProvisioningOpenstackNodeSet(worker *job.IWorker, db db.DB, cluster *model.OpenstackClusterTable, nodeSets []*model.NodeSetTable, k8sVersion string) error {
	// Make provision data
	data := model.OpenstackClusterSet{}
	data.FromTable(cluster, nodeSets)
	data.K8s.VersionName = k8sVersion

	// Processing template
	temp, err := template.ParseFiles(getTemplatePath(cluster.BootstrapProvider, "nodeset"))
	if err != nil {
		logger.Errorf("Template has errors. cause(%s)", err.Error())
		return err
	}

	// TODO: 진행상황을 어떻게 클라이언트에 보여줄 것인가?
	var buff bytes.Buffer
	err = temp.Execute(&buff, data)
	if err != nil {
		logger.Errorf("Template execution failed. cause(%s)", err.Error())
		return err
	}

	logger.Infof("processed nodeset templating yaml (%s)", buff.String())

	// 처리된 템플릿을 Kubernetes로 전송
	err = kubemethod.Apply(*cluster.Name, buff.String())
	if err != nil {
		logger.Errorf("NodeSet creation failed. (cause: %s)", err.Error())
		return err
	}

	logger.Infof("Openstack Cluseter [%s] - NodeSet provision submitted.", cluster.Name)
	return nil
}

// K8sVersionUpgradingOpenstackCluster - 오픈스택 클러스터 K8s Version Upgrading
func K8sVersionUpgradingOpenstackCluster(worker *job.IWorker, database db.DB, cluster *model.OpenstackClusterTable, nodeSets []*model.NodeSetTable, k8sVersion string, upgradeInfo *model.K8sUpgradeInfo) error {
	// Make provision data
	data := model.OpenstackClusterSet{}
	data.FromTable(cluster, nodeSets)
	data.K8s.VersionName = k8sVersion
	data.Openstack.ImageName = upgradeInfo.Image

	// Processing Template for control plane
	temp := getFunctionalTemplate(getTemplatePath(cluster.BootstrapProvider, "upgrade_controlplanes"))
	var controlPlanesBuff bytes.Buffer
	err := temp.Execute(&controlPlanesBuff, data)
	if err != nil {
		logger.Errorf("Template execution failed [upgrade controlplanes]. cause(%s)", err.Error())
		return err
	}
	logger.Infof("processed control-plane upgrade templating yaml (%s)", controlPlanesBuff.String())

	// Processing Template for workers
	temp = getFunctionalTemplate(getTemplatePath(cluster.BootstrapProvider, "upgrade_workers"))
	var workersBuff bytes.Buffer
	err = temp.Execute(&workersBuff, data)
	if err != nil {
		logger.Errorf("Template execution failed [upgrade workers]. cause(%s)", err.Error())
		return err
	}
	logger.Infof("processed worker upgrade templating yaml (%s)", workersBuff.String())

	// K8sVersionUpgrade 실행 (Background)
	masterSetName := data.Cluster.Name + "-" + data.Nodes.MasterSets[0].Name
	err = job.InvokeK8sVersionUpgrade(worker, database, data.Cluster.Name, data.Cluster.Namespace, masterSetName, controlPlanesBuff.String(), workersBuff.String())
	if err != nil {
		logger.WithError(err).Infof("JOB::Upgrading kubernetest version on Openstack Cluster [%s] failed.", *cluster.Name)
		return err
	}

	return nil
}
