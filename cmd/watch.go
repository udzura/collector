package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "The IP watcher",
	Long: `The IP watcher.
This subcommand is intended to run under "consul watch".`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("watch called")
	},
}

func init() {
	RootCmd.AddCommand(watchCmd)
	// watchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
