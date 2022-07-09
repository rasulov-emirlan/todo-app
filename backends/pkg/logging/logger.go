package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
	closer func() error
}

type ErrCloser struct {
	Errors []error
}

func (e *ErrCloser) Error() string {
	err := ""
	for i, v := range e.Errors {
		err += v.Error()
		if i != len(e.Errors)-1 {
			err += ","
		}
	}
	return err
}

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = "debug"
	// InfoLevel is the default logging priority.
	InfoLevel = "info"
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = "warn"
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = "error"
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel = "dpanic"
	// PanicLevel logs a message, then panics.
	PanicLevel = "panic"
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = "fatal"
)

type Field = zap.Field

func String(key, data string) Field {
	return zap.String(key, data)
}

func NewLogger(level string, output string) (*Logger, error) {
	var (
		logger Logger
		err    error
	)
	config := zap.NewProductionConfig()

	if output != "stdout" {
		f, err := os.OpenFile(output, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
		if err != nil {
			return nil, err
		}
		logger.closer = f.Close
	}

	config.Encoding = "json"
	lvl := zapcore.DebugLevel
	switch level {
	case InfoLevel:
		lvl = zapcore.InfoLevel
	case WarnLevel:
		lvl = zapcore.WarnLevel
	case PanicLevel:
		lvl = zapcore.PanicLevel
	case FatalLevel:
		lvl = zapcore.FatalLevel
	case DPanicLevel:
		lvl = zap.DPanicLevel
	}

	config.Level = zap.NewAtomicLevelAt(lvl)

	logger.logger, err = config.Build()
	if err != nil {
		if errr := logger.closer(); errr != nil {
			return nil, &ErrCloser{
				Errors: []error{errr, err},
			}
		}
		return nil, err
	}
	return &logger, nil
}

func (l *Logger) Close() error {
	return l.closer()
}

func (l *Logger) Sync() error {
	return l.logger.Sync()
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...Field) {
	l.logger.Panic(msg, fields...)
}

func (l *Logger) DPanic(msg string, fields ...Field) {
	l.logger.DPanic(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, fields...)
}
