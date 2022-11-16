package logger

import (
	"lostinsoba/ninhydrin/internal/monitoring/logger/zap"
)

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Warn(args ...interface{})
}

func NewLogger(level string, labels map[string]string) Logger {
	return zap.NewLogger(level, labels)
}
