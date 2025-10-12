package logger

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func Init(env string) error {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	logger, err := config.Build()
	if err != nil {
		return err
	}

	Logger = logger
	zap.ReplaceGlobals(logger)

	return nil
}

func Sync() {
	if Logger != nil {
		Logger.Sync()
	}
}
