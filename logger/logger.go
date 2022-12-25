package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugaredLogger *zap.SugaredLogger
var loglevel *zap.AtomicLevel

func init() {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	loglevel = &cfg.Level

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	sugaredLogger = logger.Sugar()
}

func Get() *zap.SugaredLogger {
	return sugaredLogger
}

func EnableDebugLog() {
	loglevel.SetLevel(zap.DebugLevel)
}
