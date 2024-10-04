package logger

import (
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	common_logger "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"go.uber.org/zap"
	"time"
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

// Logger methods interface
type Logger interface {
	common_logger.Logger
	HttpMiddlewareAccessLogger(method string, uri string, status int, size int64, time time.Duration)
	GrpcMiddlewareAccessLogger(method string, time time.Duration, metaData map[string][]string, err error)
	GrpcClientInterceptorLogger(method string, req interface{}, reply interface{}, time time.Duration, metaData map[string][]string, err error)
	EventAppeared(groupName string, event *esdb.ResolvedEvent, workerID int)
}

func (l *ExtendedLogger) HttpMiddlewareAccessLogger(method, uri string, status int, size int64, time time.Duration) {
	l.Logger().Info(
		constants.HTTP,
		zap.String(constants.METHOD, method),
		zap.String(constants.URI, uri),
		zap.Int(constants.STATUS, status),
		zap.Int64(constants.SIZE, size),
		zap.Duration(constants.TIME, time),
	)
}

func (l *ExtendedLogger) GrpcMiddlewareAccessLogger(method string, time time.Duration, metaData map[string][]string, err error) {
	if err != nil {
		l.Logger().Info(
			constants.GRPC,
			zap.String(constants.METHOD, method),
			zap.Duration(constants.TIME, time),
			zap.Any(constants.METADATA, metaData),
			zap.String(constants.ERROR, err.Error()),
		)
		return
	}
	l.Logger().Info(constants.GRPC, zap.String(constants.METHOD, method), zap.Duration(constants.TIME, time), zap.Any(constants.METADATA, metaData))
}

func (l *ExtendedLogger) GrpcClientInterceptorLogger(method string, req, reply interface{}, time time.Duration, metaData map[string][]string, err error) {
	if err != nil {
		l.Logger().Info(
			constants.GRPC,
			zap.String(constants.METHOD, method),
			zap.Any(constants.REQUEST, req),
			zap.Any(constants.REPLY, reply),
			zap.Duration(constants.TIME, time),
			zap.Any(constants.METADATA, metaData),
			zap.String(constants.ERROR, err.Error()),
		)
		return
	}
	l.Logger().Info(
		constants.GRPC,
		zap.String(constants.METHOD, method),
		zap.Any(constants.REQUEST, req),
		zap.Any(constants.REPLY, reply),
		zap.Duration(constants.TIME, time),
		zap.Any(constants.METADATA, metaData),
	)
}

func (l *ExtendedLogger) EventAppeared(groupName string, event *esdb.ResolvedEvent, workerID int) {
	l.Logger().Info(
		"EventAppeared",
		zap.String(constants.GroupName, groupName),
		zap.String(constants.StreamID, event.OriginalEvent().StreamID),
		zap.String(constants.EventID, event.OriginalEvent().EventID.String()),
		zap.String(constants.EventType, event.OriginalEvent().EventType),
		zap.Uint64(constants.EventNumber, event.OriginalEvent().EventNumber),
		zap.Time(constants.CreatedDate, event.OriginalEvent().CreatedDate),
		zap.String(constants.UserMetadata, string(event.OriginalEvent().UserMetadata)),
		zap.Int(constants.WorkerID, workerID),
	)
}
