package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type Logger struct {
	*zap.SugaredLogger
}

// New creates a new logger ready for production using the default configuration.
func NewProduction() *Logger {
	l, _ := zap.NewProduction()
	return NewWithZap(l)
}

// New creates a new logger ready for development using the default configuration.
func NewDevelopment() *Logger {
	l, _ := zap.NewDevelopment()
	return NewWithZap(l)
}

// NewWithZap creates a new logger using the preconfigured zap logger.
func NewWithZap(l *zap.Logger) *Logger {
	return &Logger{l.Sugar()}
}

// NewForTest returns a new logger and the corresponding observed logs which can be used in unit tests to verify log entries.
func NewForTest() (*Logger, *observer.ObservedLogs) {
	core, recorded := observer.New(zapcore.InfoLevel)
	core = zapcore.NewTee(
		core,
		zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), zapcore.Lock(os.Stdout), zapcore.DebugLevel),
	)
	return NewWithZap(zap.New(core)), recorded
}
