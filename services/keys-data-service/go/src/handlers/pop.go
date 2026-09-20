package handlers

import (
	"errors"

	"keys-data-service/libraries"
	"keys-data-service/utilities"

	"github.com/gofiber/fiber/v3"
	"github.com/tidwall/buntdb"

	"keys-data-service/data"
	"keys-data-service/engines"
)

func PopKeyHandler(cache *data.Cache) fiber.Handler {
	return func(fiberContext fiber.Ctx) error {
		reference := fiberContext.Params("reference")

		libraries.LogTrace("keys/handlers.pop-key", "invoked",
			map[string]any{
				"reference": reference,
			})

		if reference == "" {
			libraries.LogIssue("keys/handlers.pop-key", "undefined reference")

			return utilities.IssueResponse(fiberContext, fiber.StatusBadRequest, "pop key with issue, undefined reference")
		}

		key, issue := engines.PopKey(cache, reference)
		if errors.Is(issue, buntdb.ErrNotFound) {
			libraries.LogIssue("keys/handlers.pop-key", "key not found",
				map[string]any{
					"reference": reference,
				})

			return utilities.IssueResponse(fiberContext, fiber.StatusNotFound, "pop key with issue, undetermined key")
		}

		if issue != nil {
			libraries.LogIssue("keys/handlers.pop-key", "engine call failed",
				map[string]any{
					"error":     issue.Error(),
					"reference": reference,
				})

			return utilities.IssueResponse(fiberContext, fiber.StatusInternalServerError, "pop key with issue, "+issue.Error())
		}

		libraries.LogData("keys/handlers.pop-key", "done",
			map[string]any{
				"reference": reference,
			})

		return fiberContext.JSON(fiber.Map{"key": key})
	}
}
