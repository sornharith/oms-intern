package logger

import (
	"github.com/sirupsen/logrus"
)

func LogError(log *logrus.Logger, err error, msg string, fields logrus.Fields) {
	fields["error"] = err.Error()
	log.WithFields(fields).Error(msg)
}

func LogInfo(log *logrus.Logger, msg string, fields logrus.Fields) {
	log.WithFields(fields).Info(msg)
}
