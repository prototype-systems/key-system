package routes

import (
	"github.com/gofiber/fiber/v3"

	"keys-data-service/handlers"
)

func RegisterService(keysGroup fiber.Router) {
	keysGroup.Get("/version",
		handlers.VersionHandler())

	keysGroup.Get("/health",
		handlers.HealthHandler())

	keysGroup.Post("/stop",
		handlers.StopServiceHandler())

	keysGroup.Post("/abort",
		handlers.AbortServiceHandler())
}
