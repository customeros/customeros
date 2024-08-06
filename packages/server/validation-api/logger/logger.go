package logger

import (
	commonlogger "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
)

type ExtendedLogger struct {
	commonlogger.AppLogger
}

func NewExtendedAppLogger(cfg *commonlogger.Config) *ExtendedLogger {
	appLogger := commonlogger.NewAppLogger(cfg)
	return &ExtendedLogger{
		AppLogger: *appLogger,
	}
}

type Logger interface {
	commonlogger.Logger
}
