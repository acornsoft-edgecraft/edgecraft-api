// Package middleware - Defines the handlers for middleware functionality
package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api/response"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"

	// lsession "github.com/acornsoft-edgecraft/edgecraft-api/pkg/session"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// GetID - Returns the id from given JWT claims
func GetID(claims jwt.MapClaims) (string, error) {
	id, ok := claims["id"].(string)
	if !ok {
		return "", fmt.Errorf("unauthorized")
	}
	return id, nil
}

// Preflight - Handles for browser preflight requests
func Preflight(c echo.Context) error {
	return c.String(http.StatusOK, "")
}

// HealthCheck Method
// @Tags        Common
// @Summary     Health check k8s-api-gateway
// @Description for health check
// @ID          health-check
// @Produce     json
// @Success     200 {object} response.ReturnData
// @Router      /health [get]
func HealthCheck(c echo.Context) error {
	return response.WriteWithFields(c, struct {
		Status string `json:"status"`
	}{"ok"})
}

// CustomLogger - Custom Logger middleware
func CustomLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()

			start := time.Now()
			// process next func
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()
			latencyHuman := stop.Sub(start).String()

			errMsg := ""
			if err != nil {
				// Error may contain invalid JSON e.g. `"`
				b, _ := json.Marshal(err.Error())
				b = b[1 : len(b)-1]
				errMsg = string(b)
			}

			byteIn := req.Header.Get(echo.HeaderContentLength)
			if byteIn == "" {
				byteIn = "0"
			}

			byteOut := strconv.FormatInt(res.Size, 10)

			//logger.WithFields(logger.Fields{
			//	"01.host":       req.Host,
			//	"02.address":    req.RemoteAddr,
			//	"03.method":     req.Method,
			//	"04.requestURI": req.RequestURI,
			//	"05.proto":      req.Proto,
			//	"06.useragent":  req.UserAgent(),
			//	"07.status":     res.Status,
			//	"08.err":        errMsg,
			//	"09.latency":    latencyHuman,
			//	"10.byte_in":    byteIn,
			//	"11.byte_out":   byteOut,
			//}).Infof("HTTP process information")

			logger.WithFields(logger.Fields{
				"00": start.Format("2006-01-02 15:04:05"),
				//"01": req.Host,
				"01": req.Method,
				"02": req.RequestURI,
				"03": req.RemoteAddr,
				//"05": req.Proto,
				"04": res.Status,
				"05": latencyHuman,
				"06": errMsg,
				"07": req.UserAgent(),
				"10": byteIn,
				"11": byteOut,
			}).Infof("HTTP process information")

			return
		}
	}

}

// custom config for CORS
var customCORSConfig = echo_middleware.CORSConfig{
	AllowHeaders:     []string{"Accept", "Content-Type", "Authorization"},
	ExposeHeaders:    []string{echo.HeaderContentDisposition},
	AllowOrigins:     []string{"*"},
	AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
	AllowCredentials: true,
}

// CustomCORS - middleware function for echo's CORS
func CustomCORS() echo.MiddlewareFunc {
	return echo_middleware.CORSWithConfig(customCORSConfig)
}

// CustomJWT - middleware function for JWT
func CustomJWT(secret string) echo.MiddlewareFunc {
	return echo_middleware.JWTWithConfig(echo_middleware.JWTConfig{
		SigningKey:  []byte(secret),
		TokenLookup: "header:Authorization",
	})
}

// Session Interceptor
// func SessionInterceptor(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) (err error) {

// 		if ignoreUrl(c.Request().URL) {
// 			// 예외 URL 이면 통과 시킨다.
// 			return next(c)

// 		} else {

// 			store := lsession.GetStore()
// 			session, _ := store.Get(c.Request(), "liteedge-session")

// 			if key, ok := session.Values["userRoleID"].(string); ok && key != "" {

// 				logger.Debug(fmt.Sprintf("SessionInterceptor >> userId : %v | userRoleId : %s | userName : %s ",
// 					session.Values["userUID"].(int), session.Values["userRoleID"].(string), session.Values["userName"].(string)))

// 				// 세션이 존재하면 by pass
// 				return next(c)
// 			} else {

// 				// 세션이 존재하지 않는 경우 에러 전달
// 				c.Response().Writer.WriteHeader(http.StatusUnauthorized)
// 				message := fmt.Sprintf("{\"error\":true, \"code\": %v, \"message\":%s}",
// 					common.SessionNotFound, common.GetMessageByCode(common.SessionNotFound))
// 				c.Response().Writer.Write([]byte(message))
// 				return nil
// 			}
// 		}
// 	}
// }

func ignoreUrl(url *url.URL) bool {

	var ignores = [...]string{
		"/api/v1/health",
		"/api/v1/auth",
		"/api/v1/auth/changepw",
		"/api/v1/auth/resetpw",
		"/api/v1/device/status",
		"/api/v1/device/start",
		"/api/v1/device/stop",
	}

	for _, path := range ignores {

		if url.Path == path {
			return true
		}
	}
	return false
}

// JwtMiddleware - Handles authentication via jwt's
//func JwtMiddleware(secretKey string, next http.HandlerFunc) http.HandlerFunc {
//	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
//		// Get authentication header
//		authHeader := req.Header.Get("Authorization")
//		if authHeader == "" {
//			rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
//			rw.WriteHeader(http.StatusUnauthorized)
//			return
//		}
//
//		// Check if authentication token is present
//		authHeaderParts := strings.Split(authHeader, " ")
//		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
//			rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
//			rw.WriteHeader(http.StatusUnauthorized)
//			return
//		}
//
//		// Validate authentication token
//		claims, err := common.ParseJWT(secretKey, authHeaderParts[1])
//		if err != nil {
//			rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
//			rw.WriteHeader(http.StatusUnauthorized)
//			return
//		}
//
//		ctx := context.WithValue(req.Context(), common.ContextJWTKey, claims)
//		next.ServeHTTP(rw, req.WithContext(ctx))
//	})
//}
