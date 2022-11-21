/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import (
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
	"github.com/gofrs/uuid"
)

// NodeInfo - Data for Node
type NodeInfo struct {
	Name      string  `json:"node_name" example:"sadf"`
	IpAddress string  `json:"ip_address" example:"sadf"`
	Labels    *Labels `json:"labels"`
}

// ToTable - Node 정보를 테이블로 설정
func (ni *NodeInfo) ToTable(nodeTable *NodeTable) {
	nodeTable.Name = utils.StringPtr(ni.Name)
	nodeTable.Ipaddress = utils.StringPtr(ni.IpAddress)
	nodeTable.Labels = ni.Labels
}

// FromTable - 테이블 정보를 Node 정보로 설정
func (ni *NodeInfo) FromTable(nodeTable *NodeTable) {
	ni.Name = *nodeTable.Name
	ni.IpAddress = *nodeTable.Ipaddress
	ni.Labels = nodeTable.Labels
}

// NodeSpecificInfo - Data for Node Spec
type NodeSpecificInfo struct {
	NodeUid       string `json:"node_uid" example:""`
	Type          int    `json:"type" example:"1"`
	BaremetalHost *BaremetalHostInfo
	Node          *NodeInfo
}

// NewKey - Make new UUID V4
func (nsi *NodeSpecificInfo) NewKey() {
	if nsi.NodeUid == "" {
		nsi.NodeUid = uuid.Must(uuid.NewV4()).String()
	}
}

// ToTable - Node Specific 정보를 테이블로 설정
func (nsi *NodeSpecificInfo) ToTable(nodeTable *NodeTable, isUpdate bool, user string, at time.Time) {
	if isUpdate {
		if nodeTable.NodeUid == nil {
			nodeTable.NodeUid = utils.StringPtr(nsi.NodeUid)
		}
		nodeTable.Updater = utils.StringPtr(user)
		nodeTable.Updated = utils.TimePtr(at)
	} else {
		nsi.NewKey()
		nodeTable.NodeUid = utils.StringPtr(nsi.NodeUid)
		nodeTable.Creator = utils.StringPtr(user)
		nodeTable.Created = utils.TimePtr(at)
	}

	nodeTable.Type = utils.IntPrt(nsi.Type)

	nsi.BaremetalHost.ToTable(nodeTable)
	nsi.Node.ToTable(nodeTable)
}

// FromTable - 테이블 정보를 Node Specific 정보로 설정
func (nsi *NodeSpecificInfo) FromTable(nodeTable *NodeTable) {
	nsi.NodeUid = *nodeTable.NodeUid
	nsi.Type = *nodeTable.Type

	nsi.BaremetalHost = &BaremetalHostInfo{}
	nsi.Node = &NodeInfo{}

	nsi.BaremetalHost.FromTable(nodeTable)
	nsi.Node.FromTable(nodeTable)
}

// NodesInfo - Data for Nodes
type NodesInfo struct {
	UseLoadBalancer     bool                `json:"use_loadbalancer" example:"false"`
	LoadBalancerAddress string              `json:"loadbalancer_address" example:""`
	LoadbalancerPort    string              `json:"loadbalancer_port" example:""`
	MasterNodes         []*NodeSpecificInfo `json:"master_nodes"`
	WorkerNodes         []*NodeSpecificInfo `json:"worker_nodes"`
}

// ToTable - Nodes 정보를 테이블로 설정
func (ni *NodesInfo) ToTable(clusterTable *ClusterTable, isUpdate bool, user string, at time.Time) (nodeTables []*NodeTable) {
	clusterTable.LoadbalancerUse = utils.BoolPtr(ni.UseLoadBalancer)
	clusterTable.LoadbalancerAddress = utils.StringPtr(ni.LoadBalancerAddress)
	clusterTable.LoadbalancerPort = utils.StringPtr(ni.LoadbalancerPort)

	// Master Table 구성
	for _, node := range ni.MasterNodes {
		nodeTable := &NodeTable{}
		node.ToTable(nodeTable, isUpdate, user, at)
		nodeTable.Type = utils.IntPrt(1)
		nodeTable.Status = utils.IntPrt(1)
		nodeTables = append(nodeTables, nodeTable)
	}

	// Worker Table 구성
	for _, node := range ni.WorkerNodes {
		nodeTable := &NodeTable{}
		node.ToTable(nodeTable, isUpdate, user, at)
		nodeTable.Type = utils.IntPrt(2)
		nodeTable.Status = utils.IntPrt(1)
		nodeTables = append(nodeTables, nodeTable)
	}

	return
}

// FromTable - 테이블 정보를 Nodes 정보로 설정
func (ni *NodesInfo) FromTable(clusterTable *ClusterTable, nodes []*NodeTable) {
	ni.UseLoadBalancer = *clusterTable.LoadbalancerUse
	ni.LoadBalancerAddress = *clusterTable.LoadbalancerAddress
	ni.LoadbalancerPort = *clusterTable.LoadbalancerPort

	ni.MasterNodes = []*NodeSpecificInfo{}
	ni.WorkerNodes = []*NodeSpecificInfo{}

	for _, node := range nodes {
		var nsi *NodeSpecificInfo = &NodeSpecificInfo{}
		nsi.FromTable(node)

		if *node.Type == common.NodeTypeMaster {
			ni.MasterNodes = append(ni.MasterNodes, nsi)
		} else {
			ni.WorkerNodes = append(ni.WorkerNodes, nsi)
		}
	}
}
