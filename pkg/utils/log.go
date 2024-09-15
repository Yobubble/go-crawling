package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func LogInit() {
	Log = logrus.New()
	Log.SetFormatter(&logrus.TextFormatter{})
	Log.Level = logrus.TraceLevel
	Log.Out = os.Stdout
}