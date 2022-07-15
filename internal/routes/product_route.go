package routes

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupProductRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	// Routes for GET method:
	route.Get("/products", controllers.GetProducts)   // get list of all books
	route.Get("/product/:id", controllers.GetProduct) // get one book by ID

	// Routes for POST method:
	route.Post("/product", controllers.NewProduct) // new product item

	// Routes for PUT method:
	route.Put("/product/:id", controllers.UpdateProduct) // update product item

	// Routes for DELETE method:
	route.Delete("/product/:id", controllers.DeleteProduct) // delete product item
}
