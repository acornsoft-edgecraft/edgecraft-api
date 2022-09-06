package postgresdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"testing"
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/db"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	"github.com/gofrs/uuid"
)

func getDbObject() (db.DB, error) {
	//TODO config 설정\
	// dbConfig := &Config{
	// 	Type:         "postgres",
	// 	DatabaseName: "edgecraft",
	// 	SchemaName:   "edgecraft",
	// 	Host:         "192.168.77.42",
	// 	Port:         "31000",
	// 	UserName:     "edgecraft",
	// 	Password:     "edgecraft",
	// 	MaxIdleConns: 5,
	// 	MaxOpenConns: 100,
	// }
	// - localhost
	dbConfig := &Config{
		Type:         "postgres",
		DatabaseName: "edgecraft",
		SchemaName:   "edgecraft",
		Host:         "localhost",
		Port:         "5432",
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

	// 출력
	var buffer bytes.Buffer
	err = PrettyEncode(getClouds, &buffer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(buffer.String())
}

func Test_GetAllCodeGroup(t *testing.T) {
	db, _ := getDbObject()
	defer db.CloseConnection()
	getClouds, err := db.GetAllCodeGroup()
	if err != nil {
		fmt.Printf("error : %s", err)
	}

	// 출력
	var buffer bytes.Buffer
	err = PrettyEncode(getClouds, &buffer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(buffer.String())
}

func Test_CreateCodeGroup(t *testing.T) {

	item1 := "CLOUD_STATUS"
	item2 := "클라우드 상태 코드"
	item3 := true
	item11 := "CLOUD_TYPE"
	item12 := "클라우드 유형"
	item13 := true
	item21 := "K8S_VERSION"
	item22 := "Kubernetes Versio n"
	item23 := true
	item31 := "IMAGE_CHECKSUM_TYPE"
	item32 := "Image Checksum Type"
	item33 := true
	item41 := "IMAGE_FORMAT"
	item42 := "Image Format"
	item43 := true
	item51 := "BOOT_MODE"
	item52 := "Boot Mode"
	item53 := true
	item61 := "NODE_TYPE"
	item62 := "Node TypeNode Type"
	item63 := true

	data := model.CodeGroup{
		CodeGroupName:        &item1,
		CodeGroupDescription: &item2,
		UseYn:                &item3,
	}
	data1 := model.CodeGroup{
		CodeGroupName:        &item11,
		CodeGroupDescription: &item12,
		UseYn:                &item13,
	}
	data2 := model.CodeGroup{
		CodeGroupName:        &item21,
		CodeGroupDescription: &item22,
		UseYn:                &item23,
	}
	data3 := model.CodeGroup{
		CodeGroupName:        &item31,
		CodeGroupDescription: &item32,
		UseYn:                &item33,
	}
	data4 := model.CodeGroup{
		CodeGroupName:        &item41,
		CodeGroupDescription: &item42,
		UseYn:                &item43,
	}
	data5 := model.CodeGroup{
		CodeGroupName:        &item51,
		CodeGroupDescription: &item52,
		UseYn:                &item53,
	}
	data6 := model.CodeGroup{
		CodeGroupName:        &item61,
		CodeGroupDescription: &item62,
		UseYn:                &item63,
	}
	arrayData := []model.CodeGroup{
		data1,
		data2,
		data3,
		data4,
		data5,
		data6,
	}

	db, _ := getDbObject()
	defer db.CloseConnection()
	for _, i := range arrayData {
		err := db.CreateCodeGroup(&i)
		if err != nil {
			fmt.Printf("error : %s", err)
		}
	}
	err := db.CreateCodeGroup(&data)
	if err != nil {
		fmt.Printf("error : %s", err)
	}

	// 출력
	var buffer bytes.Buffer
	err = PrettyEncode(arrayData, &buffer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(buffer.String())
}

func Test_SearchCodeGroup(t *testing.T) {

	aa := "CLOUD_STATUS"
	bb := true
	cc := "클라우드 상태 코드"
	searchCodeGroup := model.CodeGroup{
		CodeGroupName:        &aa,
		UseYn:                &bb,
		CodeGroupDescription: &cc,
	}

	db, _ := getDbObject()
	defer db.CloseConnection()
	getClouds, err := db.SearchCodeGroup(searchCodeGroup)
	if err != nil {
		fmt.Printf("error : %s", err)
	}

	// 출력
	var buffer bytes.Buffer
	err = PrettyEncode(getClouds, &buffer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(buffer.String())
}

func Test_Transaction(t *testing.T) {
	db, _ := getDbObject()
	defer db.CloseConnection()

	// 모델 생성
	req := getReqCloud()

	// Transaction 시작
	txdb, err := db.BeginTransaction()
	if err != nil {
		t.Error(err)
	}

	// 여러개의 query 실행
	err = txdb.CreateCloud(req)
	if err != nil {
		txErr := txdb.Rollback()
		if txErr != nil {
			t.Error("DB Rollback Failed.", txErr)
		}
		t.Error(err)
	}

	// transaction 종료 - txdb.Rollback()
	txErr := txdb.Commit()
	if txErr != nil {
		t.Error("DB Commit Failed.", txErr)
	}

	// 출력
	var buffer bytes.Buffer
	err = PrettyEncode(req, &buffer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(buffer.String())
}

func Test_DeleteCloud(t *testing.T) {
	db, _ := getDbObject()
	defer db.CloseConnection()

	// 모델 생성
	req := uuid.FromStringOrNil("2e622c54-556d-4e7c-8b4f-1619dd01d781")

	count, err := db.DeleteCloud(req)
	if err != nil {
		fmt.Printf("error : %s", err)
	}

	// 출력
	var buffer bytes.Buffer
	err = PrettyEncode(count, &buffer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(buffer.String())
}

func Parse(s string) {
	panic("unimplemented")
}

// - Test Model value
func getReqCloud() *model.Cloud {
	// cloudUID := uuid.New()
	cloudName := "test-6"
	cloudType := "baremetal"
	cloudDesc := "test-description"
	cloudStatus := ""
	creator := "user-1"
	createdAt := time.Now()
	req := &model.Cloud{
		// CloudUID:    &cloudUID,
		CloudName:   &cloudName,
		CloudType:   &cloudType,
		CloudDesc:   &cloudDesc,
		CloudStatus: &cloudStatus,
		Creator:     &creator,
		CreatedAt:   &createdAt,
	}
	// req := &model.Cloud{
	// 	CloudUID:  uuid.New(),
	// 	CloudName: "test-6",
	// 	CloudType: "baremetal",
	// 	CloudDesc: "test-description",
	// 	// CloudStatus: null,
	// 	Creator:   "user-1",
	// 	CreatedAt: time.Now(),
	// 	// TestInt:   ,
	// }
	// req := &model.Cloud{
	// 	CloudUID:    cloudUID,
	// 	CloudName:   utils.NullString{String: cloudName, Valid: true},
	// 	CloudType:   utils.NullString{String: cloudType, Valid: true},
	// 	CloudDesc:   utils.NullString{String: cloudDesc, Valid: true},
	// 	CloudStatus: utils.NullString{String: cloudStatus, Valid: true},
	// 	Creator:     utils.NullString{String: creator, Valid: true},
	// 	CreatedAt:   utils.NullTime{Time: createdAt, Valid: true},
	// 	Updater:     utils.NullString{},
	// 	UpdatedAt:   utils.NullTime{Time: time.Time{}, Valid: true},
	// }
	return req
}

// -- public funtionable
func PrettyEncode(data interface{}, out io.Writer) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "    ")
	if err := enc.Encode(data); err != nil {
		return err
	}
	return nil
}
