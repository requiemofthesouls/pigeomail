package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(cfg Config, fields []Field) (Wrapper, error) {
	var (
		core zapcore.Core
		err  error
	)
	if core, err = cfg.getZapCore(); err != nil {
		return nil, fmt.Errorf("get zapcore: %v", err)
	}

	var options []zap.Option
	if options, err = cfg.getZapOptions(); err != nil {
		return nil, fmt.Errorf("get options: %v", err)
	}

	options = append(
		options,
		zap.Fields(fields...),
		zap.AddCallerSkip(1),
	)

	return &wrapper{
		zap: zap.New(zapcore.NewTee(core), options...),
	}, nil
}

func NewFromZap(l *zap.Logger) Wrapper {
	return &wrapper{
		zap: l,
	}
}

type (
	Wrapper interface {
		GetLogger() *zap.Logger
		With(fields ...Field) Wrapper
		Debug(msg string, fields ...Field)
		Info(msg string, fields ...Field)
		Warn(msg string, fields ...Field)
		Error(msg string, fields ...Field)
		Sync() error
	}

	wrapper struct {
		zap *zap.Logger
	}
)

func (w *wrapper) GetLogger() *zap.Logger {
	return w.zap
}

func (w *wrapper) With(fields ...Field) Wrapper {
	return &wrapper{
		zap: w.zap.With(fields...),
	}
}

func (w *wrapper) Debug(msg string, fields ...Field) {
	w.zap.Debug(msg, fields...)
}

func (w *wrapper) Info(msg string, fields ...Field) {
	w.zap.Info(msg, fields...)
}

func (w *wrapper) Warn(msg string, fields ...Field) {
	w.zap.Warn(msg, fields...)
}

func (w *wrapper) Error(msg string, fields ...Field) {
	w.zap.Error(msg, fields...)
}

func (w *wrapper) Sync() error {
	return w.zap.Sync()
}
