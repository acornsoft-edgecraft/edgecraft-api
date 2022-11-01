/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package api

import (
	"bytes"
	"text/template"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/kubemethod"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
)

// ProvisioningOpenstackCluster - 오픈스택 클러스터 Provisioning
func ProvisioningOpenstackCluster(cluster *model.OpenstackClusterTable, nodeSets []*model.NodeSetTable, k8sVersion string) error {
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

	// 표준 출력으로 처리
	// TODO: 파일 문자열로 K8s에 전송처리
	// TODO: 진행상황을 어떻게 클라이언트에 보여줄 것인가?
	var buff bytes.Buffer
	err = temp.Execute(&buff, data)
	if err != nil {
		logger.Errorf("Template execution failed. cause(%s)", err.Error())
		return err
	}

	logger.Infof("processed templating yaml (%s)", buff.String())

	// 처리된 템플릿을 Kubernetes로 전송
	err = kubemethod.Apply(*cluster.Name, buff.String())
	if err != nil {
		logger.Errorf("Workload cluster creation failed. (cause: %s)", err.Error())
		return err
	}

	//TODO: 진행 상태 검증은 어떤 방식으로?
	//go checkProvisioning(cluster.Name)
	//TODO: 성공된 후에 Provision 상태 갱신

	logger.Infof("Openstack Cluster [%s] provision submitted.", cluster.Name)
	return nil
}
