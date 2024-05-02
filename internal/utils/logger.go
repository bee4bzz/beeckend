package utils

import (
	zaplog "github.com/gaetanDubuc/beeckend/pkg/log"
)

// TODO: Should be in a log folder
func NewLogger() *zaplog.Logger {
	config, err := LoadConfig(".")

	if err != nil {
		panic("failed to load config")
	}

	var logger *zaplog.Logger
	if config.AppEnv == "development" {
		logger = zaplog.NewDevelopment()
	} else {
		logger = zaplog.NewProduction()
	}
	return logger
}
