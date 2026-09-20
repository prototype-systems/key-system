package engines

import (
	"keys-data-service/channels"
)

func StartService() error {
	go channels.SignalStart()

	return nil
}
