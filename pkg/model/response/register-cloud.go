package model

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
)

// used pointer
type RegisterCloud struct {
	Cloud       model.Cloud          `json:"cloud"`
	Cluster     Cluster              `json:"cluster"`
	Nodes       Nodes                `json:"nodes"`
	EtcdStorage EtcdStorage          `json:"etcd_storage"`
	OpenStack   model.CloudOpenStack `json:"openstack"`
}

type Cluster struct {
	K8s       model.K8s       `json:"k8s"`
	Baremetal model.Baremetal `json:"baremetal"`
}

type Nodes struct {
	CloudClusterLoadbalancerUse     *bool        `json:"use_loadbalancer"`
	CloudClusterLoadbalancerAddress *string      `json:"loadbalancer_address"`
	CloudClusterLoadbalancerPort    *string      `json:"loadbalancer_port"`
	MasterNode                      []MasterNode `json:"master_nodes"`
	WorkerNode                      []WorkerNode `json:"worker_nodes"`
}

type MasterNode struct {
	Baremetal Baremetal `json:"baremetal"`
	Node      Node      `json:"node"`
}
type WorkerNode struct {
	Baremetal Baremetal `json:"baremetal"`
	Node      Node      `json:"node"`
}

type Baremetal struct {
	HostName             *string `json:"host_name"`
	BmcAddress           *string `json:"bmc_address"`
	BootMacAddress       *string `json:"boot_mac_address"`
	BootMode             *string `json:"boot_mode"`
	OonlinePower         *bool   `json:"online_power"`
	ExternalProvisioning *bool   `json:"external_provisioning"`
}

type Node struct {
	NodeName  *string      `json:"node_name"`
	IpAddress *string      `json:"ip_address"`
	Labels    model.Labels `json:"labels"`
}

type Label struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type EtcdStorage struct {
	Etcd         model.Etcd         `json:"etcd"`
	StorageClass model.StorageClass `json:"storage_class"`
}
