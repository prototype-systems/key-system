package routes

import (
	"encoding/json"
	"fmt"

	"keys-data-service/handlers"
	"keys-data-service/models"
	"keys-data-service/settings"
)

func Invoke(data string) error {
	var message models.Message

	if issue := json.Unmarshal([]byte(data), &message); issue != nil {
		return fmt.Errorf("invoke parse message with issue, %w", issue)
	}

	switch message.Category {
	case settings.ConfigureLogCategory:
		go func() {
			if issue := handlers.ConfigureLogHandler(message.Data); issue != nil {
				panic("configure log on invoke with issue, " + issue.Error())
			}
		}()

	default:
		return fmt.Errorf("unknown invoke category %q", message.Category)
	}

	return nil
}
