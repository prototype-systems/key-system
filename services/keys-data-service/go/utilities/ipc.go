package utilities

import (
	"keys-data-service/settings"
	"path/filepath"
)

func GetIpcSocketAddress() string {
	if CheckWindows() {
		return settings.WindowsIpcSocketAddress
	}

	return filepath.Join(filepath.Dir(GetExecutablePath()), settings.IpcSocketAddress)
}
