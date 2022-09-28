/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package postgresdb

import "github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"

const getClustersSQL = `
SELECT 
	A.*
FROM 
	"edgecraft"."tbl_cloud_cluster" A
WHERE
	A.cloud_uid = $1
`

const deleteCloudClusters = `
DELETE
FROM "edgecraft"."tbl_cloud_cluster" A
WHERE
	A.cloud_uid = $1
`

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
func (db *DB) GetCluster(cloudUid, clusterUid string) (*model.ClusterTable, error) {
	obj, err := db.GetClient().Get(&model.ClusterTable{}, cloudUid, clusterUid)
	if err != nil {
		return nil, err
	}
	if obj != nil {
		res := obj.(*model.ClusterTable)
		return res, nil
	}
	return nil, nil
}

// SelectClusters - 클러스터들 조회
func (db *DB) SelectClusters(cloudUid string) ([]*model.ClusterTable, error) {
	clusters, err := db.GetClient().Select(&model.ClusterTable{}, getClustersSQL, cloudUid)
	if err != nil {
		return nil, err
	}

	var clusterTables []*model.ClusterTable = []*model.ClusterTable{}
	for _, cluster := range clusters {
		clusterTables = append(clusterTables, cluster.(*model.ClusterTable))
	}

	return clusterTables, nil
}

// DeleteCloudClusters - Delete clusters on cloud
func (db *DB) DeleteCloudClusters(cloudUid string) (int64, error) {
	result, err := db.GetClient().Exec(deleteCloudClusters, cloudUid)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}
