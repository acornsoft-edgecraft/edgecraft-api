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
	GetClouds() ([]model.CloudList, error)
	GetCloud(string) (*model.CloudTable, error)
	InsertCloud(*model.CloudTable) error
	UpdateCloud(*model.CloudTable) (int64, error)
	DeleteCloud(string) (int64, error)

	// tbl_cloud_cluster
	GetClusters(string) ([]*model.ClusterTable, error)
	GetCluster(string, string) (*model.ClusterTable, error)
	InsertCluster(*model.ClusterTable) error
	UpdateCluster(*model.ClusterTable) (int64, error)
	DeleteCluster(string, string) (int64, error)
	DeleteClusters(string) (int64, error)

	// tbl_cloud_node
	GetNodes(string, string) ([]*model.NodeTable, error)
	GetNodesByCloud(string) ([]*model.NodeTable, error)
	GetNode(string, string, string) (*model.NodeTable, error)
	GetNodeByCloud(string, string) (*model.NodeTable, error)

	InsertNode(*model.NodeTable) error
	UpdateNode(*model.NodeTable) (int64, error)
	DeleteNode(string) (int64, error)
	DeleteNodes(string, string) (int64, error)
	DeleteCloudNodes(string) (int64, error)

	// tbl_cluster (Openstack)
	GetOpenstackClusters(string) ([]model.OpenstackClusterList, error)
	GetOpenstackCluster(string, string) (*model.OpenstackClusterTable, error)
	InsertOpenstackCluster(*model.OpenstackClusterTable) error
	UpdateOpenstackCluster(*model.OpenstackClusterTable) (int64, error)
	DeleteOpenstackCluster(string, string) (int64, error)

	// tbl_nodeset (Openstack)
	GetNodeSets(string) ([]*model.NodeSetTable, error)
	GetNodeSet(string, string) (*model.NodeSetTable, error)
	InsertNodeSet(*model.NodeSetTable) error
	UpdateNodeSet(*model.NodeSetTable) (int64, error)
	DeleteNodeSet(string, string) (int64, error)
	DeleteNodeSets(string) (int64, error)

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
