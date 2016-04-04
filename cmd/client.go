package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

var deviceName string

var ipFinder = regexp.MustCompile(`inet \d+\.\d+\.\d+\.\d+`)

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
	checkOut, err := check.Output()
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

	var sign string
	if exitStatus == 0 {
		sign = "OK"
	} else {
		sign = "NG"
	}

	ipcmd := exec.Command("/sbin/ip", "address", "show", deviceName)
	ipOut, err := ipcmd.Output()
	if err != nil {
		logger.Errorf("Somthing is wrong with getting ip: %s", err.Error())
		return -1
	}

	res := ipFinder.Find(ipOut)
	if len(res) == 0 {
		logger.Errorf("Somthing is wrong with getting ip: response is empty")
		return -1
	}
	gip := strings.TrimPrefix(string(res), "inet ")

	rep := strings.NewReplacer("\n", " ", "\t", " ")
	checkOutEscaped := rep.Replace(string(checkOut))
	fmt.Fprintf(os.Stdout, "status:%s\tcode:%d\tcommand_out:%s\tipaddr:%s",
		sign, exitStatus, checkOutEscaped, gip)

	logger.Infof("client called: %v", args)
	return exitStatus
}
