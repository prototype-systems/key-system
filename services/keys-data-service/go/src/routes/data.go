package routes

import (
	"github.com/gofiber/fiber/v3"

	"keys-data-service/data"
	"keys-data-service/handlers"
)

func RegisterData(keysGroup fiber.Router, cache *data.Cache) {
	keysGroup.Post("/add",
		handlers.AddKeyHandler(cache))

	keysGroup.Post("/push",
		handlers.PushKeyHandler(cache))

	keysGroup.Delete("/pop/:reference",
		handlers.PopKeyHandler(cache))

	keysGroup.Get("/",
		handlers.GetKeysHandler(cache))

	keysGroup.Get("/:group/count",
		handlers.GroupCountKeys(cache))

	keysGroup.Get("/:group/:key",
		handlers.GroupGetKey(cache))

	keysGroup.Post("/:group/:key",
		handlers.GroupSetKey(cache))

	keysGroup.Get("/:reference",
		handlers.GetKeyHandler(cache))

	keysGroup.Post("/persist",
		handlers.PersistHandler(cache))
}
