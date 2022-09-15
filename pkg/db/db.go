package db

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	"github.com/gofrs/uuid"
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
	GetAllCloud() ([]model.ResCloud, error)
	GetSearchCloud(u model.Cloud) ([]model.Cloud, error)
	GetCloud(uuid.UUID) (*model.Cloud, error)
	UpdateCloud(*model.Cloud) (int64, error)
	CreateCloud(*model.Cloud) error
	DeleteCloud(uuid.UUID) (int64, error)

	// tbl_cloud_cluster
	GetAllCloudCluster() ([]model.CloudCluster, error)
	GetCloudCluster(uuid.UUID) (*model.CloudCluster, error)
	SelectCloudCluster(uuid.UUID) (*model.CloudCluster, error)
	SelectEtcdCloudCluster(uuid.UUID) (*model.Etcd, error)
	SelectK8sCloudCluster(uuid.UUID) (*model.K8s, error)
	SelectBaremetalCloudCluster(uuid.UUID) (*model.ClusterBaremetal, error)
	SelectNodeCloudCluster(uuid.UUID) (*model.ClusterNodes, error)
	UpdateCloudCluster(*model.CloudCluster) (int64, error)
	CreateCloudCluster(*model.CloudCluster) error
	DeleteCloudCluster(uuid.UUID) (int64, error)
	DeleteAllCloudCluster(uuid.UUID) (int64, error)

	// tbl_cloud_node
	GetAllCloudNode() ([]model.CloudNode, error)
	GetCloudNode(uuid.UUID, uuid.UUID) (*model.CloudNode, error)
	SelectCloudNode(uuid.UUID, uuid.UUID) ([]model.CloudNode, error)
	SelectMasterCloudNode(uuid.UUID, uuid.UUID) ([]model.CloudNode, error)
	SelectWorkerCloudNode(uuid.UUID, uuid.UUID) ([]model.CloudNode, error)
	UpdateCloudNode(*model.CloudNode) (int64, error)
	UpdateCloudNodes([]*model.CloudNode) (int64, error)
	CreateCloudNode(*model.CloudNode) error
	DeleteCloudNode(uuid.UUID) (int64, error)
	DeleteAllCloudNode(uuid.UUID) (int64, error)

	// tbl_code_group
	CreateCodeGroup(*model.CodeGroup) error
	GetAllCodeGroup() ([]model.CodeGroup, error)
	GetCodeGroup(uuid.UUID) (*model.CodeGroup, error)
	SelectCodeGroup(uuid.UUID) ([]model.CodeGroup, error)
	SearchCodeGroup(model.CodeGroup) ([]model.CodeGroup, error)

	// tbl_code
	GetAllCode() ([]model.Code, error)
	GetCode(uuid.UUID) (*model.Code, error)
	SelectCode(uuid.UUID) ([]model.Code, error)
	SearchCode(model.Code) ([]model.Code, error)

	// tbl_user
	GetUserByEmail(email string) (*model.User, error)
	// GetUserById(id string) (*model.User, error)
}
