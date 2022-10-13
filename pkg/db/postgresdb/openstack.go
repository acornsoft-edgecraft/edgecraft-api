/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package postgresdb

import "github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"

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
	obj, err := db.GetClient().Get(&model.ClusterTable{}, cloudId, clusterId)
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
func (db *DB) GetOpenstackClusters(cloudId string) ([]*model.OpenstackClusterTable, error) {
	clusters, err := db.GetClient().Select(&model.OpenstackClusterTable{}, getOpenstackClustersSQL, cloudId)
	if err != nil {
		return nil, err
	}

	var clusterTables []*model.OpenstackClusterTable = []*model.OpenstackClusterTable{}
	for _, cluster := range clusters {
		clusterTables = append(clusterTables, cluster.(*model.OpenstackClusterTable))
	}

	return clusterTables, nil
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
	result, err := db.GetClient().Exec(deleteOpenstackClusters, cloudId)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}
