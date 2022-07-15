package services

import (
	"errors"

	"github.com/acornsoft-edgecraft/edgecraft-api/internal/database"
	"github.com/acornsoft-edgecraft/edgecraft-api/internal/entities"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
	"github.com/google/uuid"
)

func NewTestAPI(testAPI *entities.TestAPI) (*entities.TestAPI, error) {
	var err error
	tx := database.Conn.Create(&testAPI)
	if tx.Error != nil {
		err = tx.Error
		utils.WarningLog.Println(err.Error())
	}
	return testAPI, err
}

func GetTestAPI(testAPIId uuid.UUID) (*entities.TestAPI, error) {
	var err error
	var testAPI *entities.TestAPI
	tx := database.Conn.Find(&testAPI, testAPIId)
	if tx.Error != nil {
		err = tx.Error
		utils.WarningLog.Println(err.Error())
	} else if testAPI.ID == nil {
		err = errors.New("product not found")
	}
	return testAPI, err
}

func GetTestJsonbs(testAPIId uuid.UUID) (interface{}, error) {
	var err error
	var testAPI *entities.TestAPI
	tx := database.Conn.Find(&testAPI, testAPIId)
	if tx.Error != nil {
		err = tx.Error
		utils.WarningLog.Println(err.Error())
	} else if testAPI.ID == nil {
		err = errors.New("jsonb not found")
	}
	return testAPI.API_response.Data, err
}

func GetTestJsonb(testAPIId uuid.UUID, parm string) (interface{}, error) {
	type Ttt struct {
		Data interface{} `json:"data"`
	}

	var err error
	var testAPI *entities.TestAPI
	var ttt Ttt
	tx := database.Conn.Table("api_response").Find(&ttt)
	// tx := database.Conn.Select(datatypes.JSONQuery("api_response").HasKey("data").HasKey(parm)).Find(&testAPI, testAPIId)
	if tx.Error != nil {
		err = tx.Error
		utils.WarningLog.Println(err.Error())
	} else if testAPI.ID == nil {
		err = errors.New("jsonb not found")
	}
	return ttt, err
}

func GetTestAPIs() (*[]entities.TestAPI, error) {
	var err error
	var testAPIs *[]entities.TestAPI
	tx := database.Conn.Order("api_url").Find(&testAPIs)
	if tx.Error != nil {
		err = tx.Error
		utils.WarningLog.Println(err.Error())
	} else if *testAPIs == nil {
		err = errors.New("testAPIs not found")
	}
	return testAPIs, err
}

func UpdateTestAPI(newTestAPI *entities.TestAPI, testAPIId uuid.UUID) (*entities.TestAPI, error) {
	testAPI, err := GetTestAPI(testAPIId)
	if err == nil {
		tx := database.Conn.Model(testAPI).Updates(newTestAPI)
		if tx.Error != nil {
			err = tx.Error
			utils.WarningLog.Println(err.Error())
		}
	}
	return testAPI, err
}

func DeleteTestAPI(testAPIId uuid.UUID) (*entities.TestAPI, error) {
	testAPI, err := GetTestAPI(testAPIId)
	if err == nil {
		tx := database.Conn.Delete(testAPI)
		if tx.Error != nil {
			err = tx.Error
			utils.WarningLog.Println(err.Error())
		}
	}
	return testAPI, err
}
