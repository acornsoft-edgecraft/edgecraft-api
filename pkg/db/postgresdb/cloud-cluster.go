package postgresdb

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	"github.com/gofrs/uuid"
)

const getAllCloudClusterSQL = `
SELECT *
FROM tbl_cloud_cluster c 
`

const selectCloudClusterSQL = `
SELECT *
FROM tbl_cloud_cluster c
WHERE
cloud_uid = $1
`

// RegisterCloudCluster - Registration a new CloudCluster
func (db *DB) CreateCloudCluster(create *model.CloudCluster) error {
	return db.GetClient().Insert(create)
}

// SelectCloudCluster - Returns a matching value for cloud clusters
func (db *DB) SelectCloudCluster(uid uuid.UUID) (*model.CloudCluster, error) {
	var res *model.CloudCluster
	_, err := db.GetClient().Select(&res, selectCloudClusterSQL, uid)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SelectK8sCloudCluster - Returns a matching value for cloud clusters
func (db *DB) SelectK8sCloudCluster(uid uuid.UUID) (*model.K8s, error) {
	// var res *model.K8s
	obj, err := db.GetClient().Get(&model.K8s{}, uid)
	if err != nil {
		return nil, err
	}
	if obj != nil {
		res := obj.(*model.K8s)
		return res, nil
	}
	return nil, nil
}

// SelectBaremetalCloudCluster - Returns a matching value for cloud clusters
func (db *DB) SelectBaremetalCloudCluster(uid uuid.UUID) (*model.ClusterBaremetal, error) {
	obj, err := db.GetClient().Get(&model.ClusterBaremetal{}, uid)
	if err != nil {
		return nil, err
	}
	if obj != nil {
		res := obj.(*model.ClusterBaremetal)
		return res, nil
	}
	return nil, nil
}

// SelectNodeCloudCluster - Returns a matching value for cloud clusters
func (db *DB) SelectNodeCloudCluster(uid uuid.UUID) (*model.ClusterNodes, error) {
	obj, err := db.GetClient().Get(&model.ClusterNodes{}, uid)
	if err != nil {
		return nil, err
	}
	if obj != nil {
		res := obj.(*model.ClusterNodes)
		return res, nil
	}
	return nil, nil
}

// SelectNodeCloudCluster - Returns a matching value for cloud clusters
func (db *DB) SelectEtcdCloudCluster(uid uuid.UUID) (*model.Etcd, error) {
	obj, err := db.GetClient().Get(&model.Etcd{}, uid)
	if err != nil {
		return nil, err
	}
	if obj != nil {
		res := obj.(*model.Etcd)
		return res, nil
	}
	return nil, nil
}

// GetCloudCluster - Returns a CloudCluster
func (db *DB) GetCloudCluster(uid uuid.UUID) (*model.CloudCluster, error) {

	obj, err := db.GetClient().Get(&model.CloudCluster{}, uid)
	if err != nil {
		return nil, err
	}
	if obj != nil {
		res := obj.(*model.CloudCluster)
		return res, nil
	}
	return nil, nil
}

// GetAllCloud - Returns all Cloud list
func (db *DB) GetAllCloudCluster() ([]model.CloudCluster, error) {
	var res []model.CloudCluster
	_, err := db.GetClient().Select(&res, getAllCloudClusterSQL)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateCloud - saves the given RegisterCloud struct
func (db *DB) UpdateCloudCluster(req *model.CloudCluster) (int64, error) {
	// Find and Update
	count, err := db.GetClient().Update(req)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// DeleteCloud - deletes the RegisterCloud with the given id
func (db *DB) DeleteCloudCluster(uid uuid.UUID) (int64, error) {
	count, err := db.GetClient().Delete(&model.CloudCluster{CloudClusterUid: &uid})
	if err != nil {
		return -1, err
	}
	return count, nil
}
