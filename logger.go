package main

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func setupLogger(logLevel string) (*zap.Logger, error) {
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(getZapCoreLevel(logLevel)),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			EncodeLevel: numberLevelEncoder,
			LevelKey:    "level",
			MessageKey:  "msg",
			TimeKey:     "time",
			EncodeTime:  epochMillisTimeEncoder,
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
}

func numberLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt(getLogLevel(l))
}

func getLogLevel(l zapcore.Level) int {
	switch l {
	case zapcore.ErrorLevel:
		return 50
	case zapcore.InfoLevel:
		return 30
	case zapcore.DebugLevel:
		return 20
	}
	return 30
}

func epochMillisTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	nanos := t.UnixNano()
	millis := nanos / int64(time.Millisecond)
	enc.AppendInt64(millis)
}

func getZapCoreLevel(logLevel string) zapcore.Level {
	switch logLevel {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	}
	return zapcore.InfoLevel
}
