package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "collector's client command",
	Long: `collector's client command.
This subcommand ins intendec to run under consul checker.`,
	Run: func(cmd *cobra.Command, args []string) {

		// TODO: Work your own magic here
		fmt.Println("client called")
	},
}

func init() {
	RootCmd.AddCommand(clientCmd)
}
