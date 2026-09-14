package handlers

import (
	"fmt"
	"strconv"

	"keys-data-service/libraries"
	"keys-data-service/settings"
	"keys-data-service/utilities"

	"github.com/gofiber/fiber/v3"

	"keys-data-service/data"
	"keys-data-service/engines"
)

func GetKeysHandler(cache *data.Cache) fiber.Handler {
	return func(fiberContext fiber.Ctx) error {
		libraries.LogTrace("keys/handlers.get-keys", "invoked")

		skip, limit, issue := parsePageData(fiberContext)
		if issue != nil {
			libraries.LogIssue("keys/handlers.get-keys", "parse page data failed",
				map[string]any{
					"error": issue.Error(),
				})

			return utilities.IssueResponse(fiberContext, fiber.StatusBadRequest, issue.Error())
		}

		libraries.LogTrace("keys/handlers.get-keys", "page data parsed",
			map[string]any{
				"skip":  skip,
				"limit": limit,
			})

		keys, issue := engines.GetKeys(cache, skip, limit)
		if issue != nil {
			libraries.LogIssue("keys/handlers.get-keys", "engine call failed",
				map[string]any{
					"error": issue.Error(),
					"skip":  skip,
					"limit": limit,
				})

			return utilities.IssueResponse(fiberContext, fiber.StatusInternalServerError, issue.Error())
		}

		total, issue := engines.GetKeysCount(cache)
		if issue != nil {
			libraries.LogIssue("keys/handlers.get-keys", "engine call failed",
				map[string]any{
					"error": issue.Error(),
				})

			return utilities.IssueResponse(fiberContext, fiber.StatusInternalServerError, issue.Error())
		}

		pages, page := libraries.ComputePageData(skip, limit, total)

		libraries.LogData("keys/handlers.get-keys", "done",
			map[string]any{
				"skip":  skip,
				"limit": limit,
				"total": total,
				"pages": pages,
				"page":  page,
				"count": len(keys),
			})

		return fiberContext.JSON(fiber.Map{
			"keys": keys,

			"skip":  skip,
			"limit": limit,
			"total": total,
			"pages": pages,
			"page":  page,
		})
	}
}

func parsePageData(fiberContext fiber.Ctx) (int, int, error) {
	skip := 0
	limit := 10

	if skipParameter := fiberContext.Query("skip"); skipParameter != "" {
		skipValue, issue := strconv.Atoi(skipParameter)
		if issue != nil {
			return 0, 0, fmt.Errorf("parse page data with issue, invalid skip integer, %w", issue)
		}

		if skipValue < 0 {
			return 0, 0, fmt.Errorf("parse page data with issue, negative skip value")
		}

		skip = skipValue
	}

	if limitParameter := fiberContext.Query("limit"); limitParameter != "" {
		limitValue, issue := strconv.Atoi(limitParameter)
		if issue != nil {
			return 0, 0, fmt.Errorf("parse page data with issue, invalid limit integer, %w", issue)
		}

		if limitValue <= 0 {
			return 0, 0, fmt.Errorf("parse page data with issue, negative limit value")
		}

		if limitValue > settings.DefaultMaximumPageLimit {
			return 0, 0, fmt.Errorf("parse page data with issue, limit value exceeds page limit setting")
		}

		limit = limitValue
	}

	return skip, limit, nil
}
