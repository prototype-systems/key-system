package utilities

import (
	"keys-data-service-test/go/settings"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func ConstructCommand(command string, arguments ...string) (*exec.Cmd, error) {
	targetDirectory, issue := ResolveTargetDirectory()
	if issue != nil {
		return nil, issue
	}

	process := exec.Command(command, arguments...)
	process.Dir = targetDirectory
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr
	process.Env = ResolveEnvironment()

	return process, nil
}

func KillProcess(reference int) error {
	if CheckWindows() {
		return exec.Command("taskkill", "/PID", strconv.Itoa(reference), "/T", "/F").Run()
	}

	process, issue := os.FindProcess(reference)
	if issue != nil {
		return issue
	}

	if issue := process.Kill(); issue != nil {
		return issue
	}

	return nil
}

func CheckCommandExit(command *exec.Cmd) bool {
	done := make(chan error, 1)
	go func() { done <- command.Wait() }()

	select {
	case issue := <-done:
		if issue != nil {
			return false
		}

		return true

	case <-time.After(settings.TestCommandTimeout):
		return false
	}
}
