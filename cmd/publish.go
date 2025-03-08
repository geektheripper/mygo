package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/geektheripper/mygo/cmd/utils"
	"github.com/geektheripper/mygo/pkg/local_repo"
	"github.com/geektheripper/mygo/pkg/virtual_repo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dustin/go-humanize"
)

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish a package to the remote repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoPath := utils.MustGetRepo()
		packageName, packagePath := utils.MustGetPackageNamePath(args[0])
		includes := viper.GetStringSlice("includes")
		message := viper.GetString("message")
		version := viper.GetString("version")

		upgradeType := ""
		if viper.GetBool("minor") {
			upgradeType = "minor"
		} else if viper.GetBool("major") {
			upgradeType = "major"
		} else if viper.GetBool("patch") {
			upgradeType = "patch"
		}

		lrepo, err := local_repo.LoadLocalRepo(repoPath)
		if err != nil {
			logger.Fatalf("failed to load local repo: %v", err)
		}

		remote := viper.GetString("remote")
		if !strings.Contains(remote, ":") {
			remote = lrepo.GetRemoteURL(remote)
		}

		if !utils.ValidateGitRemoteURL(remote) {
			logger.Fatalf("invalid remote: %s", err)
		}

		vrepo, err := virtual_repo.NewVirtualRepo(remote, packageName)
		defer vrepo.Close()

		if err != nil {
			logger.Fatalf("failed to create virtual repo: %v", err)
		}

		packageMap, err := vrepo.ListPackages()
		if err != nil {
			logger.Fatalf("failed to fetch packages form remote: %v", err)
		}

		pkg, ok := packageMap[packageName]
		if !ok {
			if upgradeType != "" {
				logger.Fatalf("failed to apply %s, package not found in remote", upgradeType)
			}
			version = "0.0.1"
		} else {
			version = pkg.NextVersion(upgradeType).String()
			logger.Printf("next version: %s", version)
		}

		if !viper.GetBool("no-license") {
			_, err := vrepo.Import(
				filepath.Join(repoPath, "LICENSE"),
				filepath.Join(packageName, "LICENSE"),
			)
			if err != nil {
				logger.Printf("warning: failed to copy license: %v", err)
			}
		}

		for _, file := range includes {
			_, err := vrepo.Import(
				filepath.Join(repoPath, file),
				filepath.Join(packageName, file),
			)
			if err != nil {
				logger.Fatalf("failed to copy file: %v", err)
			}
		}

		report, err := vrepo.Import(packagePath, packageName)

		if err != nil {
			logger.Fatalf("failed to copy package files: %v", err)
		}

		logger.Printf("imported %d files, %s", report.Count, humanize.Bytes(report.Size))

		tag := fmt.Sprintf("%s/v%s", packageName, version)
		if err := vrepo.Publish(tag, message); err != nil {
			logger.Fatalf("failed to publish: %v", err)
		}

		logger.Printf("published %s", tag)
	},
}

func init() {
	publishCmd.Flags().String("remote", "origin", "the remote to push to")
	viper.BindEnv("remote")

	publishCmd.Flags().StringP("version", "v", "", "publish a specific version")
	publishCmd.Flags().Bool("minor", false, "publish a minor version")
	publishCmd.Flags().Bool("major", false, "publish a major version")
	publishCmd.Flags().Bool("patch", false, "publish a patch version")
	publishCmd.MarkFlagsMutuallyExclusive("version", "minor", "major", "patch")

	publishCmd.Flags().StringP("message", "m", "Publish", "the commit message")

	publishCmd.Flags().Bool("no-license", false, "default copy license from root, use this to skip")
	publishCmd.Flags().StringArray("includes", []string{}, "include files to the package from root")

	viper.BindPFlags(publishCmd.Flags())

	rootCmd.AddCommand(publishCmd)
}
