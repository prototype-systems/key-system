package models

type LogSetting struct {
	LogTerminal      bool   `json:"logTerminal,omitempty"`
	LogFile          bool   `json:"logFile,omitempty"`
	LogIpc           bool   `json:"logIpc,omitempty"`
	LogDirectoryPath string `json:"logDirectoryPath,omitempty"`
	LogDataFilePath  string `json:"logDataFilePath,omitempty"`
	LogIssueFilePath string `json:"logIssueFilePath,omitempty"`
	LogTraceFilePath string `json:"logTraceFilePath,omitempty"`
	LogIpcAddress    string `json:"logIpcAddress,omitempty"`
}
