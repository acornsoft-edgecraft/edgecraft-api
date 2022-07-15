package controllers

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/internal/entities"
	"github.com/acornsoft-edgecraft/edgecraft-api/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func NewTestAPI(ctx *fiber.Ctx) error {
	var newTestAPI *entities.TestAPI
	if err := ctx.BodyParser(&newTestAPI); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	testAPI, err := services.NewTestAPI(newTestAPI)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(testAPI)
}

// GetTestAPI func gets product by given ID or 404 error.
// @Description  Get product by given ID.
// @Summary      get product by given ID
// @Tags         TestAPI
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "TestAPI ID"
// @Success      200  {object}  entities.TestAPI
// @Router       /test/api/{id} [get]
func GetTestAPI(ctx *fiber.Ctx) error {
	// Catch API ID from URL.
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Service
	testAPI, err := services.GetTestAPI(id)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(testAPI)
}

func GetTestJsonbs(ctx *fiber.Ctx) error {
	// Catch API ID from URL.
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Service
	testAPI, err := services.GetTestJsonbs(id)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(testAPI)
}

func GetTestJsonb(ctx *fiber.Ctx) error {
	// Catch API ID from URL.
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Catch query parameters
	parm := ctx.Query("name")

	// Service
	testAPI, err := services.GetTestJsonb(id, parm)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(testAPI)
}

// GetTestAPIs func gets all exists products.
// @Description  Get all exists products.
// @Summary      get all exists products
// @Tags         TestAPIs
// @Accept       json
// @Produce      json
// @Success      200  {array} entities.TestAPI
// @Router       /test/apis [get]
func GetTestAPIs(ctx *fiber.Ctx) error {
	testAPIs, err := services.GetTestAPIs()
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(testAPIs)
}

func UpdateTestAPI(ctx *fiber.Ctx) error {
	var testAPI *entities.TestAPI
	if err := ctx.BodyParser(&testAPI); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	testAPI, err := services.UpdateTestAPI(testAPI, *testAPI.ID)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(testAPI)
}

func DeleteTestAPI(ctx *fiber.Ctx) error {
	// Catch API ID from URL.
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Service
	testAPI, err := services.DeleteTestAPI(id)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(testAPI)
}
