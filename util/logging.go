package utils

import "go.uber.org/zap"

var Logger *zap.SugaredLogger

func InitializeLogger() {
	if log, err := zap.NewDevelopment(); err != nil {
		panic("Logger cannot be initialized.")
	} else {
		Logger = log.Sugar()
	}
}
