package utils

import (
	"path/filepath"

	"github.com/geektheripper/mygo/internal/log"
	"github.com/spf13/viper"
)

var logger = log.GetLogger()

func MustGetRepo() string {
	repo := viper.GetString("repo")
	if repo == "" {
		logger.Fatalf("repo is not set")
	}

	if !IsGitRepo(repo) {
		logger.Fatalf("repo is not a git repository: %s", repo)
	}

	return repo
}

func MustGetPackageNamePath(packageName string) (string, string) {
	if packageName == "" {
		logger.Fatalf("package name is required")
	}

	if filepath.IsAbs(packageName) {
		logger.Fatalf("package name must be a relative path: %s", packageName)
	}

	return packageName, filepath.Join(MustGetRepo(), packageName)
}
