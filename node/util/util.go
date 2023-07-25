package util

import (
	"fmt"

	"go.uber.org/zap"
)

var sugarLogger *zap.SugaredLogger

func GetLogger() *zap.SugaredLogger {
	if sugarLogger != nil {
		return sugarLogger
	}

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("%+v", err)
	}
	sugarLogger = logger.Sugar()
	return sugarLogger
}
