package model

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	"github.com/gofrs/uuid"
)

// used pointer
type RegisterCloud struct {
	Cloud       model.Cloud          `json:"cloud"`
	Cluster     Cluster              `json:"cluster"`
	Nodes       Nodes                `json:"nodes"`
	EtcdStorage model.EtcdStorage    `json:"etcd_storage"`
	OpenStack   model.CloudOpenStack `json:"openstack"`
}

type Cluster struct {
	CloudClusterUid *uuid.UUID             `json:"cloud_cluster_uid"`
	K8s             model.K8s              `json:"k8s"`
	Baremetal       model.ClusterBaremetal `json:"baremetal"`
}

type Nodes struct {
	model.ClusterNodes
	MasterNode []model.MasterNode `json:"master_nodes"`
	WorkerNode []model.WorkerNode `json:"worker_nodes"`
}
