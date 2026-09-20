package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v3"

	"keys-data-service/data"
	"keys-data-service/engines"
	"keys-data-service/libraries"
	"keys-data-service/utilities"
)

func GroupSetKey(cache *data.Cache) fiber.Handler {
	return func(fiberContext fiber.Ctx) error {
		group := fiberContext.Params("group")
		name := fiberContext.Params("key")

		libraries.LogTrace("keys/handlers.group-set-key", "invoked",
			map[string]any{
				"group": group,
				"name":  name,
			})

		if group == "" {
			libraries.LogIssue("keys/handlers.group-set-key", "undefined group")

			return utilities.IssueResponse(fiberContext, fiber.StatusBadRequest, "group set key with issue, undefined group")
		}

		if name == "" {
			libraries.LogIssue("keys/handlers.group-set-key", "undefined key name")

			return utilities.IssueResponse(fiberContext, fiber.StatusBadRequest, "group set key with issue, undefined key name")
		}

		value, issue := parseGroupSetKeyValue(fiberContext)
		if issue != nil {
			libraries.LogIssue("keys/handlers.group-set-key", "parse value failed",
				map[string]any{
					"error": issue.Error(),
					"group": group,
					"name":  name,
				})

			return utilities.IssueResponse(fiberContext, fiber.StatusBadRequest, issue.Error())
		}

		reference, issue := engines.GroupSetKey(cache, group, name, value)
		if issue != nil {
			libraries.LogIssue("keys/handlers.group-set-key", "engine call failed",
				map[string]any{
					"error": issue.Error(),
					"group": group,
					"name":  name,
				})

			return utilities.IssueResponse(fiberContext, fiber.StatusInternalServerError, issue.Error())
		}

		libraries.LogData("keys/handlers.group-set-key", "done",
			map[string]any{
				"reference": reference,
				"group":     group,
				"name":      name,
			})

		return fiberContext.JSON(fiber.Map{"reference": reference})
	}
}

func parseGroupSetKeyValue(fiberContext fiber.Ctx) (json.RawMessage, error) {
	var body struct {
		Value json.RawMessage `json:"value"`
	}

	if issue := fiberContext.Bind().JSON(&body); issue != nil {
		return nil, fmt.Errorf("parse group set key value with issue, %w", issue)
	}

	if libraries.CheckNullValue(body.Value) {
		return nil, fmt.Errorf("parse group set key value with issue, undefined key value")
	}

	return body.Value, nil
}
