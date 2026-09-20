package engines

import (
	"fmt"

	"keys-data-service/libraries"
	"keys-data-service/models"
)

func ConfigureLog(setting models.LogSetting) error {
	if issue := libraries.ConfigureLog(libraries.LogSetting{
		LogTerminal:      setting.LogTerminal,
		LogFile:          setting.LogFile,
		LogIpc:           setting.LogIpc,
		LogDirectoryPath: setting.LogDirectoryPath,
		LogDataFilePath:  setting.LogDataFilePath,
		LogIssueFilePath: setting.LogIssueFilePath,
		LogTraceFilePath: setting.LogTraceFilePath,
		LogIpcAddress:    setting.LogIpcAddress,
	}); issue != nil {
		return fmt.Errorf("configure log with issue, %w", issue)
	}

	return nil
}
