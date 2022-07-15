package main

import (
	"os"

	"github.com/acornsoft-edgecraft/edgecraft-api/internal/routes"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/configs"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/middleware"

	server "github.com/acornsoft-edgecraft/edgecraft-api/cmd/edgecraft-apiserver/app"

	"github.com/gofiber/fiber/v2"

	_ "github.com/acornsoft-edgecraft/edgecraft-api/api" // load API Docs files (Swagger)

	_ "github.com/joho/godotenv/autoload" // load .env file automatically
)

// @title                       API
// @version                     1.0
// @description                 This is an auto-generated API Docs.
// @termsOfService              http://swagger.io/terms/
// @contact.name                API Support
// @contact.email               your@mail.com
// @license.name                Apache 2.0
// @license.url                 http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath                    /api
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
func main() {
	// Define Fiber config.
	config := configs.FiberConfig()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	// Middlewares.
	middleware.FiberMiddleware(app) // Register Fiber's middleware(cors, logger) for app.

	// Routes.
	routes.SwaggerRoute(app)       // Register a route for API Docs (Swagger).
	routes.PublicRoutes(app)       // Register a public routes for app.
	routes.PrivateRoutes(app)      // Register a private routes for app.
	routes.SetupProductRoutes(app) // Register a SetupProductRoutes route for app (used ORM(Gorm)).
	routes.TestAPIs(app)           // Register a TestRoutes Test API Response route for app (used ORM(Gorm)).
	routes.MockServerRoutes(app)   // Register a TestRoutes Test API Response route for app (used ORM(Gorm)).

	// Not Found Route.
	routes.NotFoundRoute(app) // Register route for 404 Error.

	// Start server (with or without graceful shutdown).
	if os.Getenv("STAGE_STATUS") == "dev" {
		server.StartServer(app)
	} else {
		server.StartServerWithGracefulShutdown(app)
	}
}
