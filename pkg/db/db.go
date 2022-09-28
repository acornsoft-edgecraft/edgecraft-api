package db

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	"gopkg.in/gorp.v2"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// Config - Represents the configuration of the database interface

// SQLExecutor - gorp SqlExecutor 상속받은 Interface
type SQLExecutor interface {
	gorp.SqlExecutor
	UpdateColumns(filter gorp.ColumnFilter, list ...interface{}) (int64, error)
}

// DB - Interface with must be implemented by all db drivers
type DB interface {
	// Connection manage
	CloseConnection() error
	BeginTransaction() (DB, error)
	Commit() error
	Rollback() error
	GetClient() SQLExecutor

	// tbl_cloud
	GetCloudList() ([]model.CloudList, error)
	GetCloud(string) (*model.CloudTable, error)
	InsertCloud(*model.CloudTable) error
	UpdateCloud(*model.CloudTable) (int64, error)
	DeleteCloud(string) (int64, error)

	// tbl_cloud_cluster
	// GetAllCloudCluster() ([]model.CloudCluster, error)
	GetCluster(string, string) (*model.ClusterTable, error)
	SelectClusters(string) ([]*model.ClusterTable, error)
	InsertCluster(*model.ClusterTable) error
	UpdateCluster(*model.ClusterTable) (int64, error)
	DeleteCloudClusters(string) (int64, error)

	// SelectCloudCluster(uuid.UUID) (*model.CloudCluster, error)
	// SelectEtcdCloudCluster(uuid.UUID) (*model.Etcd, error)
	// SelectK8sCloudCluster(uuid.UUID) (*model.K8s, error)
	// SelectBaremetalCloudCluster(uuid.UUID) (*model.Baremetal, error)
	// SelectNodeCloudCluster(uuid.UUID) (*model.ClusterNodes, error)
	// UpdateCloudCluster(*model.CloudCluster) (int64, error)
	// CreateCloudCluster(*model.CloudCluster) error
	// DeleteCloudCluster(uuid.UUID) (int64, error)
	// DeleteAllCloudCluster(uuid.UUID) (int64, error)

	// tbl_cloud_node
	GetNode(string, string, string) (*model.NodeTable, error)
	SelectNodes(string, string) ([]*model.NodeTable, error)
	InsertNode(*model.NodeTable) error
	DeleteNodes(string, string) (int64, error)
	DeleteCloudNodes(string) (int64, error)

	// GetAllCloudNode() ([]model.CloudNode, error)
	// GetCloudNode(uuid.UUID, uuid.UUID) (*model.CloudNode, error)

	// SelectCloudNode(uuid.UUID, uuid.UUID) ([]model.CloudNode, error)
	// SelectMasterCloudNode(uuid.UUID, uuid.UUID) ([]model.CloudNode, error)
	// SelectWorkerCloudNode(uuid.UUID, uuid.UUID) ([]model.CloudNode, error)
	// UpdateCloudNode(*model.CloudNode) (int64, error)
	// UpdateCloudNodes([]*model.CloudNode) (int64, error)
	// CreateCloudNode(*model.CloudNode) error
	// DeleteCloudNode(uuid.UUID) (int64, error)
	// DeleteAllCloudNode(uuid.UUID) (int64, error)

	// tbl_code_group
	GetCodeGroupList() ([]*model.CodeGroupTable, error)
	GetCodeGroup(string) (*model.CodeGroupTable, error)
	InsertCodeGroup(*model.CodeGroupTable) error
	UpdateCodeGroup(*model.CodeGroupTable) (int64, error)
	DeleteCodeGroup(string) (int64, error)

	// tbl_code
	GetCodeList() ([]*model.CodeTable, error)
	GetCodeListByGroup(string) ([]*model.CodeTable, error)
	GetCode(string, int) (*model.CodeTable, error)
	InsertCode(*model.CodeTable) error
	UpdateCode(*model.CodeTable) (int64, error)
	DeleteCode(string, int) (int64, error)
	DeleteCodeByGroup(string) (int64, error)

	// tbl_user
	GetUserByEmail(email string) (*model.User, error)
	// GetUserById(id string) (*model.User, error)
}
