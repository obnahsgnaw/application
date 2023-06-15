package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const EncodingJson EncodingType = "json"
const EncodingConsole EncodingType = "console"

type EncodingType string

func (t EncodingType) String() string {
	return string(t)
}

// NewEncoderConfig encoder config
func NewEncoderConfig() zapcore.EncoderConfig {
	encoderConf := zap.NewProductionEncoderConfig()
	encoderConf.EncodeTime = zapcore.ISO8601TimeEncoder

	return encoderConf
}

// NewLoggerConfig logger config
func NewLoggerConfig(level zapcore.Level, output []string, errPath []string, develop bool, encoding EncodingType, encoderConf zapcore.EncoderConfig) zap.Config {
	conf := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: develop,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         encoding.String(),
		EncoderConfig:    encoderConf,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	if len(output) > 0 {
		conf.OutputPaths = output
	}
	if len(errPath) > 0 {
		conf.ErrorOutputPaths = errPath
	}

	return conf
}

// New a new logger
func New(name string, conf zap.Config, opts ...zap.Option) (*zap.Logger, error) {
	logger, err := conf.Build(opts...)
	if err != nil {
		return nil, err
	}

	if name != "" {
		logger = logger.Named(name)
	}

	return logger, nil
}

// NewJsonLogger json encoder logger
func NewJsonLogger(name string, level zapcore.Level, output []string, errPath []string, develop bool) (*zap.Logger, error) {
	return New(name, NewLoggerConfig(level, output, errPath, develop, EncodingJson, NewEncoderConfig()))
}

// NewConsoleLogger console encoder logger
func NewConsoleLogger(name string, level zapcore.Level, output []string, errPath []string, develop bool) (*zap.Logger, error) {
	return New(name, NewLoggerConfig(level, output, errPath, develop, EncodingConsole, NewEncoderConfig()))
}

// NewCliLogger cli console logger
func NewCliLogger(name string, level zapcore.Level, develop bool) (*zap.Logger, error) {
	encoderConf := NewEncoderConfig()
	encoderConf.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return New(name, NewLoggerConfig(level, nil, nil, develop, EncodingConsole, encoderConf))
}

// NewMultiLogger merge multi logger in logger
func NewMultiLogger(l *zap.Logger, ol ...*zap.Logger) *zap.Logger {
	if len(ol) > 0 {
		var cores []zapcore.Core
		cores = append(cores, l.Core())
		for _, ll := range ol {
			cores = append(cores, ll.Core())
		}
		l = l.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(cores...)
		}))
	}

	return l
}
