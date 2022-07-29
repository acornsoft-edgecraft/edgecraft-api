package api

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/response"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/labstack/echo/v4"
)

func (a *API) AllCloudListHandler(c echo.Context) error {
	res, err := a.Db.GetAllCloud()
	if err != nil {
		return response.Errorf(c, common.CodeFailedDatabase, err)
	}
	return response.Write(c, nil, res)
}
