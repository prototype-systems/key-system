package utilities

import (
	"keys-data-service-test/settings"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetRootDirectory() (string, error) {
	output, issue := exec.Command("git", "rev-parse", "--show-toplevel").CombinedOutput()
	if issue != nil {
		return "", issue
	}

	return filepath.Clean(strings.TrimSpace(string(output))), nil
}

func ResolveTestDirectory() (string, error) {
	homeDirectory, issue := os.UserHomeDir()
	if issue != nil {
		return "", issue
	}

	testDirectory := filepath.Join(homeDirectory, settings.RootTestDirectory, settings.TargetNamespace)
	if issue := os.RemoveAll(testDirectory); issue != nil {
		return "", issue
	}

	if issue := os.MkdirAll(testDirectory, 0755); issue != nil {
		return "", issue
	}

	return testDirectory, nil
}

func DropTestDirectory(testDirectory string) error {
	if issue := os.RemoveAll(testDirectory); issue != nil {
		return issue
	}

	return nil
}

func ResolveTargetDirectory() (string, error) {
	repositoryRoot, issue := GetRootDirectory()
	if issue != nil {
		return "", issue
	}

	targetDirectory := filepath.Join(repositoryRoot, settings.TargetTestDirectory)

	return targetDirectory, nil
}
