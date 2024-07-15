package middleware

import (
	"memrizr/account/observability/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware(c *gin.Context) {
	start := time.Now()

	fields := logrus.Fields{
		"method":   c.Request.Method,
		"path":     c.FullPath(),
		"query":    c.Request.URL.RawQuery,
		"remoteIP": c.ClientIP(),
	}

	// Log HTTP request
	//logger.LogInfo("HTTP request", fields)

	// Process the request
	c.Next()

	// Log HTTP response
	fields["status"] = c.Writer.Status()
	fields["latency"] = time.Since(start).Seconds()

	logger.LogInfo("HTTP response", fields)
}
