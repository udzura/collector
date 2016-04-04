package cmd

import (
	// "fmt"
	"errors"
	"os"
	"os/exec"
	"syscall"

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
		ret := runCheck(args)
		os.Exit(ret)
	},
}

func init() {
	RootCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringVarP(&deviceName, "dev", "D", "eth0", "Device name which has the instance's global IP")
}

func runCheck(args []string) int {
	var exitStatus int
	if len(args) < 1 {
		logger.Errorln("Check command is empty.")
		return -1
	}

	check := exec.Command(args[0], args[1:]...)
	_, err := check.Output()
	if err != nil {
		if e2, ok := err.(*exec.ExitError); ok {
			if s, ok := e2.Sys().(syscall.WaitStatus); ok {
				exitStatus = s.ExitStatus()
			} else {
				panic(errors.New("Unimplemented for system where exec.ExitError.Sys() is not syscall.WaitStatus."))
			}
		}
	} else {
		exitStatus = 0
	}

	logger.Infof("client called: %v", args)
	logger.Infof("command exit with %d", exitStatus)
	return exitStatus
}
