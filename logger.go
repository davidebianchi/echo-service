package main

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func setupLogger() (*zap.Logger, error) {
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevel(),
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
