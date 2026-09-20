package handlers

import (
	"keys-data-service/libraries"
	"keys-data-service/settings"

	"github.com/gofiber/fiber/v3"
)

func VersionHandler() fiber.Handler {
	return func(fiberContext fiber.Ctx) error {
		defer func() {
			libraries.LogData("keys/handlers.version", "done")
		}()

		libraries.LogTrace("keys/handlers.version", "invoked")

		return fiberContext.JSON(fiber.Map{
			"version": fiber.Map{
				"number":  settings.Version,
				"service": settings.ServiceName,
			},
		})
	}
}
