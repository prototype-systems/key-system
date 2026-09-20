package utilities

import (
	"keys-data-service/channels"
	"keys-data-service/settings"
	"os"
)

func GetPort() string {
	port := os.Getenv(settings.ServicePortVariable)
	if port == "" {
		port = settings.DefaultPort
	}

	return port
}

func KillService() error {
	go channels.SignalKill()

	return nil
}
