package util

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var sugarLogger *zap.SugaredLogger

func GetLogger() (*zap.SugaredLogger, error) {
	if sugarLogger != nil {
		return sugarLogger, nil
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	config := zap.NewProductionConfig()
	config.InitialFields = map[string]interface{}{
		"node": hostname,
	}
	config.EncoderConfig.TimeKey = zapcore.OmitKey
	config.EncoderConfig.CallerKey = zapcore.OmitKey
	config.EncoderConfig.StacktraceKey = zapcore.OmitKey

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	sugarLogger = logger.Sugar()
	return sugarLogger, nil
}
