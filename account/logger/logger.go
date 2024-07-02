package logger

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

func Setup() {
	LogFile, err := os.OpenFile("./logger/temp/logs.log", os.O_CREATE|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open log file; %v", err)
	}
	log := logrus.New()
	log.SetOutput(LogFile)
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)

}

func GinLogger(logger *logrus.Logger) gin.HandlerFunc {
	return ginlogrus.Logger(logger)
}

func GinRecovery(logger *logrus.Logger) gin.HandlerFunc {
	return gin.Recovery()
}
