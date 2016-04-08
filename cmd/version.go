package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "0.3.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Long:  `Show version.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version: %s\n", Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
