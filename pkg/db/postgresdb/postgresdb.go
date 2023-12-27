package postgresdb

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/db"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	_ "github.com/lib/pq" // db driver for postgreSQL
	"gopkg.in/gorp.v2"
)

// ===== [ Constants and Variables ] =====

var (
	// verify interface compliance
	_ = []gorp.Dialect{
		gorp.PostgresDialect{},
	}
)

// ===== [ Types ] =====

// Config - Represents the postgreSQL configuration
type Config struct {
	Type         string `yaml:"type"`
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	DatabaseName string `yaml:"database_name"`
	SchemaName   string `yaml:"schema_name"`
	UserName     string `yaml:"username"`
	Password     string `yaml:"password"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

// DB - Represents the structure of the database
type DB struct {
	config *Config
	client *gorp.DbMap
	tx     *gorp.Transaction // transaction 사용시에만 셋팅되는 변수
}

// ===== [ Implementations ] =====

// CloseConnection - Closes the database connection
func (db *DB) CloseConnection() error {
	db.client.TraceOff()

	return db.client.Db.Close()
}

// BeginTransaction transaction 처리 시작 하면서 사용할 새로운 DB 객체를 transaction 포함하여 만든다.
func (db *DB) BeginTransaction() (db.DB, error) {
	transaction, err := db.client.Begin()
	if err != nil {
		return nil, err
	}
	return &DB{
		config: db.config,
		client: db.client,
		tx:     transaction,
	}, nil
}

// Commit transaction commit
func (db *DB) Commit() error {
	if db.tx == nil {
		return fmt.Errorf("not exist a transaction")
	}
	err := db.tx.Commit()
	if err == nil {
		db.tx = nil
	}
	return err
}

// Rollback transaction rollback
func (db *DB) Rollback() error {
	if db.tx == nil {
		return fmt.Errorf("not exist a transaction")
	}
	err := db.tx.Rollback()
	if err == nil {
		db.tx = nil
	}
	return err
}

// TransactionScope - 데이터베이스 처리를 위한 트랜잭션 Scope 실행
func (db *DB) TransactionScope(stmts func(txDB db.DB) error) error {
	rollbacked := false

	// transaction start
	tx, err := db.client.Begin()
	if err != nil {
		return err
	}

	// 종료할 때 Commit 처리
	defer func() {
		// rollback이 아닌 경우만 Commit 처리
		if tx != nil && !rollbacked {
			err := tx.Commit()
			if err != nil {
				logger.WithError(err).Info("Commit failed with error.")
			}
			logger.Info("Database changes commited.")
		}
	}()

	txDB := &DB{
		config: db.config,
		client: db.client,
		tx:     tx,
	}

	// 블럭 명령 실행
	err = stmts(txDB)
	if err != nil {
		logger.WithError(err).Info("Database changes will be rollback")

		txErr := txDB.Rollback()
		if txErr != nil {
			logger.WithError(txErr).Info("Rollback failed with error.")
		}
		rollbacked = true

		return err
	}

	return nil
}

// GetClient DB 처리에 사용할 client를 반환한다. transaction 사용할때는 transaction 객체, 사용하지 않을때는 dbMap 객체를 리턴한다.
func (db *DB) GetClient() db.SQLExecutor {
	if db.tx != nil {
		return db.tx
	}
	return db.client
}

// ===== [ Private Functions ] =====

// 테이블과 Struct 맵핑(abc순으로 정리)
func newDbMap(conf *Config) *gorp.DbMap {
	dbmap := &gorp.DbMap{Db: connect(conf), Dialect: gorp.PostgresDialect{}}
	dbmap.Db.SetMaxIdleConns(conf.MaxIdleConns)
	dbmap.Db.SetMaxOpenConns(conf.MaxOpenConns)

	// For Code (manual key)
	dbmap.AddTableWithName(model.CodeGroupTable{}, "tbl_code_group").SetKeys(false, "group_id")
	dbmap.AddTableWithName(model.CodeTable{}, "tbl_code").SetKeys(false, "group_id", "code")

	// For Cloud (manual key)
	dbmap.AddTableWithName(model.CloudTable{}, "tbl_cloud").SetKeys(false, "cloud_uid")
	dbmap.AddTableWithName(model.ClusterTable{}, "tbl_cloud_cluster").SetKeys(false, "cloud_uid", "cluster_uid")
	dbmap.AddTableWithName(model.NodeTable{}, "tbl_cloud_node").SetKeys(false, "cloud_uid", "cluster_uid", "node_uid")

	// For Openstack (manual key)
	dbmap.AddTableWithName(model.OpenstackClusterTable{}, "tbl_cluster").SetKeys(false, "cloud_uid", "cluster_uid")
	dbmap.AddTableWithName(model.NodeSetTable{}, "tbl_nodeset").SetKeys(false, "cluster_uid", "nodeset_uid")
	dbmap.AddTableWithName(model.OpenstackBenchmarksTable{}, "tbl_cluster_benchmarks").SetKeys(false, "cluster_uid", "benchmarks_uid")
	dbmap.AddTableWithName(model.BackResTable{}, "tbl_cluster_backres").SetKeys(false, "cloud_uid", "cluster_uid", "backres_uid")

	dbmap.AddTableWithName(model.UserTable{}, "tbl_user").SetKeys(false, "user_uid")
	// // SetKeys(isAutoIncr bool, fieldNames ...string)
	// // SetKeys(true) means we have a auto increment primary key, which
	// // will get automatically bound to your struct post-insert

	// dbmap.AddTableWithName(model.Cloud{}, "tbl_cloud").SetKeys(true, "cloud_uid")
	// dbmap.AddTableWithName(model.CloudCluster{}, "tbl_cloud_cluster").SetKeys(true, "cloud_cluster_uid")
	// dbmap.AddTableWithName(model.CloudCluster{}, "tbl_cloud_cluster").SetKeys(false, "cloud_uid")
	// dbmap.AddTableWithName(model.K8s{}, "tbl_cloud_cluster").SetKeys(false, "cloud_uid")
	// dbmap.AddTableWithName(model.Baremetal{}, "tbl_cloud_cluster").SetKeys(false, "cloud_uid")
	// dbmap.AddTableWithName(model.Etcd{}, "tbl_cloud_cluster").SetKeys(false, "cloud_uid")
	// dbmap.AddTableWithName(model.ClusterNodes{}, "tbl_cloud_cluster").SetKeys(false, "cloud_uid")
	// dbmap.AddTableWithName(model.CloudNode{}, "tbl_cloud_node").SetKeys(true, "cloud_node_uid")
	// dbmap.AddTableWithName(model.DelCloudNode{}, "tbl_cloud_node").SetKeys(false, "cloud_uid")
	// dbmap.AddTableWithName(model.CodeGroup{}, "tbl_code_group").SetKeys(true, "code_group_uid")
	// dbmap.AddTableWithName(model.Code{}, "tbl_code").SetKeys(true, "code_uid")
	// // dbmap.AddTableWithName(model.MasterNode{}, "tbl_cloud_node").SetKeys(false, "cloud_uid", "cloud_cluster_uid")
	// // dbmap.AddTableWithName(model.WorkerNode{}, "tbl_cloud_node").SetKeys(false, "cloud_uid", "cloud_cluster_uid")

	return dbmap
}

// connect db connection
func connect(conf *Config) *sql.DB {
	// dsn := fmt.Sprintf("postgres://%v:%v@%s:%s/%v?sslmode=disable", conf.UserName, conf.Password, conf.Host, conf.Port, conf.DatabaseName)
	dsn := fmt.Sprintf("user=%s password='%s' host=%s port=%s dbname=%s sslmode=disable search_path=%s", conf.UserName, conf.Password, conf.Host, conf.Port, conf.DatabaseName, conf.SchemaName)
	db, err := sql.Open(conf.Type, dsn)
	if err != nil {
		panic("Error connecting to db: " + err.Error())
	}

	return db
}

// ===== [ Public Functions ] =====

// NewConnection - Creates a new database connection
func NewConnection(conf *Config) (db.DB, error) {

	// create a PostgreSQL DbMap
	client := newDbMap(conf)

	// Loggin Enabling
	client.TraceOn("[GORP]", log.New(os.Stdout, "PostgreSQL", log.Lmicroseconds))

	// Checking the connection
	err := client.Db.PingContext(context.Background())
	if err != nil {
		return nil, err
	}

	return &DB{
		config: conf,
		client: client,
	}, nil
}
