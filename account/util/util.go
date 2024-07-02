package util

import (
	"encoding/json"
	_ "fmt"
	"net/http"
	"runtime"

	"github.com/sirupsen/logrus"
)

func LogError(logger *logrus.Logger, err error, msg string) {
	pc, fn, line, _ := runtime.Caller(1)
	fields := logrus.Fields{
		"file": fn,
		"line": line,
		"func": runtime.FuncForPC(pc).Name(),
	}
	logger.WithFields(fields).WithError(err).Error(msg)
}

func LogInfo(logger *logrus.Logger, msg string) {
	pc, fn, line, _ := runtime.Caller(1)
	fields := logrus.Fields{
		"file": fn,
		"line": line,
		"func": runtime.FuncForPC(pc).Name(),
	}
	logger.WithFields(fields).Info(msg)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func LogResponse(logger *logrus.Logger, status int, method, path string) {
	logger.WithFields(logrus.Fields{
		"status": status,
		"method": method,
		"path":   path,
	}).Info("Request handled")
}
