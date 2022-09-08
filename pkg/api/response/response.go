// Package response -
package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/labstack/echo/v4"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ReturnData - Represents the structure of an response data or error information
type ReturnData struct {
	Error   bool        `json:"isError"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// getFields - Returns logging fields from request
func getFields(req *http.Request) logger.Fields {
	return logger.Fields{
		"host":       req.Host,
		"address":    req.RemoteAddr,
		"method":     req.Method,
		"requestURI": req.RequestURI,
		"proto":      req.Proto,
		"userAgent":  req.UserAgent(),
	}
}

// ===== [ Public Functions ] =====

// Errorf - Returns an new error response
func Errorf(c echo.Context, code int, err error) error {
	msg := common.GetMessageByCode(code)
	logger.WithFields(getFields(c.Request())).WithError(err).Debug("Processed Code: " + strconv.Itoa(code) + ", Message: " + msg)

	// Message Applying
	if err != nil {
		msg = fmt.Sprintf("%v\n (%v)", msg, err.Error())
	}

	returnData := ReturnData{
		Error:   true,
		Code:    code,
		Message: msg,
		Data:    nil,
	}

	return c.JSON(http.StatusOK, &returnData)
}

// Errorf - Returns an new error response
func ErrorfReqRes(c echo.Context, reqBody interface{}, code int, err error) error {
	msg := common.GetMessageByCode(code)
	logger.WithFields(getFields(c.Request())).WithError(err).Debug("Processed Code: " + strconv.Itoa(code) + ", Message: " + msg)

	// Message Applying
	if err != nil {
		msg = fmt.Sprintf("%v\n (%v)", msg, err.Error())
	}

	returnData := ReturnData{
		Error:   true,
		Code:    code,
		Message: msg,
		Data:    nil,
	}

	req := c.Request()
	reqJson, _ := json.Marshal(reqBody)
	retJson, _ := json.Marshal(returnData)
	formatter := fmt.Sprintf("%v %v %v %v %v %v ", code, req.RemoteAddr, req.Method, req.RequestURI, string(reqJson), string(retJson))
	logger.Error(formatter)

	return c.JSON(http.StatusOK, &returnData)
}

// Write - Writes an new json response
func WriteWithFields(c echo.Context, data interface{}) error {
	logger.WithFields(getFields(c.Request())).Debug(data)

	returnData := ReturnData{
		Error:   false,
		Code:    common.CodeOK,
		Message: common.GetMessageByCode(common.CodeOK),
		Data:    data,
	}

	return c.JSON(http.StatusOK, &returnData)
}

// Write - Writes an new json response
func Write(c echo.Context, reqBody interface{}, data interface{}) error {
	req := c.Request()
	returnData := ReturnData{
		Error:   false,
		Code:    common.CodeOK,
		Message: common.GetMessageByCode(common.CodeOK),
		Data:    data,
	}
	reqJson, _ := json.Marshal(reqBody)
	retJson, _ := json.Marshal(returnData)

	formatter := fmt.Sprintf("%v %v %v %v %v %v ", common.CodeOK, req.RemoteAddr, req.Method, req.RequestURI, string(reqJson), string(retJson))
	logger.Info(formatter)
	return c.JSON(http.StatusOK, &returnData)
}
