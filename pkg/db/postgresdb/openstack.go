/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package postgresdb

import "github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"

/***********************
 * Openstack Cluster
 ***********************/
// InsertOpenstackCluster - Insert a new Openstack Cluster
func (db *DB) InsertOpenstackCluster(cluster *model.OpenstackClusterTable) error {
	return db.GetClient().Insert(cluster)
}

// UpdateOpenstackCluster - Update a Openstack Cluster
func (db *DB) UpdateOpenstackCluster(cluster *model.OpenstackClusterTable) (int64, error) {
	count, err := db.GetClient().Update(cluster)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// GetOpenstackCluster - Query a Openstack Cluster
func (db *DB) GetOpenstackCluster(cloudId, clusterId string) (*model.OpenstackClusterTable, error) {
	obj, err := db.GetClient().Get(&model.OpenstackClusterTable{}, cloudId, clusterId)
	if err != nil {
		return nil, err
	}
	if obj != nil {
		res := obj.(*model.OpenstackClusterTable)
		return res, nil
	}
	return nil, nil
}

// GetOpenstackClusters - Query all Openstack Clusters
func (db *DB) GetOpenstackClusters(cloudId string) ([]model.OpenstackClusterList, error) {
	var list []model.OpenstackClusterList
	_, err := db.GetClient().Select(&list, getOpenstackClustersSQL, cloudId)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// DeleteOpenstackCluster - Delete a Openstack Cluster
func (db *DB) DeleteOpenstackCluster(cloudId string, clusterId string) (int64, error) {
	cnt, err := db.GetClient().Delete(&model.OpenstackClusterTable{CloudUid: &cloudId, ClusterUid: &clusterId})
	if err != nil {
		return -1, err
	}
	return cnt, nil
}

// DeleteOpenstackClusters - Delete all Openstack Cluster on Cloud
func (db *DB) DeleteOpenstackClusters(cloudId string) (int64, error) {
	result, err := db.GetClient().Exec(deleteOpenstackClustersSQL, cloudId)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

// UpdateOpenstackClusterStatus - Update status of Openstack Cluster
func (db *DB) UpdateOpenstackClusterStatus(cloudId, clusterId string, state int) (int64, error) {
	result, err := db.GetClient().Exec(updateProvisionStatusSQL, cloudId, clusterId, state)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

/***********************
 * NodeSet
 ***********************/

// GetNodeSets - Query all Openstack NodeSets
func (db *DB) GetNodeSets(clusterId string) ([]*model.NodeSetTable, error) {
	nodeSets, err := db.GetClient().Select(&model.NodeSetTable{}, getNodeSetsSQL, clusterId)
	if err != nil {
		return nil, err
	}

	var nodeSetTables []*model.NodeSetTable = []*model.NodeSetTable{}
	for _, nodeSet := range nodeSets {
		nodeSetTables = append(nodeSetTables, nodeSet.(*model.NodeSetTable))
	}

	return nodeSetTables, nil
}

// InsertNodeSet - Insert a new Openstack NodeSet
func (db *DB) InsertNodeSet(nodeSet *model.NodeSetTable) error {
	return db.GetClient().Insert(nodeSet)
}

// UpdateNodeSet - Update a Openstack NodeSet
func (db *DB) UpdateNodeSet(nodeSet *model.NodeSetTable) (int64, error) {
	count, err := db.GetClient().Update(nodeSet)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// GetNodeSet - Query a Openstack NodeSet
func (db *DB) GetNodeSet(clusterId, nodeSetId string) (*model.NodeSetTable, error) {
	obj, err := db.GetClient().Get(&model.NodeSetTable{}, clusterId, nodeSetId)
	if err != nil {
		return nil, err
	}
	if obj != nil {
		res := obj.(*model.NodeSetTable)
		return res, nil
	}
	return nil, nil
}

// DeleteNodeSet - Delete a Openstack NodeSet
func (db *DB) DeleteNodeSet(clusterId, nodeSetId string) (int64, error) {
	cnt, err := db.GetClient().Delete(&model.NodeSetTable{ClusterUid: &clusterId, NodeSetUid: &nodeSetId})
	if err != nil {
		return -1, err
	}
	return cnt, nil
}

// DeleteNodeSets - Delete all Openstack NodeSet on Cluster
func (db *DB) DeleteNodeSets(clusterId string) (int64, error) {
	result, err := db.GetClient().Exec(deleteNodeSetsSQL, clusterId)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

/***********************
 * Benchmarks
 ***********************/

// GetOpenstackBenchmarksList - Query all Openstack Benchmarks
func (db *DB) GetOpenstackBenchmarksList(cloudId, clusterId string) ([]model.OpenstackBenchmarksList, error) {
	var list []model.OpenstackBenchmarksList
	_, err := db.GetClient().Select(&list, getOpenstackBenchmarksListSQL, cloudId, clusterId)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetOpenstackBenchmarks - Query a Openstack Benchmarks
func (db *DB) GetOpenstackBenchmarks(cloudId, clusterId, benchmarkId string) (*model.OpenstackBenchmarksTable, error) {
	obj, err := db.GetClient().Get(&model.OpenstackBenchmarksTable{}, cloudId, clusterId, benchmarkId)
	if err != nil {
		return nil, err
	}
	if obj != nil {
		res := obj.(*model.OpenstackBenchmarksTable)
		return res, nil
	}
	return nil, nil
}

// InsertOpenstackCluster - Insert a new Openstack Cluster
func (db *DB) InsertOpenstackBenchmarks(benchmarks *model.OpenstackBenchmarksTable) error {
	return db.GetClient().Insert(benchmarks)
}
