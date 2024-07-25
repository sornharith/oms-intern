package logger

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

var (
	LogFile *os.File
	Log     *logrus.Logger
)

func Setup() {
	// var err error
	// LogFile, err = os.OpenFile("./observability/logger/temp/logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err != nil {
	// 	logrus.Fatalf("Failed to open log file: %v", err)
	// }
	Log = logrus.New()
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetLevel(logrus.InfoLevel)
}

func GinLogger(logger *logrus.Logger) gin.HandlerFunc {
	if logger == nil {
		logrus.Error("Logger is not initialized")
		return func(c *gin.Context) {} // Return no-op handler
	}
	return ginlogrus.Logger(logger)
}

func GinRecovery(logger *logrus.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok && logger != nil {
			logger.WithFields(logrus.Fields{
				"module": "gin",
				"error":  err,
			}).Error("Panic recovered")
		} else {
			logrus.Error("Panic recovered but logger is nil")
		}
		c.AbortWithStatus(500)
	})
}

func LogError(err error, msg string, fields logrus.Fields) {
	fields["error"] = err.Error()
	Log.WithFields(fields).Error(msg)
}

func LogInfo(msg string, fields logrus.Fields) {
	Log.WithFields(fields).Info(msg)
}
