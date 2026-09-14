package utilities

import (
	"keys-data-service/settings"

	"github.com/gofiber/fiber/v3"
)

func IssueResponse(fiberContext fiber.Ctx, statusCode int, description string) error {
	return fiberContext.Status(statusCode).JSON(fiber.Map{
		"issue": fiber.Map{
			"description": description,
			"method":      fiberContext.Method(),
			"path":        fiberContext.Path(),
		},

		"service": settings.ServiceName,
		"version": settings.Version,
	})
}

func OperationResponse(fiberContext fiber.Ctx, procedure, status string) error {
	return fiberContext.JSON(fiber.Map{
		"operation": fiber.Map{
			"status":    status,
			"procedure": procedure,
		},

		"service": settings.ServiceName,
		"version": settings.Version,
	})
}
