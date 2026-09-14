package libraries

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	LogDataType  = "D"
	LogIssueType = "I"
	LogTraceType = "T"
)

const (
	LogDataCategory  = "log-data"
	LogIssueCategory = "log-issue"
	LogTraceCategory = "log-trace"
)

type LogSetting struct {
	LogTerminal      bool
	LogFile          bool
	LogIpc           bool
	LogDirectoryPath string
	LogDataFilePath  string
	LogIssueFilePath string
	LogTraceFilePath string
	LogIpcAddress    string
}

type LogMessage struct {
	Scope   string
	Message string
	Data    json.RawMessage
}

type Message struct {
	Category string
	Data     json.RawMessage
}

var (
	logSetting        *LogSetting
	logSettingControl sync.RWMutex
)

var defaultLogSetting = LogSetting{
	LogTerminal: true,
}

func SetLogSetting(setting *LogSetting) {
	logSettingControl.Lock()
	defer logSettingControl.Unlock()

	logSetting = setting
}

func GetLogSetting() *LogSetting {
	logSettingControl.RLock()
	defer logSettingControl.RUnlock()

	if logSetting == nil {
		return &defaultLogSetting
	}

	return logSetting
}

func SetupLog(setting LogSetting) (*LogSetting, error) {
	logDirectory, issue := ResolveDirectory(setting.LogDirectoryPath)
	if issue != nil {
		return nil, fmt.Errorf("setup log resolve directory with issue, %w", issue)
	}

	logDataFile, issue := ResolveFile(setting.LogDataFilePath)
	if issue != nil {
		return nil, fmt.Errorf("setup log resolve data file with issue, %w", issue)
	}

	logIssueFile, issue := ResolveFile(setting.LogIssueFilePath)
	if issue != nil {
		return nil, fmt.Errorf("setup log resolve issue file with issue, %w", issue)
	}

	logTraceFile, issue := ResolveFile(setting.LogTraceFilePath)
	if issue != nil {
		return nil, fmt.Errorf("setup log resolve trace file with issue, %w", issue)
	}

	return &LogSetting{
		LogTerminal:      setting.LogTerminal,
		LogFile:          setting.LogFile,
		LogIpc:           setting.LogIpc,
		LogDirectoryPath: logDirectory,
		LogDataFilePath:  logDataFile,
		LogIssueFilePath: logIssueFile,
		LogTraceFilePath: logTraceFile,
		LogIpcAddress:    setting.LogIpcAddress,
	}, nil
}

func ConfigureLog(setting LogSetting) error {
	logDirectory, issue := ResolveDirectory(setting.LogDirectoryPath)
	if issue != nil {
		return fmt.Errorf("configure log resolve directory with issue, %w", issue)
	}
	setting.LogDirectoryPath = logDirectory

	logDataFile, issue := ResolveFile(setting.LogDataFilePath)
	if issue != nil {
		return fmt.Errorf("configure log resolve data file with issue, %w", issue)
	}
	setting.LogDataFilePath = logDataFile

	logIssueFile, issue := ResolveFile(setting.LogIssueFilePath)
	if issue != nil {
		return fmt.Errorf("configure log resolve issue file with issue, %w", issue)
	}
	setting.LogIssueFilePath = logIssueFile

	logTraceFile, issue := ResolveFile(setting.LogTraceFilePath)
	if issue != nil {
		return fmt.Errorf("configure log resolve trace file with issue, %w", issue)
	}
	setting.LogTraceFilePath = logTraceFile

	SetLogSetting(&setting)

	return nil
}

func formatLogMessage(
	logType string,
	reference string,
	scope string,
	message string,
) string {
	return fmt.Sprintf("# %s %s %s %s %s",
		time.Now().UTC().Format("2006-01-02T15:04:05.00Z"),
		logType, scope, message,
		fmt.Sprintf("urn:cla:%s:%s", scope, reference))
}

func LogData(scope string, message string, data ...map[string]any) error {
	setting := GetLogSetting()

	if setting == nil {
		return fmt.Errorf("log data setting undefined")
	}

	issueChannel := make(chan error, 2)

	if setting.LogTerminal {
		go func() {
			issueChannel <- LogTerminalData(setting, scope, message, data...)
		}()
	}

	if setting.LogFile {
		go func() {
			issueChannel <- LogFileData(setting, scope, message, data...)
		}()
	}

	if setting.LogIpc {
		go func() {
			issueChannel <- LogIpcData(setting, scope, message, data...)
		}()
	}

	if issue := <-issueChannel; issue != nil {
		return issue
	}

	return nil
}

