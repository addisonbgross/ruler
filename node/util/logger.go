package util

import (
	t "Git/ruler/node/types"

	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugarLogger *zap.SugaredLogger
var info *t.NodeInfo

func InitLogger(n *t.NodeInfo) {
	info = n
}

func GetLogger() *zap.SugaredLogger {
	if info == nil {
		panic("Logger not supplied with NodeInfo")
	}

	if sugarLogger != nil {
		return sugarLogger
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
	return sugarLogger
}
