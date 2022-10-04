/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package postgresdb

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
)

// findNodesByKey - 지정된 조건에 맞는 Node 정보 조회
func (db *DB) findNodesByKey(cloudId, clusterId, nodeId *string) ([]*model.NodeTable, error) {
	var nodes []*model.NodeTable
	var nodeCond *model.NodeTable = &model.NodeTable{}

	if cloudId != nil {
		nodeCond.CloudUid = cloudId
	}
	if clusterId != nil {
		nodeCond.ClusterUid = clusterId
	}
	if nodeId != nil {
		nodeCond.NodeUid = nodeId
	}

	sql, err := utils.RenderTmpl("getNodeByCond", getNodeSQL, nodeCond)
	if err != nil {
		return nil, err
	}

	args, _ := utils.StructToMap(nodeCond)
	_, err = db.GetClient().Select(&nodes, sql, args)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

// deleteNodesByKey - 지정된 조건에 맞는 Node 정보 삭제
func (db *DB) deleteNodesByKey(cloudId, clusterId, nodeId *string) (int64, error) {
	var nodeCond *model.NodeTable = &model.NodeTable{}

	if cloudId != nil {
		nodeCond.CloudUid = cloudId
	}
	if clusterId != nil {
		nodeCond.ClusterUid = clusterId
	}
	if nodeId != nil {
		nodeCond.NodeUid = nodeId
	}

	sql, err := utils.RenderTmpl("deleteNodeByCond", deleteNodeSQL, nodeCond)
	if err != nil {
		return -1, err
	}

	args, _ := utils.StructToMap(nodeCond)
	result, err := db.GetClient().Exec(sql, args)
	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}

// GetNodes - Nodes 조회
func (db *DB) GetNodes(cloudId, clusterId string) ([]*model.NodeTable, error) {
	return db.findNodesByKey(&cloudId, &clusterId, nil)
}

// GetNodesByCloud - Cloud에 속한 Nodes 조회
func (db *DB) GetNodesByCloud(cloudId string) ([]*model.NodeTable, error) {
	return db.findNodesByKey(&cloudId, nil, nil)
}

// GetNode - 단일 Node 조회
func (db *DB) GetNode(cloudId, clusterId, nodeId string) (*model.NodeTable, error) {
	nodes, err := db.findNodesByKey(&cloudId, &clusterId, &nodeId)
	if err != nil {
		return nil, err
	}
	if len(nodes) > 0 {
		return nodes[0], nil
	}

	return nil, nil
}

// GetNodeByCloud - Cloud에 속한 단일 Node 조회
func (db *DB) GetNodeByCloud(cloudId, nodeId string) (*model.NodeTable, error) {
	nodes, err := db.findNodesByKey(&cloudId, nil, &nodeId)
	if err != nil {
		return nil, err
	}
	if len(nodes) > 0 {
		return nodes[0], nil
	}

	return nil, nil
}

// InsertNode - Insert a new Baremetal Node
func (db *DB) InsertNode(node *model.NodeTable) error {
	return db.GetClient().Insert(node)
}

// UpdateNode - Update a Node
func (db *DB) UpdateNode(node *model.NodeTable) (int64, error) {
	count, err := db.GetClient().Update(node)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// DeleteNode - Delete a node by node uid
func (db *DB) DeleteNode(nodeId string) (int64, error) {
	return db.deleteNodesByKey(nil, nil, &nodeId)
}

// DeleteNodes - Delete a nodes by cloud/cluster uid
func (db *DB) DeleteNodes(cloudId, clusterId string) (int64, error) {
	return db.deleteNodesByKey(&cloudId, &clusterId, nil)
}

// DeleteCloudNodes - Delete all nodes on cloud
func (db *DB) DeleteCloudNodes(cloudId string) (int64, error) {
	return db.deleteNodesByKey(&cloudId, nil, nil)
}
