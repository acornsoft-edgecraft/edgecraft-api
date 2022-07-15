package routes

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public routes.
func MockServerRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/mock")

	// Routes for GET All:
	route.All("/*", controllers.GetMockRoute)
}