func LogIssue(scope string, message string, data ...map[string]any) error {
	setting := GetLogSetting()

	if setting == nil {
		return fmt.Errorf("log issue setting undefined")
	}

	issueChannel := make(chan error, 2)

	if setting.LogTerminal {
		go func() {
			issueChannel <- LogTerminalIssue(setting, scope, message, data...)
		}()
	}

	if setting.LogFile {
		go func() {
			issueChannel <- LogFileIssue(setting, scope, message, data...)
		}()
	}

	if setting.LogIpc {
		go func() {
			issueChannel <- LogIpcIssue(setting, scope, message, data...)
		}()
	}

	if issue := <-issueChannel; issue != nil {
		return issue
	}

	return nil
}

func LogTrace(scope string, message string, data ...map[string]any) error {
	setting := GetLogSetting()

	if setting == nil {
		return fmt.Errorf("log trace setting undefined")
	}

	issueChannel := make(chan error, 2)

	if setting.LogTerminal {
		go func() {
			issueChannel <- LogTerminalTrace(setting, scope, message, data...)
		}()
	}

	if setting.LogFile {
		go func() {
			issueChannel <- LogFileTrace(setting, scope, message, data...)
		}()
	}

	if setting.LogIpc {
		go func() {
			issueChannel <- LogIpcTrace(setting, scope, message, data...)
		}()
	}

	if issue := <-issueChannel; issue != nil {
		return issue
	}

	return nil
}

func LogTerminalData(setting *LogSetting, scope string, message string, data ...map[string]any) error {
	if setting == nil {
		setting = GetLogSetting()
	}

	if setting == nil {
		return fmt.Errorf("log terminal data setting undefined")
	}

	reference, issue := ConstructReference(scope, message, data)
	if issue != nil {
		return fmt.Errorf("log data construct reference with issue, %w", issue)
	}

	if _, issue := fmt.Fprintln(os.Stdout, formatLogMessage(LogDataType, reference, scope, message)); issue != nil {
		return fmt.Errorf("write terminal data with issue, %w", issue)
	}

	return nil
}

func LogTerminalIssue(setting *LogSetting, scope string, message string, data ...map[string]any) error {
	if setting == nil {
		setting = GetLogSetting()
	}

	if setting == nil {
		return fmt.Errorf("log terminal issue setting undefined")
	}

	reference, issue := ConstructReference(scope, message, data)
	if issue != nil {
		return fmt.Errorf("log issue construct reference with issue, %w", issue)
	}

	if _, issue := fmt.Fprintln(os.Stdout, formatLogMessage(LogIssueType, reference, scope, message)); issue != nil {
		return fmt.Errorf("write terminal issue with issue, %w", issue)
	}

	return nil
}

func LogTerminalTrace(setting *LogSetting, scope string, message string, data ...map[string]any) error {
	if setting == nil {
		setting = GetLogSetting()
	}

	if setting == nil {
		return fmt.Errorf("log terminal trace setting undefined")
	}

	reference, issue := ConstructReference(scope, message, data)
	if issue != nil {
		return fmt.Errorf("log trace construct reference with issue, %w", issue)
	}

	if _, issue := fmt.Fprintln(os.Stdout, formatLogMessage(LogTraceType, reference, scope, message)); issue != nil {
		return fmt.Errorf("write terminal trace with issue, %w", issue)
	}

	return nil
}

func appendLogFile(filePath string, message string) error {
	file, issue := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if issue != nil {
		return fmt.Errorf("append log file on open with issue, %w", issue)
	}
	defer file.Close()

	if _, issue = fmt.Fprintln(file, message); issue != nil {
		return fmt.Errorf("append log file on write with issue, %w", issue)
	}

	return nil
}

func LogFileData(setting *LogSetting, scope string, message string, data ...map[string]any) error {
	if setting == nil {
		setting = GetLogSetting()
	}

	if setting == nil {
		return fmt.Errorf("log file data setting undefined")
	}

	reference, issue := ConstructReference(scope, message, data)
	if issue != nil {
		return fmt.Errorf("log file data on construct reference with issue, %w", issue)
	}

	return appendLogFile(setting.LogDataFilePath, formatLogMessage(LogDataType, reference, scope, message))
}

func LogFileIssue(setting *LogSetting, scope string, message string, data ...map[string]any) error {
	if setting == nil {
		setting = GetLogSetting()
	}

	if setting == nil {
		return fmt.Errorf("log file issue setting undefined")
	}

	reference, issue := ConstructReference(scope, message, data)
	if issue != nil {
		return fmt.Errorf("log file issue on construct reference with issue, %w", issue)
	}

	return appendLogFile(setting.LogIssueFilePath, formatLogMessage(LogIssueType, reference, scope, message))
}

