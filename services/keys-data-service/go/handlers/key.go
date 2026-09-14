package handlers

import (
	"github.com/gofiber/fiber/v3"

	"keys-data-service/data"
	"keys-data-service/engines"
	"keys-data-service/libraries"
	"keys-data-service/utilities"
)

func GetKeyHandler(cache *data.Cache) fiber.Handler {
	return func(fiberContext fiber.Ctx) error {
		reference := fiberContext.Params("reference")

		libraries.LogTrace("keys/handlers.get-key", "invoked")

		if reference == "" {
			libraries.LogIssue("keys/handlers.get-key", "undefined reference")

			return utilities.IssueResponse(fiberContext, fiber.StatusBadRequest, "get key with issue, undefined key reference")
		}

		key, issue := engines.GetKey(cache, reference)
		if issue != nil {
			libraries.LogIssue("keys/handlers.get-key", "engine call failed",
				map[string]any{
					"error":     issue.Error(),
					"reference": reference,
				})

			return utilities.IssueResponse(fiberContext, fiber.StatusNotFound, issue.Error())
		}

		libraries.LogData("keys/handlers.get-key", "done",
			map[string]any{
				"reference": reference,
			})

		return fiberContext.JSON(fiber.Map{"key": key})
	}
}
