package handlers

import (
	"encoding/json"
	"fmt"

	"keys-data-service/engines"
	"keys-data-service/libraries"
	"keys-data-service/models"
)

func ConfigureLogHandler(data json.RawMessage) error {
	libraries.LogTrace("keys/handlers.configure-log", "invoked")

	var logSetting models.LogSetting

	if issue := json.Unmarshal([]byte(data), &logSetting); issue != nil {
		libraries.LogIssue("keys/handlers.configure-log", "parse setting failed",
			map[string]any{
				"error": issue.Error(),
			})

		return fmt.Errorf("configure log on parse message with issue, %w", issue)
	}

	if issue := engines.ConfigureLog(logSetting); issue != nil {
		libraries.LogIssue("keys/handlers.configure-log", "engine call failed",
			map[string]any{
				"error": issue.Error(),
			})

		return issue
	}

	libraries.LogData("keys/handlers.configure-log", "done")

	return nil
}