func LogFileTrace(setting *LogSetting, scope string, message string, data ...map[string]any) error {
	if setting == nil {
		setting = GetLogSetting()
	}

	if setting == nil {
		return fmt.Errorf("log file trace setting undefined")
	}

	reference, issue := ConstructReference(scope, message, data)
	if issue != nil {
		return fmt.Errorf("log file trace on construct reference with issue, %w", issue)
	}

	return appendLogFile(setting.LogTraceFilePath, formatLogMessage(LogTraceType, reference, scope, message))
}

func LogIpcData(setting *LogSetting, scope string, message string, data ...map[string]any) error {
	if setting == nil {
		setting = GetLogSetting()
	}

	if setting == nil {
		return fmt.Errorf("log ipc data setting undefined")
	}

	ipcSetting := GetIpc().Setting

	logMessage := LogMessage{
		Scope:   scope,
		Message: message,
	}

	if len(data) > 0 && data[0] != nil {
		logDataBytes, issue := json.Marshal(data[0])
		if issue != nil {
			return fmt.Errorf("log ipc data on parse data with issue, %w", issue)
		}

		logMessage.Data = logDataBytes
	}

	logMessageBytes, issue := json.Marshal(logMessage)
	if issue != nil {
		return fmt.Errorf("log ipc data on parse log with issue, %w", issue)
	}

	messageBytes, issue := json.Marshal(Message{
		Category: LogDataCategory,
		Data:     logMessageBytes,
	})
	if issue != nil {
		return fmt.Errorf("log ipc data on parse message with issue, %w", issue)
	}

	sendDataPipe, issue := SendIpcData(setting.LogIpcAddress, string(messageBytes), ipcSetting).Data()
	if issue != nil {
		return fmt.Errorf("log ipc data on send with issue, %w", issue)
	}

	rawIssue := <-sendDataPipe
	if rawIssue != "" {
		return fmt.Errorf("log ipc data on send with external issue, %s", rawIssue)
	}

	return nil
}

func LogIpcIssue(setting *LogSetting, scope string, message string, data ...map[string]any) error {
	if setting == nil {
		setting = GetLogSetting()
	}

	if setting == nil {
		return fmt.Errorf("log ipc issue setting undefined")
	}

	ipcSetting := GetIpc().Setting

	logMessage := LogMessage{
		Scope:   scope,
		Message: message,
	}

	if len(data) > 0 && data[0] != nil {
		logDataBytes, issue := json.Marshal(data[0])
		if issue != nil {
			return fmt.Errorf("log ipc issue on parse data with issue, %w", issue)
		}

		logMessage.Data = logDataBytes
	}

	logMessageBytes, issue := json.Marshal(logMessage)
	if issue != nil {
		return fmt.Errorf("log ipc issue on parse log with issue, %w", issue)
	}

	messageBytes, issue := json.Marshal(Message{
		Category: LogIssueCategory,
		Data:     logMessageBytes,
	})
	if issue != nil {
		return fmt.Errorf("log ipc issue on parse message with issue, %w", issue)
	}

	sendDataPipe, issue := SendIpcData(setting.LogIpcAddress, string(messageBytes), ipcSetting).Data()
	if issue != nil {
		return fmt.Errorf("log ipc issue on send with issue, %w", issue)
	}

	rawIssue := <-sendDataPipe
	if rawIssue != "" {
		return fmt.Errorf("log ipc issue on send with external issue, %s", rawIssue)
	}

	return nil
}

func LogIpcTrace(setting *LogSetting, scope string, message string, data ...map[string]any) error {
	if setting == nil {
		setting = GetLogSetting()
	}

	if setting == nil {
		return fmt.Errorf("log ipc trace setting undefined")
	}

	ipcSetting := GetIpc().Setting

	logMessage := LogMessage{
		Scope:   scope,
		Message: message,
	}

	if len(data) > 0 && data[0] != nil {
		logDataBytes, issue := json.Marshal(data[0])
		if issue != nil {
			return fmt.Errorf("log ipc trace on parse data with issue, %w", issue)
		}

		logMessage.Data = logDataBytes
	}

	logMessageBytes, issue := json.Marshal(logMessage)
	if issue != nil {
		return fmt.Errorf("log ipc trace on parse log with issue, %w", issue)
	}

	messageBytes, issue := json.Marshal(Message{
		Category: LogTraceCategory,
		Data:     logMessageBytes,
	})
	if issue != nil {
		return fmt.Errorf("log ipc trace on parse message with issue, %w", issue)
	}

	sendDataPipe, issue := SendIpcData(setting.LogIpcAddress, string(messageBytes), ipcSetting).Data()
	if issue != nil {
		return fmt.Errorf("log ipc trace on send with issue, %w", issue)
	}

	rawIssue := <-sendDataPipe
	if rawIssue != "" {
		return fmt.Errorf("log ipc trace on send with external issue, %s", rawIssue)
	}

	return nil
}
