package zap

import (
	zapLogger "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const Kind = "zap"

const logLevelDefault = zapcore.DebugLevel

const (
	settingLevel = "level"
)

func NewLogger(settings map[string]string, labels map[string]string) *zapLogger.SugaredLogger {
	levelStr := settings[settingLevel]
	logLevel, err := zapcore.ParseLevel(levelStr)
	if err != nil {
		logLevel = logLevelDefault
	}

	fields := make([]zapLogger.Field, 0, len(labels))
	for label, labelValue := range labels {
		fields = append(fields, zapLogger.String(label, labelValue))
	}
	zapLoggerFields := zapLogger.Fields(fields...)

	config := zapLogger.NewProductionConfig()
	config.DisableStacktrace = true
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.Level.SetLevel(logLevel)

	log, _ := config.Build(zapLoggerFields)
	return log.Sugar()
}
