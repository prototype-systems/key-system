package handlers

import (
	"keys-data-service/engines"
	"keys-data-service/libraries"
	"keys-data-service/utilities"

	"github.com/gofiber/fiber/v3"
)

func AbortServiceHandler() fiber.Handler {
	return func(fiberContext fiber.Ctx) error {
		libraries.LogTrace("keys/handlers.abort-service", "invoked")

		if issue := utilities.OperationResponse(fiberContext, "abort", "initiated"); issue != nil {
			libraries.LogIssue("keys/handlers.abort-service", "operation response failed",
				map[string]any{
					"error": issue.Error(),
				})

			return issue
		}

		if issue := engines.AbortService(); issue != nil {
			libraries.LogIssue("keys/handlers.abort-service", "engine call failed",
				map[string]any{
					"error": issue.Error(),
				})

			return issue
		}

		libraries.LogData("keys/handlers.abort-service", "done")

		return nil
	}
}
