package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/geektheripper/go-gutils/git/git_utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logger = log.New(os.Stderr)

var rootCmd = &cobra.Command{
	Use:   "mygo",
	Short: "MyGO: A Lifetime Golang Monorepo Manager",
	Long:  "MyGO is a tool helps you managing golang sub packages.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}

func init() {
	viper.SetEnvPrefix("MYGO")

	defaultRepo, _ := git_utils.FindGitRoot(".")
	rootCmd.PersistentFlags().StringP("repo", "r", defaultRepo, "the remote repository to manage")
	viper.BindEnv("repo")

	viper.BindPFlags(rootCmd.PersistentFlags())
}
