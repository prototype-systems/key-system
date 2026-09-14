package engines

import (
	"keys-data-service/channels"
)

func StopService() error {
	go channels.SignalStop()

	return nil
}
