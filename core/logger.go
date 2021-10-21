package core

import (
	"path"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger interafce
type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
}

// ZapLogger structure
type ZapLogger struct {
	Zap *zap.SugaredLogger
}

// NewLogger set up logger
func NewLogger(config Config) Logger {

	zapConfig := zap.NewDevelopmentConfig()

	if config.App.Debug {
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	if config.Log.Server != "" {
		zapConfig.OutputPaths = []string{
			path.Join(config.File.Current, config.Log.Server),
		}
	} else {
		zapConfig.OutputPaths = []string{
			path.Join(config.File.Current, "server.log"),
		}
	}

	logger, _ := zapConfig.Build()

	sugar := logger.Sugar()

	return &ZapLogger{
		Zap: sugar,
	}
}

func (log *ZapLogger) Debug(args ...interface{}) {
	log.Zap.Debug(args)
}

func (log *ZapLogger) Info(args ...interface{}) {
	log.Zap.Info(args)
}

func (log *ZapLogger) Warn(args ...interface{}) {
	log.Zap.Warn(args)
}

func (log *ZapLogger) Error(args ...interface{}) {
	log.Zap.Error(args)
}
