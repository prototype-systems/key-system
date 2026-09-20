package utilities

import (
	"keys-data-service/settings"
	"os"
	"strconv"
)

func CheckServiceFileDescriptor() bool {
	return os.Getenv(settings.FileDescriptorUseVariable) == "1" ||
		os.Getenv("USE_INHERITED_FD") == "1" ||
		os.Getenv("SOCKET_FD") != ""
}

func GetServiceFileDescriptor() uintptr {
	if fileDescriptor := os.Getenv("SOCKET_FD"); fileDescriptor != "" {
		if descriptorValue, issue := strconv.Atoi(fileDescriptor); issue == nil && descriptorValue >= 0 {
			return uintptr(descriptorValue)
		}
	}

	return uintptr(3)
}

func CheckChild() bool {
	return len(os.Args) > 1 && os.Args[1] == "child"
}

func GetExecutablePath() string {
	if exe, err := os.Executable(); err == nil {
		return exe
	}

	return os.Args[0]
}
