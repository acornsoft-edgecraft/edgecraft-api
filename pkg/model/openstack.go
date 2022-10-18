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
func (ocs *OpenstackClusterSet) ToTable(cloudId string, isUpdate bool, user string, at time.Time) (clusterTable *OpenstackClusterTable, nodesetTables []*NodeSetTable) {
	clusterTable = &OpenstackClusterTable{}

	clusterTable.CloudUid = utils.StringPtr(cloudId)
	ocs.Cluster.ToTable(clusterTable, isUpdate, user, at)
	ocs.K8s.ToOpenstackTable(clusterTable)

	// TODO: ETCD/STORAGE 정보에 대한 정의 필요함.
	ocs.EtcdStorage.ToOpenstackTable(clusterTable)

	// NodeSet Table은 Delete & Insert 방식이므로 Update 개념 없음.
	nodesetTables = ocs.Nodes.ToTable(clusterTable, isUpdate, user, at)

	return
}

// OpenstackInfo - Configuration for Openstack
type OpenstackInfo struct {
	Cloud               string `json:"openstack_cloud" example:"openstack"`
	ProviderConfB64     string `json:"openstack_cloud_provider_conf_b64" example:"W0dsb2JhbF0KYXV0aC11cmw9aHR0cDovLzE5Mi4xNjguNzcuMTEvaWRlbnRpdHkKdXNlcm5hbWU9InN1bm1pIgpwYXNzd29yZD0iZmtmZms0NDgiCnRlbmFudC1pZD0iNTQyZTdhMDRmNjkxNDgyOWI0M2U3N2Y5ZWYxMmI3NzkiCnRlbmFudC1uYW1lPSJlZGdlY3JhZnQiCmRvbWFpbi1uYW1lPSJEZWZhdWx0IgpyZWdpb249IlJlZ2lvbk9uZSIK"`
	YamlB64             string `json:"openstack_cloud_yaml_b64" example:"Y2xvdWRzOgogIG9wZW5zdGFjazoKICAgIGF1dGg6CiAgICAgIGF1dGhfdXJsOiBodHRwOi8vMTkyLjE2OC43Ny4xMS9pZGVudGl0eQogICAgICB1c2VybmFtZTogInN1bm1pIgogICAgICBwYXNzd29yZDogImZrZmZrNDQ4IgogICAgICBwcm9qZWN0X2lkOiA1NDJlN2EwNGY2OTE0ODI5YjQzZTc3ZjllZjEyYjc3OQogICAgICBwcm9qZWN0X25hbWU6ICJlZGdlY3JhZnQiCiAgICAgIHVzZXJfZG9tYWluX25hbWU6ICJEZWZhdWx0IgogICAgcmVnaW9uX25hbWU6ICJSZWdpb25PbmUiCiAgICBpbnRlcmZhY2U6ICJwdWJsaWMiCiAgICBpZGVudGl0eV9hcGlfdmVyc2lvbjogMwo="`
	CACertB64           string `json:"openstack_cloud_cacert_b64" example:"Cg=="`
	DNSNameServers      string `json:"dns_nameservers" example:"168.126.63.1"`
	FailureDomain       string `json:"failure_domain" example:""` // nova
	ImageName           string `json:"image_name" example:"ubuntu-2004-kube-v1.23.3"`
	SSHKeyName          string `json:"ssh_key_name" example:"sunmi"`
	ExternalNetworkID   string `json:"external_network_id" example:""` // public
	APIServerFloatingIP string `json:"api_server_floating_ip" example:""`
	NodeCidr            string `json:"node_cidr" example:"10.96.0.0/24"`
	UseBastionHost      bool   `json:"use_bastion_host" example:"false"`
	BastionFlavor       string `json:"bastion_flavor" example:""`
	BastionImageName    string `json:"bastion_image_name" example:""`
	BastionSSHKeyName   string `json:"bastion_ssh_key_name" example:""`
	BastionFloatingIP   string `json:"bastion_floating_ip" example:""`
}

// OpenstackClusterInfo - Basic data for openstack cluster
type OpenstackClusterInfo struct {
	ClusterUid string `json:"cluster_uid" example:""`
	Name       string `json:"name" example:"os-cluster-#1"`
	Desc       string `json:"desc" example:"Openstack Test Cluster #1"`
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
		clusterTable.Status = utils.IntPrt(1)
	}

	clusterTable.Name = utils.StringPtr(osc.Name)
	clusterTable.Desc = utils.StringPtr(osc.Desc)
}

