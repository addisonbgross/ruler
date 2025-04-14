package util

import (
	"fmt"
	sh "ruler/node/shared"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugarLogger *zap.SugaredLogger

func GetLogger() (*zap.SugaredLogger, error) {
	if sugarLogger != nil {
		return sugarLogger, nil
	}

	config := zap.NewProductionConfig()
	config.InitialFields = map[string]interface{}{
		"node": sh.NodeID,
	}
	config.EncoderConfig.TimeKey = zapcore.OmitKey
	config.EncoderConfig.CallerKey = zapcore.OmitKey
	config.EncoderConfig.StacktraceKey = zapcore.OmitKey

	logger, err := config.Build()
	if err != nil {
		fmt.Printf("%+v", err)
	}
	sugarLogger = logger.Sugar()
	return sugarLogger, nil
}
