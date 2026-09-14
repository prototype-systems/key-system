package routes

import (
	"github.com/gofiber/fiber/v3"

	"keys-data-service/handlers"
)

func RegisterControls(keysGroup fiber.Router) {
	keysGroup.Get("/version",
		handlers.VersionHandler())

	keysGroup.Get("/health",
		handlers.HealthHandler())

	keysGroup.Post("/start",
		handlers.StartServiceHandler())

	keysGroup.Post("/kill",
		handlers.KillServiceHandler())
}
