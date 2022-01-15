package main

import (
	"github.com/sirupsen/logrus"
)

func logInfo(args ...interface{}) {
	logrus.Info(args)
}

func logWarn(args ...interface{}) {
	logrus.Warn(args)
}

func logError(args ...interface{}) {
	logrus.Error(args)
}

func logFatal(args ...interface{}) {
	logrus.Fatal(args)
}
