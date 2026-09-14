package engines

import (
	"keys-data-service/channels"
)

func KillService() error {
	go channels.SignalKill()

	return nil
}
