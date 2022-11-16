package zap

import (
	zapLogger "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const logLevelDefault = zapcore.DebugLevel

func NewLogger(level string, labels map[string]string) *zapLogger.SugaredLogger {
	logLevel, err := zapcore.ParseLevel(level)
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
