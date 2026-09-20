package handlers

import (
	"keys-data-service/engines"
	"keys-data-service/libraries"
	"keys-data-service/utilities"

	"github.com/gofiber/fiber/v3"
)

func StopServiceHandler() fiber.Handler {
	return func(fiberContext fiber.Ctx) error {
		libraries.LogTrace("keys/handlers.stop-service", "invoked")

		if issue := utilities.OperationResponse(fiberContext, "stop", "initiated"); issue != nil {
			libraries.LogIssue("keys/handlers.stop-service", "operation response failed",
				map[string]any{
					"error": issue.Error(),
				})

			return issue
		}

		if issue := engines.StopService(); issue != nil {
			libraries.LogIssue("keys/handlers.stop-service", "engine call failed",
				map[string]any{
					"error": issue.Error(),
				})

			return issue
		}

		libraries.LogData("keys/handlers.stop-service", "done")

		return nil
	}
}
