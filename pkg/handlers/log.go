package handlers

import "github.com/sirupsen/logrus"

var (
	logger *logrus.Logger
)

//SetLogger used to pass in the logger object from the trackingbot package to avoid circular dependency
func SetLogger(l *logrus.Logger) *logrus.Logger {
	logger = l
	return logger
}
