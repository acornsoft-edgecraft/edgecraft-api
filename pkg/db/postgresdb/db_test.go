package postgresdb

import (
	"fmt"
	"testing"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/db"
)

func getDbObject() (db.DB, error) {
	//TODO config 설정
	// type: "postgres"
	// host: "edge-dev.acornsoft.io"
	// port: "31219"
	// database_name: "liteedge"
	// username: "liteedge"
	// password: "liteedge2020"
	// max_idle_conns: 5
	// max_open_conns: 100
	dbConfig := &Config{
		Type:         "postgres",
		DatabaseName: "edgecraft",
		SchemaName:   "edgecraft",
		Host:         "192.168.77.42",
		Port:         "31000",
		UserName:     "edgecraft",
		Password:     "edgecraft",
		MaxIdleConns: 5,
		MaxOpenConns: 100,
	}
	//TODO DB connection 생성
	db, err := NewConnection(dbConfig)
	if err != nil {
		return nil, err
	}
	return db, err
}

func Test_GetAllCloud(t *testing.T) {
	db, _ := getDbObject()
	defer db.CloseConnection()
	getClouds, err := db.GetAllCloud()
	if err != nil {
		fmt.Printf("error : %s", err)
	}
	fmt.Println(getClouds)
}
