package engines

import (
	"keys-data-service/channels"
)

func AbortService() error {
	go channels.SignalAbort()

	return nil
}
