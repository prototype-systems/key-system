package utilities

import (
	"keys-data-service-test/settings"
	"os"
)

func ResolveEnvironment() []string {
	return append(os.Environ(),
		"KEYS_DATA_SERVICE_PORT="+settings.TargetPort,
		"KEYS_CACHE_PATH="+settings.TargetTestDirectory,
	)
}
