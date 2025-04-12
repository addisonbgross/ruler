package util

import (
	"errors"
	t "ruler-node/internal/types"

	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugarLogger *zap.SugaredLogger
var info *t.NodeInfo

func InitLogger(n *t.NodeInfo) {
	info = n
}

func GetLogger() (*zap.SugaredLogger, error) {
	if info == nil {
		return nil, errors.New("logger not supplied with NodeInfo")
	}

	if sugarLogger != nil {
		return sugarLogger, nil
	}

	config := zap.NewProductionConfig()
	config.InitialFields = map[string]interface{}{
		"node": fmt.Sprintf("%s:%s", info.Ip, info.Port),
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
