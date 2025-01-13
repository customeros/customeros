package logger

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"

type ExtendedLogger struct {
	logger.AppLogger
}

func NewExtendedAppLogger(cfg *logger.Config) *ExtendedLogger {
	appLogger := logger.NewAppLogger(cfg)
	return &ExtendedLogger{
		AppLogger: *appLogger,
	}
}

type Logger interface {
	logger.Logger
}
