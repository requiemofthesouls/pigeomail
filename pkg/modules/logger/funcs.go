package logger

import (
	"time"

	"go.uber.org/zap"
)

func String(key string, val string) Field {
	return zap.String(key, val)
}

func Int64(key string, val int64) Field {
	return zap.Int64(key, val)
}

func Int32(key string, val int32) Field {
	return zap.Int32(key, val)
}

func Time(key string, val time.Time) Field {
	return zap.Time(key, val)
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
