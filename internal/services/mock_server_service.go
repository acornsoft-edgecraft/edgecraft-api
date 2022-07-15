package services

import (
	"errors"

	"github.com/acornsoft-edgecraft/edgecraft-api/internal/database"
	"github.com/acornsoft-edgecraft/edgecraft-api/internal/entities"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
)

func GetMockRoute(url string, method string) (*entities.TestAPI, error) {
	var err error
	var mockRoute *entities.TestAPI
	tx := database.Conn.Where("api_url = ? AND api_method = ?", url[5:], method).Find(&mockRoute)
	if tx.Error != nil {
		err = tx.Error
		utils.WarningLog.Println(err.Error())
	} else if mockRoute.API_url == "" {
		err = errors.New("API Url not found")
	}

	return mockRoute, err
}
