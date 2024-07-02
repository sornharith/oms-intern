package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

type key int

const loggerKey key = iota

func WithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) *logrus.Entry {
	logger, ok := ctx.Value(loggerKey).(*logrus.Entry)
	if !ok {
		logger = logrus.NewEntry(logrus.New())
	}
	return logger
}

func Error(ctx context.Context, err error, msg string, fields logrus.Fields) {
	logger := FromContext(ctx).WithFields(fields)
	logger.WithError(err).Error(msg)
}

func Info(ctx context.Context, msg string, fields logrus.Fields) {
	logger := FromContext(ctx).WithFields(fields)
	logger.Info(msg)
}
