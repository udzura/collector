package cmd

import (
	"os"

	"github.com/Sirupsen/logrus"
)

var logger = logrus.New()

func init() {
	logrus.SetOutput(os.Stderr)
	logrus.SetLevel(logrus.InfoLevel)
}
