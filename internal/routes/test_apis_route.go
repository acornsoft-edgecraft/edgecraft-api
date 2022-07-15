package routes

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public routes.
func TestAPIs(a *fiber.App) {
	// Create routes group.
	route := a.Group("/test")

	// Routes for GET method:
	route.Get("/apis", controllers.GetTestAPIs)         // get list of all APIs
	route.Get("/api/:id", controllers.GetTestAPI)       // get one API Response by ID
	route.Get("/jsonbs/:id", controllers.GetTestJsonbs) // get one API Response by ID for test JSONb
	route.Get("/jsonb/:id", controllers.GetTestJsonb)   // get one API Response by ID for test JSONb

	// Routes for POST method:
	route.Post("/api", controllers.NewTestAPI) // create a new API Route

	// Routes for PUT method:
	route.Put("/api/:id", controllers.UpdateTestAPI) // update API Response by ID

	// Routes for DELETE method:
	route.Delete("/api/:id", controllers.DeleteTestAPI) // delete one API Route by ID
}
