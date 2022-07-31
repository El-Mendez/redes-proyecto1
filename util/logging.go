package utils

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func InitializeLogger() {
	log, err := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{},
		ErrorOutputPaths: []string{},
	}.Build()

	if err != nil {
		panic("Logger cannot be initialized.")
	} else {
		Logger = log.Sugar()
	}
	//if log, err := zap.NewDevelopment(); err != nil {
	//	panic("Logger cannot be initialized.")
	//} else {
	//	Logger = log.Sugar()
	//}
}
