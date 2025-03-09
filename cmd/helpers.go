package cmd

import (
	"path/filepath"

	"github.com/geektheripper/go-gutils/git/git_utils"
	"github.com/spf13/viper"
)

func MustGetRepo() string {
	repo := viper.GetString("repo")
	if repo == "" {
		logger.Fatalf("repo is not set")
	}

	ok, err := git_utils.IsGitRepo(repo)
	if err != nil {
		logger.Fatalf("failed to check if repo is a git repository: %v", err)
	}

	if !ok {
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

	return filepath.Clean(packageName), filepath.Clean(filepath.Join(MustGetRepo(), packageName))
}
