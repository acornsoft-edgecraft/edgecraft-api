/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package api

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
)

// ProvisioningOpenstackCluster - 오픈스택 클러스터 Provisioning
func ProvisioningOpenstackCluster(cluster *model.OpenstackClusterTable, nodeSets []*model.NodeSetTable) {
	// // Make provision data
	// data := model.OpenstackCAPI{}
	// data.FromTable(cluster, nodeSets)

	// // Processing template
	// temp := template.Must(template.ParseFiles("./conf/templates/capi/openstack_cluster.yaml"))
	// err := temp.Execute(os.Stdout, data)
	// if err != nil {
	// 	logger.Errorf("Openstack Cluster provisioning failed. %v", err)
	// 	return
	// }

	logger.Infof("Openstack Cluster provisioned.")
}
