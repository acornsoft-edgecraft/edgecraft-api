/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import (
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
	"github.com/gofrs/uuid"
)

// OpenstackClusterSet - Data for openstack cluster
type OpenstackClusterSet struct {
	Cluster     *OpenstackClusterInfo `json:"cluster"`
	K8s         *KubernetesInfo       `json:"k8s"`
	Openstack   *OpenstackInfo        `json:"openstack"`
	Nodes       *OpenstackNodeSetInfo `json:"nodes"`
	EtcdStorage *EtcdStorageInfo      `json:"etcd_storage"`
}

// ToTable - Openstack Cluster Set을 대상 Table 정보로 매핑 처리
func (ocs *OpenstackClusterSet) ToTable(isUpdate bool, user string, at time.Time) (clusterTable *OpenstackClusterTable, nodesetTables []*NodesetTable) {
	clusterTable = &OpenstackClusterTable{}

	ocs.Cluster.ToTable(clusterTable, isUpdate, user, at)
	ocs.K8s.ToOpenstackTable(clusterTable)

	// TODO: ETCD/STORAGE 정보에 대한 정의 필요함.
	//ocs.EtcdStorage.ToOpenstackTable(clusterTable)

	// NodeSet Table은 Delete & Insert 방식이므로 Update 개념 없음.
	nodesetTables = ocs.Nodes.ToTable(clusterTable, isUpdate, user, at)

	return
}

// OpenstackInfo - Configuration for Openstack
type OpenstackInfo struct {
	ClusterUid          string `json:"cluster_uid"`
	Cloud               string `json:"openstack_cloud"`
	ProviderConfB64     string `json:"openstack_cloud_provider_conf_b64"`
	YamlB64             string `json:"openstack_cloud_yaml_b64"`
	CACertB64           string `json:"openstack_cloud_cacert_b64"`
	NameServers         string `json:"dns_nameservers"`
	FailureDomain       string `json:"failure_domain"`
	ImageName           string `json:"image_name"`
	SSHKeyName          string `json:"ssh_key_name"`
	ExternalNetworkID   string `json:"external_network_id"`
	APIServerFloatingIP string `json:"api_server_floating_ip"`
	UseBastionHost      bool   `json:"use_bastion_host"`
	BastionFlavor       string `json:"bastion_flavor"`
	BastionImageName    string `json:"bastion_image_name"`
	BastionSSHKeyName   string `json:"bastion_ssh_key_name"`
}

// OpenstackClusterInfo - Basic data for openstack cluster
type OpenstackClusterInfo struct {
	ClusterUid string `json:"cluster_uid"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
}

// NewKey - Make new UUID V4
func (osc *OpenstackClusterInfo) NewKey() {
	if osc.ClusterUid == "" {
		osc.ClusterUid = uuid.Must(uuid.NewV4()).String()
	}
}

// ToTable - Openstack Cluster 정보를 Table 정보로 설정
func (osc *OpenstackClusterInfo) ToTable(clusterTable *OpenstackClusterTable, isUpdate bool, user string, at time.Time) {
	if isUpdate {
		clusterTable.Updater = utils.StringPtr(user)
		clusterTable.Updated = utils.TimePtr(at)
	} else {
		osc.NewKey()
		clusterTable.ClusterUid = utils.StringPtr(osc.ClusterUid)
		clusterTable.Creator = utils.StringPtr(user)
		clusterTable.Created = utils.TimePtr(at)
	}

	clusterTable.Name = utils.StringPtr(osc.Name)
	clusterTable.Desc = utils.StringPtr(osc.Desc)
}

// NodeSetInfo - Data for Nodeset
type NodeSetInfo struct {
	Namespace string  `json:"namespace"`
	Name      string  `json:"name"`
	NodeCount int     `json:"node_count"`
	Flavor    string  `json:"flavor"`
	Labels    *Labels `json:"labels"`
}

// OpenstackNodeSetInfo - Data for Nodeset of openstack
type OpenstackNodeSetInfo struct {
	UseLoadbalancer bool           `json:"use_loadbalancer"`
	MasterSets      []*NodeSetInfo `json:"master_sets"`
	WorkerSets      []*NodeSetInfo `json:"worker_sets"`
}

// ToTable - NodeSet 정보를 Openstack 테이블 정보로 설정
func (osnsi *OpenstackNodeSetInfo) ToTable(clusterTable *OpenstackClusterTable, isUpdate bool, user string, at time.Time) []*NodesetTable {
	return nil
}

// FromTable - Openstack 테이블 정보를 NodeSet 정보로 설정
func (osnsi *OpenstackNodeSetInfo) FromTable(clusterTable *OpenstackClusterTable) {

}

// OpenstackClusterList - Cluster list for openstack
type OpenstackClusterList struct {
	ClusterUID string    `json:"cluster_uid"`
	Name       string    `json:"name"`
	Status     int       `json:"status" db:"state"`
	NodeCount  int       `json:"nodeCount"`
	Version    int       `json:"version"`
	Created    time.Time `json:"created" db:"created_at"`
}
