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
	"github.com/udzura/collector/collectorlib"
	"net/http"
	"io/ioutil"
)

var deviceName string
var url string

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
	clientCmd.Flags().StringVarP(&url, "url", "U", "", "Get instance's global IP from URL")
}

func runCheck(args []string) int {
	collectorlib.Logger.Debugf("client called: %v", args)

	var exitStatus int
	if len(args) < 1 {
		collectorlib.Logger.Errorln("Check command is empty.")
		return -1
	}

	check := exec.Command(args[0], args[1:]...)
	checkOut, err := check.CombinedOutput()
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

	gip, err := getGIP(url, deviceName)
	if err != nil {
		collectorlib.Logger.Errorf("Somthing is wrong with getting ip: %s", err.Error())
		return -1
	}

	rep := strings.NewReplacer("\n", " ", "\t", " ")
	checkOutEscaped := rep.Replace(string(checkOut))
	fmt.Fprintf(os.Stdout, "status:%s\tcode:%d\tcommand_out:%s\tipaddr:%s\n",
		sign, exitStatus, checkOutEscaped, gip)

	return exitStatus
}

func getGIP(url string, deviceName string) (gip string, err error) {
	if url != "" {
		gip, err = getGIPFromURL(url)
		if err != nil {
			return "", err
		}
	} else {
		gip, err = getGIPFromDevice(deviceName)
		if err != nil {
			return "", err
		}
	}
	return gip, nil
}

func getGIPFromURL(url string) (gip string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func getGIPFromDevice(deviceName string) (gip string, err error) {
	ipcmd := exec.Command("/sbin/ip", "address", "show", deviceName)
	ipOut, err := ipcmd.Output()
	if err != nil {
		return "", err
	}

	res := ipFinder.Find(ipOut)
	if len(res) == 0 {
		return "", errors.New("response is empty")
	}
	gip = strings.TrimPrefix(string(res), "inet ")
	return gip, nil
}
