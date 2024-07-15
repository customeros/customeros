package logger

import (
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	common_logger "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events/utils"
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
		utils.HTTP,
		zap.String(utils.METHOD, method),
		zap.String(utils.URI, uri),
		zap.Int(utils.STATUS, status),
		zap.Int64(utils.SIZE, size),
		zap.Duration(utils.TIME, time),
	)
}

func (l *ExtendedLogger) GrpcMiddlewareAccessLogger(method string, time time.Duration, metaData map[string][]string, err error) {
	if err != nil {
		l.Logger().Info(
			utils.GRPC,
			zap.String(utils.METHOD, method),
			zap.Duration(utils.TIME, time),
			zap.Any(utils.METADATA, metaData),
			zap.String(utils.ERROR, err.Error()),
		)
		return
	}
	l.Logger().Info(utils.GRPC, zap.String(utils.METHOD, method), zap.Duration(utils.TIME, time), zap.Any(utils.METADATA, metaData))
}

func (l *ExtendedLogger) GrpcClientInterceptorLogger(method string, req, reply interface{}, time time.Duration, metaData map[string][]string, err error) {
	if err != nil {
		l.Logger().Info(
			utils.GRPC,
			zap.String(utils.METHOD, method),
			zap.Any(utils.REQUEST, req),
			zap.Any(utils.REPLY, reply),
			zap.Duration(utils.TIME, time),
			zap.Any(utils.METADATA, metaData),
			zap.String(utils.ERROR, err.Error()),
		)
		return
	}
	l.Logger().Info(
		utils.GRPC,
		zap.String(utils.METHOD, method),
		zap.Any(utils.REQUEST, req),
		zap.Any(utils.REPLY, reply),
		zap.Duration(utils.TIME, time),
		zap.Any(utils.METADATA, metaData),
	)
}

func (l *ExtendedLogger) EventAppeared(groupName string, event *esdb.ResolvedEvent, workerID int) {
	l.Logger().Info(
		"EventAppeared",
		zap.String(utils.GroupName, groupName),
		zap.String(utils.StreamID, event.OriginalEvent().StreamID),
		zap.String(utils.EventID, event.OriginalEvent().EventID.String()),
		zap.String(utils.EventType, event.OriginalEvent().EventType),
		zap.Uint64(utils.EventNumber, event.OriginalEvent().EventNumber),
		zap.Time(utils.CreatedDate, event.OriginalEvent().CreatedDate),
		zap.String(utils.UserMetadata, string(event.OriginalEvent().UserMetadata)),
		zap.Int(utils.WorkerID, workerID),
	)
}
