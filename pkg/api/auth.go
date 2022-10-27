package api

import (
	"strings"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/response"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
	"github.com/labstack/echo/v4"
)

// LoginHandler - 사용자 로그인 처리
// @Tags        Auth
// @Summary     Login
// @Description User Login
// @ID          Login
// @Produce     json
// @Param       loginInfo body     model.Login true "Request json"
// @Success     200       {object} response.ReturnData
// @Router      /auth [post]
func (a *API) LoginHandler(c echo.Context) error {
	var login model.Login

	// 요청 데이터 검증
	err := getRequestData(c.Request(), &login)
	if err != nil {
		return response.Errorf(c, common.CodeInvalidUser, err)
	} else if code, ok := login.Validate(); !ok {
		return response.Errorf(c, code, err)
	}

	//logger.Infof("Receved login data : %s", utils.GetGoString(login))

	// Email은 항상 소문자로 처리
	*login.Email = strings.ToLower(*login.Email)

	// 사용자 정보 조회
	user, err := a.Db.GetUserByEmail(*login.Email)

	if user == nil || err != nil {
		return response.Errorf(c, common.CodeInvalidUser, err)
	}
	// TODO: Check status codes
	if user.Status == 2 {
		return response.Errorf(c, common.CodeFaildStatusUser, nil)
	}

	// 비밀번호 검증
	if code, ok := user.MatchPassword(*login.Password); !ok {
		return response.Errorf(c, code, nil)
	}
	// TODO: 사용자 등록 및 패스워드 변경 상황인지 검증 필요
	// TODO: 사용자 메뉴 설정
	// TODO: 공통코드 설정

	// TODO: JWT 토큰 설정
	//user.token

	// TODO: Login History 설정

	// TODO: 최종 결과 정보 설정
	var tr model.TokenResponse
	var loginUser model.LoginInfo
	tr.AccessToken = nil
	tr.RefreshToken = nil

	if err := utils.CopyTo(&loginUser, &user); err != nil {
		return response.Errorf(c, common.CodeProcessingError, err)
	}

	tr.User = &loginUser

	return response.Write(c, c.Request(), tr)
}
