/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package api

import (
	"errors"
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/response"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/db"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	"github.com/labstack/echo/v4"
)

// GetUserListHandler - 사용자 목록 조회
// @Tags        User
// @Summary     GetUserList
// @Description 전체 사용자 리스트
// @ID          GetUserList
// @Produce     json
// @Success     200     {object} response.ReturnData
// @Router      /users [get]
func (a *API) GetUserListHandler(c echo.Context) error {
	result, err := a.Db.GetUserList()
	if err != nil {
		return response.Errorf(c, common.CodeFailedDatabase, err)
	}
	return response.Write(c, nil, result)
}

// GetUserHandler - 사용자 상세 조회
// @Tags        User
// @Summary     GetUser
// @Description 사용자 상세 조회
// @ID          GetUser
// @Produce     json
// @Param       userId   path     string true "User UID"
// @Success     200     {object} response.ReturnData
// @Router      /users/{userId} [get]
func (a *API) GetUserHandler(c echo.Context) error {
	userId := c.Param("userId")
	if userId == "" {
		return response.ErrorfReqRes(c, userId, common.CodeInvalidParm, nil)
	}

	data, err := a.Db.GetUser(userId)
	data.Password = nil
	if err != nil {
		return response.Errorf(c, common.CodeFailedDatabase, err)
	}
	if data == nil {
		return response.Errorf(c, common.UserNotFound, err)
	}

	return response.Write(c, nil, data)
}

// SetUserHandler - 사용자 등록
// @Tags        User
// @Summary     SetUser
// @Description 사용자 등록
// @ID          SetUser
// @Produce     json
// @Param       user body     model.User true "User"
// @Success     200  {object} response.ReturnData
// @Router      /users [post]
func (a *API) SetUserHandler(c echo.Context) error {
	var user *model.User
	err := getRequestData(c.Request(), &user)
	if err != nil {
		return response.ErrorfReqRes(c, user, common.CodeInvalidData, err)
	}

	// TODO 중복 검증?

	userTable := &model.UserTable{}
	user.ToTable(userTable, false, "system", time.Now().UTC())

	err = a.Db.TransactionScope(func(txDB db.DB) error {
		err = txDB.InsertUser(userTable)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return response.ErrorfReqRes(c, userTable, common.CodeFailedDatabase, err)
	}

	return response.Write(c, user, nil)
}

// UpdateUserHandler - 사용자 정보 수정
// @Tags        User
// @Summary     UpdateUser
// @Description 사용자 수정
// @ID          UpdateUser
// @Produce     json
// @Param       userId 	path   string 		true "User ID"
// @Param       user 	body   model.User 	true "User"
// @Success     200  {object} response.ReturnData
// @Router      /users/{userId} [put]
func (a *API) UpdateUserHandler(c echo.Context) error {
	userId := c.Param("userId")
	if userId == "" {
		return response.ErrorfReqRes(c, userId, common.CodeInvalidParm, nil)
	}

	var user model.User
	err := getRequestData(c.Request(), &user)
	if err != nil {
		return response.ErrorfReqRes(c, user, common.CodeInvalidData, err)
	}

	// 사용자 정보 조회
	userTable, err := a.Db.GetUser(userId)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if userTable == nil {
		return response.ErrorfReqRes(c, nil, common.UserNotFound, err)
	}

	user.ToTable(userTable, true, "system", time.Now().UTC())

	err = a.Db.TransactionScope(func(txDB db.DB) error {
		cnt, err := txDB.UpdateUser(userTable)
		if err != nil {
			return err
		}
		if cnt == 0 {
			return errors.New("cannot find user for updating")
		}
		return nil
	})
	if err != nil {
		return response.ErrorfReqRes(c, userTable, common.CodeFailedDatabase, err)
	}

	return response.Write(c, user, nil)
}

// DeleteUserHandler - 사용자 정보 삭제
// @Tags        User
// @Summary     DeleteUser
// @Description 사용자 삭제
// @ID          DeleteUser
// @Produce     json
// @Param       userId 	path   string 		true "User ID"
// @Success     200  {object} response.ReturnData
// @Router      /users/{userId} [delete]
func (a *API) DeleteUserHandler(c echo.Context) error {
	userId := c.Param("userId")
	if userId == "" {
		return response.ErrorfReqRes(c, userId, common.CodeInvalidParm, nil)
	}

	// 사용자 정보 조회
	userTable, err := a.Db.GetUser(userId)
	if err != nil {
		return response.ErrorfReqRes(c, nil, common.CodeFailedDatabase, err)
	} else if userTable == nil {
		return response.ErrorfReqRes(c, nil, common.UserNotFound, err)
	}

	err = a.Db.TransactionScope(func(txDB db.DB) error {
		cnt, err := txDB.DeleteUser(userId)
		if err != nil {
			return err
		}
		if cnt == 0 {
			return errors.New("cannot find user for deleting")
		}
		return nil
	})
	if err != nil {
		return response.ErrorfReqRes(c, userId, common.CodeFailedDatabase, err)
	}

	return response.Write(c, userId, nil)
}
