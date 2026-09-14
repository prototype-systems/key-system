package handlers

import (
	"github.com/gofiber/fiber/v3"

	"keys-data-service/data"
	"keys-data-service/engines"
	"keys-data-service/libraries"
	"keys-data-service/utilities"
)

func GroupCountKeys(cache *data.Cache) fiber.Handler {
	return func(fiberContext fiber.Ctx) error {
		group := fiberContext.Params("group")

		libraries.LogTrace("keys/handlers.group-count-keys", "invoked",
			map[string]any{
				"group": group,
			})

		if group == "" {
			libraries.LogIssue("keys/handlers.group-count-keys", "undefined group")

			return utilities.IssueResponse(fiberContext, fiber.StatusBadRequest, "group count keys with issue, undefined group")
		}

		count, issue := engines.GroupCountKeys(cache, group)
		if issue != nil {
			libraries.LogIssue("keys/handlers.group-count-keys", "engine call failed",
				map[string]any{
					"error": issue.Error(),
					"group": group,
				})

			return utilities.IssueResponse(fiberContext, fiber.StatusInternalServerError, issue.Error())
		}

		libraries.LogData("keys/handlers.group-count-keys", "done",
			map[string]any{
				"group": group,
				"count": count,
			})

		return fiberContext.JSON(fiber.Map{"count": count})
	}
}
