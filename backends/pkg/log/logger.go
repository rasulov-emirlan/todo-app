package log

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(level zapcore.Level, output string) *zap.Logger {
	config := zap.NewProductionConfig()
	config.Encoding = "json"
	config.Level = zap.NewAtomicLevelAt(level)
	logger, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}
	return logger
}

func ParseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}
