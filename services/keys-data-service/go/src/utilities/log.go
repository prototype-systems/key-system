package utilities

import (
	"keys-data-service/settings"
	"os"
	"path/filepath"
)

func GetLogDirectoryPath() string {
	logDirectory := os.Getenv(settings.LogDirectoryVariable)
	if logDirectory == "" {
		logDirectory = settings.DefaultLogDirectory
	}

	return filepath.Join(filepath.Dir(GetExecutablePath()), logDirectory)
}

func GetLogDataFilePath() string {
	fileName := os.Getenv(settings.LogDataFileVariable)
	if fileName == "" {
		fileName = settings.DefaultLogDataFile
	}

	return filepath.Join(GetLogDirectoryPath(), fileName)
}

func GetLogIssueFilePath() string {
	fileName := os.Getenv(settings.LogIssueFileVariable)
	if fileName == "" {
		fileName = settings.DefaultLogIssueFile
	}

	return filepath.Join(GetLogDirectoryPath(), fileName)
}

func GetLogTraceFilePath() string {
	fileName := os.Getenv(settings.LogTraceFileVariable)
	if fileName == "" {
		fileName = settings.DefaultLogTraceFile
	}

	return filepath.Join(GetLogDirectoryPath(), fileName)
}
