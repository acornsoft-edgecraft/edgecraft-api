/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package api

import (
	"html/template"
	"os"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
)

// ProvisioningOpenstackCluster - 오픈스택 클러스터 Provisioning
func ProvisioningOpenstackCluster(cluster *model.OpenstackClusterTable, nodeSets []*model.NodeSetTable, k8sVersion string) {
	// Make provision data
	data := model.OpenstackClusterSet{}
	data.FromTable(cluster, nodeSets)
	data.K8s.VersionName = k8sVersion

	// Processing template
	temp, err := template.ParseFiles("./conf/templates/capi/openstack_cluster.yaml")
	if err != nil {
		logger.Errorf("Template has errors. %v", err)
		return
	}

	// 표준 출력으로 처리
	// TODO: 파일 문자열로 K8s에 전송처리
	// TODO: 진행상황을 어떻게 클라이언트에 보여줄 것인가?
	err = temp.Execute(os.Stdout, data)
	if err != nil {
		logger.Errorf("Template execution failed. %v", err)
		return
	}

	//TODO: Submit yaml to kubernetes

	//TODO: 성공된 후에 Provision 상태 갱신

	logger.Infof("Openstack Cluster provisioned.")
}
