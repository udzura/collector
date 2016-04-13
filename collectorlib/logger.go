package collectorlib

import (
	"os"

	"github.com/Sirupsen/logrus"
)

var Logger = logrus.New()

func init() {
	Logger.Formatter = &logrus.TextFormatter{ForceColors: true}
	Logger.Out = os.Stderr
	Logger.Level = logrus.InfoLevel
}

func SwitchToVerbose() {
	Logger.Level = logrus.DebugLevel
}