// NodeSetInfo - Data for Nodeset
type NodeSetInfo struct {
	NodeSetUid string  `json:"nodeset_uid" example:""`
	Namespace  string  `json:"namespace" example:""`
	Name       string  `json:"name" example:""`
	NodeCount  int     `json:"node_count" example:"1"`
	Flavor     string  `json:"flavor" example:"m1.medium"`
	Labels     *Labels `json:"labels"`
}

// NewKey - Make new UUID V4
func (nsi *NodeSetInfo) NewKey() {
	if nsi.NodeSetUid == "" {
		nsi.NodeSetUid = uuid.Must(uuid.NewV4()).String()
	}
}

// ToTable - NodeSet 정보를 테이블 정보로 설정
func (nsi *NodeSetInfo) ToTable(nodeSetTable *NodeSetTable, isUpdate bool, user string, at time.Time) {
	if isUpdate {
		if nodeSetTable.NodeSetUid == nil {
			nodeSetTable.NodeSetUid = utils.StringPtr(nsi.NodeSetUid)
		}
		nodeSetTable.Updater = utils.StringPtr(user)
		nodeSetTable.Updated = utils.TimePtr(at)
	} else {
		nsi.NewKey()
		nodeSetTable.NodeSetUid = utils.StringPtr(nsi.NodeSetUid)
		nodeSetTable.Creator = utils.StringPtr(user)
		nodeSetTable.Created = utils.TimePtr(at)
	}

	nodeSetTable.Namespace = utils.StringPtr(nsi.Namespace)
	nodeSetTable.Name = utils.StringPtr(nsi.Name)
	nodeSetTable.NodeCount = utils.IntPrt(nsi.NodeCount)
	nodeSetTable.Flavor = utils.StringPtr(nsi.Flavor)
	nodeSetTable.Labels = nsi.Labels
}

// FromTable - 테이블 정보를 NodeSet 정보로 설정
func (nsi *NodeSetInfo) FromTable(nodeSetTable *NodeSetTable) {
	nsi.NodeSetUid = *nodeSetTable.NodeSetUid
	nsi.Namespace = *nodeSetTable.Namespace
	nsi.Name = *nodeSetTable.Name
	nsi.NodeCount = *nodeSetTable.NodeCount
	nsi.Flavor = *nodeSetTable.Flavor
	nsi.Labels = nodeSetTable.Labels
}

// OpenstackNodeSetInfo - Data for Nodeset of openstack
type OpenstackNodeSetInfo struct {
	UseLoadbalancer bool           `json:"use_loadbalancer" example:"false"`
	MasterSets      []*NodeSetInfo `json:"master_sets"`
	WorkerSets      []*NodeSetInfo `json:"worker_sets"`
}

// ToTable - NodeSet 정보를 Openstack 테이블 정보로 설정
func (osnsi *OpenstackNodeSetInfo) ToTable(clusterTable *OpenstackClusterTable, isUpdate bool, user string, at time.Time) (nodeSetTables []*NodeSetTable) {
	clusterTable.LoadbalancerUse = utils.BoolPtr(osnsi.UseLoadbalancer)

	// MasterSet 구성
	for _, nodeSet := range osnsi.MasterSets {
		nodeSetTable := &NodeSetTable{}
		nodeSet.ToTable(nodeSetTable, isUpdate, user, at)
		nodeSetTable.ClusterUid = clusterTable.ClusterUid
		nodeSetTable.Type = utils.IntPrt(1)
		nodeSetTables = append(nodeSetTables, nodeSetTable)
	}

	// WorkerSet 구성
	for _, nodeSet := range osnsi.WorkerSets {
		nodeSetTable := &NodeSetTable{}
		nodeSet.ToTable(nodeSetTable, isUpdate, user, at)
		nodeSetTable.ClusterUid = clusterTable.ClusterUid
		nodeSetTable.Type = utils.IntPrt(2)
		nodeSetTables = append(nodeSetTables, nodeSetTable)
	}

	return
}

// FromTable - Openstack 테이블 정보를 NodeSet 정보로 설정
func (osnsi *OpenstackNodeSetInfo) FromTable(clusterTable *OpenstackClusterTable, nodeSetTables []*NodeSetTable) {
	osnsi.UseLoadbalancer = *clusterTable.LoadbalancerUse
	osnsi.MasterSets = []*NodeSetInfo{}
	osnsi.WorkerSets = []*NodeSetInfo{}

	for _, nodeSetTable := range nodeSetTables {
		var nsi *NodeSetInfo = &NodeSetInfo{}
		nsi.FromTable(nodeSetTable)

		if *nodeSetTable.Type == 1 {
			osnsi.MasterSets = append(osnsi.MasterSets, nsi)
		} else {
			osnsi.WorkerSets = append(osnsi.WorkerSets, nsi)
		}
	}
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
