package controllers

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/internal/entities"
	"github.com/acornsoft-edgecraft/edgecraft-api/internal/services"
	"github.com/gofiber/fiber/v2"
)

func NewProduct(ctx *fiber.Ctx) error {
	var newProduct *entities.Product
	if err := ctx.BodyParser(&newProduct); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	product, err := services.NewProduct(newProduct)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(product)
}

// GetProduct func gets product by given ID or 404 error.
// @Description  Get product by given ID.
// @Summary      get product by given ID
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  entities.Product
// @Router       /v1/product/{id} [get]
func GetProduct(ctx *fiber.Ctx) error {
	product, err := services.GetProduct(ctx.Params("id"))
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(product)
}

// GetProducts func gets all exists products.
// @Description  Get all exists products.
// @Summary      get all exists products
// @Tags         Products
// @Accept       json
// @Produce      json
// @Success      200  {array} entities.Product
// @Router       /v1/products [get]
func GetProducts(ctx *fiber.Ctx) error {
	products, err := services.GetProducts()
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(products)
}

func UpdateProduct(ctx *fiber.Ctx) error {
	var product *entities.Product
	if err := ctx.BodyParser(&product); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	product, err := services.UpdateProduct(product, ctx.Params("id"))
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(product)
}

func DeleteProduct(ctx *fiber.Ctx) error {
	product, err := services.DeleteProduct(ctx.Params("id"))
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(product)
}
