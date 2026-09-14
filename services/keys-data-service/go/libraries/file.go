package libraries

import (
	"os"
	"path/filepath"
)

func CheckFile(filePath string) bool {
	fileMetadata, issue := os.Stat(filePath)
	return issue == nil && !fileMetadata.IsDir()
}

func CheckDirectory(directoryPath string) bool {
	fileMetadata, issue := os.Stat(directoryPath)
	return issue == nil && fileMetadata.IsDir()
}

func ResolveDirectory(directoryPath string) (string, error) {
	if directoryPath == "" {
		return "", nil
	}

	absolutePath, issue := filepath.Abs(directoryPath)
	if issue != nil {
		return "", issue
	}

	if !CheckDirectory(absolutePath) {
		if issue = os.MkdirAll(absolutePath, os.ModePerm); issue != nil {
			return "", issue
		}
	}

	return absolutePath, nil
}

func ResolveFile(filePath string) (string, error) {
	if filePath == "" {
		return "", nil
	}

	absolutePath, issue := filepath.Abs(filePath)
	if issue != nil {
		return "", issue
	}

	if !CheckFile(absolutePath) {
		if issue = os.MkdirAll(filepath.Dir(absolutePath), os.ModePerm); issue != nil {
			return "", issue
		}

		file, issue := os.Create(absolutePath)
		if issue != nil {
			return "", issue
		}

		file.Close()
	}

	return absolutePath, nil
}
