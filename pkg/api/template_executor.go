/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package api

import (
	"bytes"
	"errors"
	"text/template"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/kubemethod"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/db"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/job"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
)

// ProvisioningOpenstackCluster - 오픈스택 클러스터 Provisioning
func ProvisioningOpenstackCluster(worker *job.IWorker, db db.DB, cluster *model.OpenstackClusterTable, nodeSets []*model.NodeSetTable, k8sVersion string) error {
	// Make provision data
	data := model.OpenstackClusterSet{}
	data.FromTable(cluster, nodeSets)
	data.K8s.VersionName = k8sVersion

	// Processing template
	temp, err := template.ParseFiles("./conf/templates/capi/openstack_cluster.yaml")
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
	err = job.InvokeProvisionCheck(worker, db, *cluster.CloudUid, *cluster.ClusterUid, *cluster.Name, *cluster.Namespace)
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
	temp, err := template.ParseFiles("./conf/templates/capi/openstack_nodeset.yaml")
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
