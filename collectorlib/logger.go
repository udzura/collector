package collectorlib

import (
	"os"

	"github.com/Sirupsen/logrus"
)

var Logger = logrus.New()

func init() {
	logrus.SetOutput(os.Stderr)
	logrus.SetLevel(logrus.InfoLevel)
}
