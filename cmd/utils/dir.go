package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func IsGitRepo(dir string) bool {
	if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
		return true
	}

	return false
}

func IsRootDir(dir string) bool {
	return filepath.Dir(dir) == dir
}

var projectRoot string

func GetProjectRoot() (string, error) {
	if projectRoot != "" {
		return projectRoot, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if IsRootDir(wd) {
			return "", fmt.Errorf("not in a git repository")
		}

		if IsGitRepo(wd) {
			projectRoot = wd
			return wd, nil
		}

		wd = filepath.Dir(wd)
	}
}

func MustGetProjectRoot() string {
	root, err := GetProjectRoot()
	if err != nil {
		logger.Fatal(err)
	}

	return root
}
