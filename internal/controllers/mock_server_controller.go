package controllers

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/internal/services"
	"github.com/gofiber/fiber/v2"
)

func GetMockRoute(ctx *fiber.Ctx) error {

	// Service
	getMockRoute, err := services.GetMockRoute(ctx.Path(), ctx.Method())
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(getMockRoute.API_response.Data)
}
