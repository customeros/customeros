package logger

import (
	common_logger "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
)

type ExtendedLogger struct {
	common_logger.AppLogger
}

func NewExtendedAppLogger(cfg *common_logger.Config) *ExtendedLogger {
	appLogger := common_logger.NewAppLogger(cfg)
	return &ExtendedLogger{
		AppLogger: *appLogger,
	}
}

type Logger interface {
	common_logger.Logger
}
