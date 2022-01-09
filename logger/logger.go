package logger

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func buildZapLogger() *zap.Logger {
	zc := zap.NewProductionConfig()
	zc.Level = zap.NewAtomicLevelAt(zapcore.Level(-10))
	zc.OutputPaths = []string{"stdout"}
	zc.ErrorOutputPaths = zc.OutputPaths
	zc.DisableStacktrace = true
	zc.DisableCaller = true
	z, _ := zc.Build()
	return z
}

func New() logr.Logger {
	log := zapr.NewLogger(buildZapLogger())
	return log
}
