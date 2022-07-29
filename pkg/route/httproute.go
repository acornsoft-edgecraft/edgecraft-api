package route

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/middleware"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/server"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	echo_swagger "github.com/swaggo/echo-swagger"
)

// SetRoutes - api routes setting
func SetRoutes(api *api.API, server *server.Instance) {
	setSwagger(api, server)
	setHTTPRoutes(api, server)
}

func setSwagger(api *api.API, server *server.Instance) {
	server.HTTPServer.GET("/swagger/*", echo_swagger.EchoWrapHandler(echo_swagger.URL("doc.json")))
}

func setHTTPRoutes(api *api.API, server *server.Instance) {
	server.HTTPServer.Use(middleware.CustomCORS())

	server.HTTPServer.Use(middleware.CustomLogger())

	// create a sub route group
	v1 := server.HTTPServer.Group("/api/v1")

	// Middleware, default
	v1.Use(echo_middleware.Recover())

	// Session Interceptor
	// v1.Use(middleware.SessionInterceptor)

	v1.OPTIONS("", middleware.Preflight)

	// format-swagger:route [method] [path pattern] [?tag1 tag2 tag3] [operation id]
	v1.GET("/health", middleware.HealthCheck)

	// Cloud
	v1.GET("/clouds", api.AllCloudListHandler)

}
