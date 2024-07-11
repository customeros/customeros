package logger

import (
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	common_logger "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events"
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
		events.HTTP,
		zap.String(events.METHOD, method),
		zap.String(events.URI, uri),
		zap.Int(events.STATUS, status),
		zap.Int64(events.SIZE, size),
		zap.Duration(events.TIME, time),
	)
}

func (l *ExtendedLogger) GrpcMiddlewareAccessLogger(method string, time time.Duration, metaData map[string][]string, err error) {
	if err != nil {
		l.Logger().Info(
			events.GRPC,
			zap.String(events.METHOD, method),
			zap.Duration(events.TIME, time),
			zap.Any(events.METADATA, metaData),
			zap.String(events.ERROR, err.Error()),
		)
		return
	}
	l.Logger().Info(events.GRPC, zap.String(events.METHOD, method), zap.Duration(events.TIME, time), zap.Any(events.METADATA, metaData))
}

func (l *ExtendedLogger) GrpcClientInterceptorLogger(method string, req, reply interface{}, time time.Duration, metaData map[string][]string, err error) {
	if err != nil {
		l.Logger().Info(
			events.GRPC,
			zap.String(events.METHOD, method),
			zap.Any(events.REQUEST, req),
			zap.Any(events.REPLY, reply),
			zap.Duration(events.TIME, time),
			zap.Any(events.METADATA, metaData),
			zap.String(events.ERROR, err.Error()),
		)
		return
	}
	l.Logger().Info(
		events.GRPC,
		zap.String(events.METHOD, method),
		zap.Any(events.REQUEST, req),
		zap.Any(events.REPLY, reply),
		zap.Duration(events.TIME, time),
		zap.Any(events.METADATA, metaData),
	)
}

func (l *ExtendedLogger) EventAppeared(groupName string, event *esdb.ResolvedEvent, workerID int) {
	l.Logger().Info(
		"EventAppeared",
		zap.String(events.GroupName, groupName),
		zap.String(events.StreamID, event.OriginalEvent().StreamID),
		zap.String(events.EventID, event.OriginalEvent().EventID.String()),
		zap.String(events.EventType, event.OriginalEvent().EventType),
		zap.Uint64(events.EventNumber, event.OriginalEvent().EventNumber),
		zap.Time(events.CreatedDate, event.OriginalEvent().CreatedDate),
		zap.String(events.UserMetadata, string(event.OriginalEvent().UserMetadata)),
		zap.Int(events.WorkerID, workerID),
	)
}
