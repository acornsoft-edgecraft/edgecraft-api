package postgresdb

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
	"github.com/gofrs/uuid"
)

const getAllCloudNodeSQL = `
SELECT *
FROM tbl_cloud_node c 
`
const selectCloudNodeSQL = `
SELECT *
FROM tbl_cloud_node c
WHERE
cloud_uid = $1
and cloud_cluster_uid = $2
`

// CreateCloudNode - Registration a new Cloud Nodes
func (db *DB) CreateCloudNode(create *model.CloudNode) error {
	return db.GetClient().Insert(create)
}

// SelectCloudNode - Returns a matching value for cloud clusters
func (db *DB) SelectCloudNode(uid uuid.UUID, clusterUid uuid.UUID) ([]model.CloudNode, error) {
	var res []model.CloudNode
	_, err := db.GetClient().Select(&res, selectCloudNodeSQL, uid, clusterUid)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SelectMasterCloudNode - Returns a SelectMasterCloudNode
func (db *DB) SelectMasterCloudNode(uid uuid.UUID, clusterUid uuid.UUID) ([]model.CloudNode, error) {
	var res []model.CloudNode
	_, err := db.GetClient().Select(&res, selectCloudNodeSQL, uid, clusterUid)
	if err != nil {
		return nil, err
	}
	return res, nil

	// obj, err := db.GetClient().Get(&model.MasterNode{}, uid, clusterUid)
	// if err != nil {
	// 	return nil, err
	// }
	// if obj != nil {
	// 	res := obj.([]model.MasterNode)
	// 	return res, nil
	// }
	// return nil, nil
}

// SelectWorkerCloudNode - Returns a SelectWorkerCloudNode
func (db *DB) SelectWorkerCloudNode(uid uuid.UUID, clusterUid uuid.UUID) ([]model.CloudNode, error) {
	obj, err := db.GetClient().Get(&model.CloudNode{}, uid, clusterUid)
	if err != nil {
		return nil, err
	}
	if obj != nil {
		res := obj.([]model.CloudNode)
		return res, nil
	}
	return nil, nil
}

// GetCloudNode - Returns a GetCloudNode
func (db *DB) GetCloudNode(uid uuid.UUID, clusterUid uuid.UUID) (*model.CloudNode, error) {
	obj, err := db.GetClient().Get(&model.CloudNode{}, uid, clusterUid)
	if err != nil {
		return nil, err
	}
	if obj != nil {
		res := obj.(*model.CloudNode)
		return res, nil
	}
	return nil, nil
}

// GetAllCloudNode - Returns all Cloud list
func (db *DB) GetAllCloudNode() ([]model.CloudNode, error) {
	var res []model.CloudNode
	_, err := db.GetClient().Select(&res, getAllCloudNodeSQL)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateCloud - saves the given RegisterCloud struct
func (db *DB) UpdateCloudNode(req *model.CloudNode) (int64, error) {
	// Find and Update
	utils.Print(req)
	count, err := db.GetClient().Update(req)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// UpdateCloudNodes - saves the given RegisterCloud struct
func (db *DB) UpdateCloudNodes(req []*model.CloudNode) (int64, error) {
	// Find and Update
	var updateNodes []interface{}
	for _, data := range req {
		updateNodes = append(updateNodes, data)
	}
	utils.Print(req)
	count, err := db.GetClient().Update(updateNodes...)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// DeleteCloud - deletes the RegisterCloud with the given id
func (db *DB) DeleteCloudNode(uid uuid.UUID) (int64, error) {
	count, err := db.GetClient().Delete(&model.CloudNode{CloudNodeUid: &uid})
	if err != nil {
		return -1, err
	}
	return count, nil
}

// DeleteCloud - deletes the RegisterCloud with the given id
func (db *DB) DeleteAllCloudNode(uid uuid.UUID) (int64, error) {
	count, err := db.GetClient().Delete(&model.DelCloudNode{CloudUid: &uid})
	if err != nil {
		return -1, err
	}
	return count, nil
}
