package handlers

import (
	"keys-data-service/engines"
	"keys-data-service/libraries"
	"keys-data-service/utilities"

	"github.com/gofiber/fiber/v3"
)

func KillServiceHandler() fiber.Handler {
	return func(fiberContext fiber.Ctx) error {
		libraries.LogTrace("keys/handlers.kill-service", "invoked")

		if issue := utilities.OperationResponse(fiberContext, "kill", "initiated"); issue != nil {
			libraries.LogIssue("keys/handlers.kill-service", "operation response failed",
				map[string]any{
					"error": issue.Error(),
				})

			return issue
		}

		if issue := engines.KillService(); issue != nil {
			libraries.LogIssue("keys/handlers.kill-service", "engine call failed",
				map[string]any{
					"error": issue.Error(),
				})

			return issue
		}

		libraries.LogData("keys/handlers.kill-service", "done")

		return nil
	}
}
