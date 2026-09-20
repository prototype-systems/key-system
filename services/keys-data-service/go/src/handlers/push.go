package handlers

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/tidwall/buntdb"

	"keys-data-service/data"
	"keys-data-service/engines"
	"keys-data-service/libraries"
	"keys-data-service/models"
	"keys-data-service/utilities"
)

func PushKeyHandler(cache *data.Cache) fiber.Handler {
	return func(fiberContext fiber.Ctx) error {
		libraries.LogTrace("keys/handlers.push-key", "invoked")

		key, issue := parsePushKey(fiberContext)
		if issue != nil {
			libraries.LogIssue("keys/handlers.push-key", "parse key failed",
				map[string]any{
					"error": issue.Error(),
				})

			return utilities.IssueResponse(fiberContext, fiber.StatusBadRequest, issue.Error())
		}

		libraries.LogTrace("keys/handlers.push-key", "key parsed",
			map[string]any{
				"reference": key.Reference,
			})

		if issue = validatePushKey(key); issue != nil {
			libraries.LogIssue("keys/handlers.push-key", "validate key failed",
				map[string]any{
					"error":     issue.Error(),
					"reference": key.Reference,
				})

			return utilities.IssueResponse(fiberContext, fiber.StatusBadRequest, issue.Error())
		}

		if issue = checkPushKey(cache, key.Reference); issue != nil {
			libraries.LogIssue("keys/handlers.push-key", "check existing key failed",
				map[string]any{
					"error":     issue.Error(),
					"reference": key.Reference,
				})

			return utilities.IssueResponse(fiberContext, fiber.StatusConflict, issue.Error())
		}

		if issue = engines.PushKey(cache, *key); issue != nil {
			libraries.LogIssue("keys/handlers.push-key", "engine call failed",
				map[string]any{
					"error":     issue.Error(),
					"reference": key.Reference,
				})

			return utilities.IssueResponse(fiberContext, fiber.StatusInternalServerError, issue.Error())
		}

		libraries.LogData("keys/handlers.push-key", "done",
			map[string]any{
				"reference": key.Reference,
				"name":      key.Name,
			})

		return fiberContext.JSON(fiber.Map{"reference": key.Reference})
	}
}

func parsePushKey(fiberContext fiber.Ctx) (*models.Key, error) {
	var key models.Key

	if issue := fiberContext.Bind().JSON(&key); issue != nil {
		return nil, fmt.Errorf("parse push key with issue, %w", issue)
	}

	return &key, nil
}

func validatePushKey(key *models.Key) error {
	if key.Reference == "" {
		return fmt.Errorf("validate push key with issue, undefined key reference")
	}

	if key.Name == "" {
		return fmt.Errorf("validate push key with issue, undefined key name")
	}

	if libraries.CheckNullValue(key.Value) {
		return fmt.Errorf("validate push key with issue, undefined key value")
	}

	return nil
}

func checkPushKey(cache *data.Cache, reference string) error {
	return cache.DB().View(func(transaction *buntdb.Tx) error {
		key, issue := transaction.Get("key:" + reference)
		if errors.Is(issue, buntdb.ErrNotFound) {
			return nil
		}

		if issue != nil {
			return fmt.Errorf("check push key with issue, %w", issue)
		}

		if key != "" {
			return fmt.Errorf("check push key with issue, key exists")
		}

		return fmt.Errorf("check push key with undetermined issue")
	})
}
