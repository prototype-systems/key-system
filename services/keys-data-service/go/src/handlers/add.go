package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v3"

	"keys-data-service/data"
	"keys-data-service/engines"
	"keys-data-service/libraries"
	"keys-data-service/models"
	"keys-data-service/utilities"
)

func AddKeyHandler(cache *data.Cache) fiber.Handler {
	return func(fiberContext fiber.Ctx) error {
		libraries.LogTrace("keys/handlers.add-key", "invoked")

		key, issue := parseAddKey(fiberContext)
		if issue != nil {
			libraries.LogIssue("keys/handlers.add-key", "parse key failed",
				map[string]any{
					"error": issue.Error(),
				})

			return utilities.IssueResponse(fiberContext, fiber.StatusBadRequest, issue.Error())
		}

		libraries.LogTrace("keys/handlers.add-key", "key parsed")

		if issue = validateAddKey(key); issue != nil {
			libraries.LogIssue("keys/handlers.add-key", "validate key failed",
				map[string]any{
					"error": issue.Error(),
				})

			return utilities.IssueResponse(fiberContext, fiber.StatusBadRequest, issue.Error())
		}

		reference, issue := engines.AddKey(cache, *key)
		if issue != nil {
			libraries.LogIssue("keys/handlers.add-key", "engine call failed",
				map[string]any{
					"error": issue.Error(),
				})

			return utilities.IssueResponse(fiberContext, fiber.StatusInternalServerError, issue.Error())
		}

		libraries.LogData("keys/handlers.add-key", "done",
			map[string]any{
				"reference": reference,
			})

		return fiberContext.JSON(fiber.Map{"reference": reference})
	}
}

func parseAddKey(fiberContext fiber.Ctx) (*models.Key, error) {
	var key models.Key

	if issue := fiberContext.Bind().JSON(&key); issue != nil {
		return nil, fmt.Errorf("parse add key with issue, %w", issue)
	}

	return &key, nil
}

func validateAddKey(key *models.Key) error {
	if key.Name == "" {
		return fmt.Errorf("validate add key with issue, undefined key name")
	}

	if libraries.CheckNullValue(key.Value) {
		return fmt.Errorf("validate add key with issue, undefined key value")
	}

	return nil
}
