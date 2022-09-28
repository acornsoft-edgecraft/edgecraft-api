/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package postgresdb

import "github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"

const getNodesSQL = `
SELECT *
FROM tbl_cloud_node c
WHERE
	cloud_uid = $1
AND	cluster_uid = $2
`

const deleteNodesSQL = `
DELETE
FROM tbl_cloud_node A
WHERE
	A.cloud_uid = $1
AND	A.cluster_uid = $2
`

const deleteCloudNodesSQL = `
DELETE
FROM tbl_cloud_node A
WHERE
	A.cloud_uid = $1
`

// GetNode - 단일 Node 조회
func (db *DB) GetNode(cloudUid, clusterUid, nodeUid string) (*model.NodeTable, error) {
	node, err := db.GetClient().Get(&model.NodeTable{}, cloudUid, clusterUid, nodeUid)
	if err != nil {
		return nil, err
	}
	if node != nil {
		return node.(*model.NodeTable), nil
	}

	return nil, nil
}

// SelectNodes - Nodes 조회
func (db *DB) SelectNodes(cloudUid, clusterUid string) ([]*model.NodeTable, error) {
	nodes, err := db.GetClient().Select(&model.NodeTable{}, getNodesSQL, cloudUid, clusterUid)
	if err != nil {
		return nil, err
	}

	var nodeTables []*model.NodeTable = []*model.NodeTable{}
	for _, node := range nodes {
		nodeTables = append(nodeTables, node.(*model.NodeTable))
	}

	return nodeTables, nil
}

// InsertNode - Insert a new Baremetal Node
func (db *DB) InsertNode(node *model.NodeTable) error {
	return db.GetClient().Insert(node)
}

// DeleteNodes - Delete a nodes by cloud/cluster uid
func (db *DB) DeleteNodes(cloudUid, clusterUid string) (int64, error) {
	result, err := db.GetClient().Exec(deleteNodesSQL, cloudUid, clusterUid)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

// DeleteCloudNodes - Delete all nodes on cloud
func (db *DB) DeleteCloudNodes(cloudUid string) (int64, error) {
	result, err := db.GetClient().Exec(deleteCloudNodesSQL, cloudUid)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}
