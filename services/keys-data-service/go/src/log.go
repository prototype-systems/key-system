package main

import (
	"keys-data-service/libraries"
	"keys-data-service/utilities"
)

func RunLog() {
	setting, issue := libraries.SetupLog(libraries.LogSetting{
		LogTerminal:      true,
		LogFile:          false,
		LogIpc:           false,
		LogDirectoryPath: utilities.GetLogDirectoryPath(),
		LogDataFilePath:  utilities.GetLogDataFilePath(),
		LogIssueFilePath: utilities.GetLogIssueFilePath(),
		LogTraceFilePath: utilities.GetLogTraceFilePath(),
		LogIpcAddress:    utilities.GetIpcSocketAddress(),
	})

	if issue != nil {
		libraries.LogIssue("keys/log", "run log setup with issue", map[string]any{"issue": issue})

		utilities.KillService()
	}

	libraries.SetLogSetting(setting)
}
