package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "create a new sub package",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		packageName, packagePath := MustGetPackageNamePath(args[0])

		_, err := os.Stat(packagePath)
		if err == nil {
			logger.Fatalf("package already exists: %s", packagePath)
		}

		if !os.IsNotExist(err) {
			logger.Fatalf("failed to check package: %s", err)
		}

		err = os.MkdirAll(packagePath, 0755)
		if err != nil {
			logger.Fatalf("failed to create package: %s", err)
		}

		logger.Printf("created package: %s", packageName)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
