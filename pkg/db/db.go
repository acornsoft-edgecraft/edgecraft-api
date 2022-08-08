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
	GetAllCloud() ([]model.Cloud, error)
	GetCloud(uuid.UUID) (*model.Cloud, error)
	UpdateRegisterCloud(*model.Cloud) (int64, error)
	CreateRegisterCloud(*model.Cloud) error
	DeleteRegisterCloud(uuid.UUID) (int64, error)
}