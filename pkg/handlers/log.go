package handlers

import "github.com/sirupsen/logrus"

var (
	logger *logrus.Logger
)

func SetLogger(l *logrus.Logger) *logrus.Logger {
	logger = l
	return logger
}
