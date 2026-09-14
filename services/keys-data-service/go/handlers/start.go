package handlers

import (
	"keys-data-service/engines"
	"keys-data-service/libraries"
	"keys-data-service/utilities"

	"github.com/gofiber/fiber/v3"
)

func StartServiceHandler() fiber.Handler {
	return func(fiberContext fiber.Ctx) error {
		libraries.LogTrace("keys/handlers.start-service", "invoked")

		if issue := utilities.OperationResponse(fiberContext, "start", "initiated"); issue != nil {
			libraries.LogIssue("keys/handlers.start-service", "operation response failed",
				map[string]any{
					"error": issue.Error(),
				})

			return issue
		}

		if issue := engines.StartService(); issue != nil {
			libraries.LogIssue("keys/handlers.start-service", "engine call failed",
				map[string]any{
					"error": issue.Error(),
				})

			return issue
		}

		libraries.LogData("keys/handlers.start-service", "done")

		return nil
	}
}
