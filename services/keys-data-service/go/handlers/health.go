package handlers

import (
	"keys-data-service/libraries"
	"keys-data-service/settings"

	"github.com/gofiber/fiber/v3"
)

func HealthHandler() fiber.Handler {
	return func(fiberContext fiber.Ctx) error {
		defer func() {
			libraries.LogData("keys/handlers.health", "done")
		}()

		libraries.LogTrace("keys/handlers.health", "invoked")

		return fiberContext.JSON(fiber.Map{
			"health": fiber.Map{
				"status": "healthy",
			},

			"service": settings.ServiceName,
			"version": settings.Version,
		})
	}
}
