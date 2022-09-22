/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
	"github.com/gofrs/uuid"
)

// NodeInfo - Data for Node
type NodeInfo struct {
	NodeName  string  `json:"node_name" example:"sadf"`
	IpAddress string  `json:"ip_address" example:"sadf"`
	Labels    *Labels `json:"labels"`
}

// ToTable - Node 정보를 테이블로 설정
func (ni *NodeInfo) ToTable(nodeTable *NodeTable) {
	nodeTable.Name = ni.NodeName
	nodeTable.Ipaddress = ni.IpAddress
	nodeTable.Labels = ni.Labels
}

// NodeSpecificInfo - Data for Node Spec
type NodeSpecificInfo struct {
	NodeUid       string `json:"node_uid" example:""`
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
func (nsi *NodeSpecificInfo) ToTable(nodeTable *NodeTable) {
	nsi.NewKey()

	nsi.BaremetalHost.ToTable(nodeTable)
	nsi.Node.ToTable(nodeTable)

	nodeTable.NodeUid = nsi.NodeUid
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
func (ni *NodesInfo) ToTable(clusterTable *ClusterTable) (nodeTables []*NodeTable) {
	utils.CopyTo(&clusterTable, ni)

	// Master Table 구성
	for _, node := range ni.MasterNodes {
		nodeTable := &NodeTable{}
		node.ToTable(nodeTable)
		nodeTable.Type = "1"
		nodeTable.Status = "1"
		nodeTables = append(nodeTables, nodeTable)
	}

	// Worker Table 구성
	for _, node := range ni.WorkerNodes {
		nodeTable := &NodeTable{}
		node.ToTable(nodeTable)
		nodeTable.Type = "2"
		nodeTable.Status = "1"
		nodeTables = append(nodeTables, nodeTable)
	}

	return
}
