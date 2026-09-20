package handlers

import (
	"keys-data-service/libraries"
	"keys-data-service/utilities"

	"github.com/gofiber/fiber/v3"

	"keys-data-service/data"
	"keys-data-service/engines"
)

func PersistHandler(cache *data.Cache) fiber.Handler {
	return func(fiberContext fiber.Ctx) error {
		libraries.LogTrace("keys/handlers.persist", "invoked")

		if issue := engines.PersistCache(cache); issue != nil {
			libraries.LogIssue("keys/handlers.persist", "engine call failed",
				map[string]any{
					"error": issue.Error(),
				})

			return utilities.IssueResponse(fiberContext, fiber.StatusInternalServerError, issue.Error())
		}

		libraries.LogData("keys/handlers.persist", "done")

		return utilities.OperationResponse(fiberContext, "persist", "completed")
	}
}
