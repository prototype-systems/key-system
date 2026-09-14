package handlers

import (
	"github.com/gofiber/fiber/v3"

	"keys-data-service/data"
	"keys-data-service/engines"
	"keys-data-service/libraries"
	"keys-data-service/utilities"
)

func GroupGetKey(cache *data.Cache) fiber.Handler {
	return func(fiberContext fiber.Ctx) error {
		group := fiberContext.Params("group")
		name := fiberContext.Params("key")

		libraries.LogTrace("keys/handlers.group-get-key", "invoked",
			map[string]any{
				"group": group,
				"name":  name,
			})

		if group == "" {
			libraries.LogIssue("keys/handlers.group-get-key", "undefined group")

			return utilities.IssueResponse(fiberContext, fiber.StatusBadRequest, "group get key with issue, undefined group")
		}

		if name == "" {
			libraries.LogIssue("keys/handlers.group-get-key", "undefined key name")

			return utilities.IssueResponse(fiberContext, fiber.StatusBadRequest, "group get key with issue, undefined key name")
		}

		value, issue := engines.GroupGetKey(cache, group, name)
		if issue != nil {
			libraries.LogIssue("keys/handlers.group-get-key", "engine call failed",
				map[string]any{
					"error": issue.Error(),
					"group": group,
					"name":  name,
				})

			return utilities.IssueResponse(fiberContext, fiber.StatusInternalServerError, issue.Error())
		}

		if value == nil {
			libraries.LogIssue("keys/handlers.group-get-key", "key not found",
				map[string]any{
					"group": group,
					"name":  name,
				})

			return utilities.IssueResponse(fiberContext, fiber.StatusNotFound, "group get key with issue, key not found")
		}

		libraries.LogData("keys/handlers.group-get-key", "done",
			map[string]any{
				"group": group,
				"name":  name,
			})

		return fiberContext.JSON(fiber.Map{"value": value})
	}
}
