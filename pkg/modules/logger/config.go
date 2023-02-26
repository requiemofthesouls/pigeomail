package logger

import (
	"fmt"
	"io"
	"net"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Address     string `mapstructure:"address"`
	Level       string `mapstructure:"level"`
	Encoding    string `mapstructure:"encoding"`
	Caller      bool   `mapstructure:"caller"`
	Stacktrace  string `mapstructure:"stacktrace"`
	Development bool   `mapstructure:"development"`
}

func (c Config) getZapCore() (zapcore.Core, error) {
	var eConf = zap.NewProductionEncoderConfig()
	if c.Development {
		eConf = zap.NewDevelopmentEncoderConfig()
		eConf.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var encoder = zapcore.NewConsoleEncoder(eConf)
	if c.Encoding == "json" {
		encoder = zapcore.NewJSONEncoder(eConf)
	}

	var writer = zapcore.AddSync(os.Stdout)
	if c.Address != "" {
		var (
			udpWriter io.Writer
			err       error
		)
		if udpWriter, err = newUDPWriter(c.Address); err != nil {
			return nil, err
		}

		writer = zap.CombineWriteSyncers(writer, zapcore.AddSync(udpWriter))
	}

	var level zap.AtomicLevel
	if c.Level != "" {
		if err := level.UnmarshalText([]byte(c.Level)); err != nil {
			return nil, err
		}
	}

	return zapcore.NewCore(encoder, writer, level), nil
}

func newUDPWriter(address string) (io.Writer, error) {
	var (
		remoteAddr *net.UDPAddr
		err        error
	)
	if remoteAddr, err = net.ResolveUDPAddr("udp", address); err != nil {
		return nil, err
	}

	var conn net.Conn
	if conn, err = net.DialUDP("udp", nil, remoteAddr); err != nil {
		return nil, err
	}

	return conn, nil
}

func (c Config) getZapOptions() ([]zap.Option, error) {
	var options = make([]zap.Option, 0, 5)
	if c.Caller {
		options = append(options, zap.AddCaller())
	}

	if c.Development {
		options = append(options, zap.Development())
	}

	var level zap.AtomicLevel
	if len(c.Stacktrace) > 0 {
		if err := level.UnmarshalText([]byte(c.Stacktrace)); err != nil {
			return nil, fmt.Errorf("unmarshall level stacktrace: %v", err)
		}

		options = append(options, zap.AddStacktrace(level))
	}

	return options, nil
}
