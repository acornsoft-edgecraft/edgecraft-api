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
	v1.GET("/clouds/:cloudUid", api.SelectCloudHandler)
	v1.POST("/clouds", api.RegisterCloudHandler)
	v1.PUT("/clouds/:cloudUid", api.UpdateCloudHandler)
	v1.DELETE("/clouds/:cloudUid", api.DeleteCloudHandler)

	// Cloud - Cluster
	v1.GET("/clouds/:cloudID/clusters", api.AllCloudClusterListHandler)

	// Cloud - Node
	v1.GET("/clouds/:cloudID/nodes", api.AllCloudNodeListHandler)
	v1.GET("/clouds/:cloudID/nodes", api.GetCloudNodeHandler)
	// v1.POST("/clouds/:cloudID/nodes", api.)
	v1.DELETE("/clouds/:cloudID/nodes", api.AllCloudNodeListHandler)

	// CodeGroup
	v1.GET("/codegroups", api.AllCodeGroupListHandler)
	v1.POST("/codegroups", api.RegisterCodeGroupHandler)
	v1.POST("/codegroups/search", api.CodeGroupSearchHandler)
	// v1.PUT("/cgroups/:cgroupUid", api.UpdateCloudHandler)
	// v1.DELETE("/cgroups/:cgroupsUid", api.DeleteCloudHandler)

	// Code
	v1.GET("/codes", api.AllCloudListHandler)
}
