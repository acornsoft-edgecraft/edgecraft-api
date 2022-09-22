/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package postgresdb

import "github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"

// InsertCluster - Insert a new Baremetal Cluster
func (db *DB) InsertCluster(cluster *model.ClusterTable) error {
	return db.GetClient().Insert(cluster)
}

// GetCloudCluster - Returns a CloudCluster
func (db *DB) GetCloudCluster(cloudUid string) (*model.ClusterTable, error) {
	obj, err := db.GetClient().Get(&model.ClusterTable{}, cloudUid)
	if err != nil {
		return nil, err
	}
	if obj != nil {
		res := obj.(*model.ClusterTable)
		return res, nil
	}
	return nil, nil
}
