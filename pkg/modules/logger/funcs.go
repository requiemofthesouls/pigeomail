package logger

import (
	"go.uber.org/zap"
)

func String(key, val string) Field {
	return zap.String(key, val)
}

func Error(err error) Field {
	return zap.Error(err)
}

func ByteString(key string, val []byte) Field {
	return zap.ByteString(key, val)
}

func Reflect(key string, val interface{}) Field {
	return zap.Reflect(key, val)
}
