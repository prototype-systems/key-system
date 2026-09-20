package utilities

import "runtime"

func CheckWindows() bool {
	return runtime.GOOS == "windows"
}

func CheckLinux() bool {
	return runtime.GOOS == "linux"
}

func CheckDarwin() bool {
	return runtime.GOOS == "darwin"
}
