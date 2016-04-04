package cmd

import (
	"github.com/spf13/cobra"
)

var deviceName string

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "collector's client command",
	Long: `collector's client command.
This subcommand ins intendec to run under consul checker.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Infof("client called: %v", args)
	},
}

func init() {
	RootCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringVarP(&deviceName, "dev", "D", "eth0", "Device name which has the instance's global IP")
}
