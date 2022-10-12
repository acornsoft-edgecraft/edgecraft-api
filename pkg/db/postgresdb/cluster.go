/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package postgresdb

import "github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"

// InsertCluster - Insert a new Baremetal Cluster
func (db *DB) InsertCluster(cluster *model.ClusterTable) error {
	return db.GetClient().Insert(cluster)
}

// UpdateCluster - Update a Cluster
func (db *DB) UpdateCluster(cluster *model.ClusterTable) (int64, error) {
	count, err := db.GetClient().Update(cluster)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// GetCluster - 단일 클러스터 조회
func (db *DB) GetCluster(cloudId, clusterId string) (*model.ClusterTable, error) {
	obj, err := db.GetClient().Get(&model.ClusterTable{}, cloudId, clusterId)
	if err != nil {
		return nil, err
	}
	if obj != nil {
		res := obj.(*model.ClusterTable)
		return res, nil
	}
	return nil, nil
}

// GetClusters - 클러스터들 조회
func (db *DB) GetClusters(cloudId string) ([]*model.ClusterTable, error) {
	clusters, err := db.GetClient().Select(&model.ClusterTable{}, getClustersSQL, cloudId)
	if err != nil {
		return nil, err
	}

	var clusterTables []*model.ClusterTable = []*model.ClusterTable{}
	for _, cluster := range clusters {
		clusterTables = append(clusterTables, cluster.(*model.ClusterTable))
	}

	return clusterTables, nil
}

// DeleteCluster - Delete a cluster
func (db *DB) DeleteCluster(cloudId string, clusterId string) (int64, error) {
	cnt, err := db.GetClient().Delete(&model.ClusterTable{CloudUid: &cloudId, ClusterUid: &clusterId})
	if err != nil {
		return -1, err
	}
	return cnt, nil
}

// DeleteCloudClusters - Delete clusters on cloud
func (db *DB) DeleteClusters(cloudId string) (int64, error) {
	result, err := db.GetClient().Exec(deleteClusters, cloudId)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}
