package utils

import (
	"encoding/json"
	"go.uber.org/zap"
	"io/ioutil"
)

var Logger *zap.SugaredLogger

// InitializeLogger creates de Debug logger following the configuration details file. If ignore is used, then all
// non-errors messages will be ignored.
func InitializeLogger(configFile string, ignore bool) {
	var cfg zap.Config

	if ignore {
		cfg = zap.NewProductionConfig()
		cfg.OutputPaths = []string{}
	} else {
		content, err := ioutil.ReadFile(configFile)
		if err != nil {
			panic("Could not read logging configuration.")
		}

		if err := json.Unmarshal(content, &cfg); err != nil {
			panic("Logging configuration is not valid.")
		}
	}

	log, err := cfg.Build()
	if err != nil {
		panic(err)
	} else {
		Logger = log.Sugar()
		Logger.Info("Program Started =============")
	}
}
