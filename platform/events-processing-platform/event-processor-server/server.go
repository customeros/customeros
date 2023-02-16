package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/domain/contacts/service"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/eventstore/store"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/eventstroredb"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/interceptors"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/middlewares"
	"os"
	"os/signal"
	"syscall"
)

type server struct {
	cfg                *config.Config
	log                logger.Logger
	interceptorManager interceptors.InterceptorManager
	mw                 middlewares.MiddlewareManager
	contactService     *service.ContactService
	//validate           *validator.Validate

	echo *echo.Echo
	//	metrics            *metrics.ESMicroserviceMetrics
	doneCh chan struct{}
}

func NewServer(cfg *config.Config, log logger.Logger) *server {
	return &server{cfg: cfg, log: log, echo: echo.New(), doneCh: make(chan struct{})}
}

func (server *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	//if err := server.validate.StructCtx(ctx, server.cfg); err != nil {
	//	return errors.Wrap(err, "cfg validate")
	//}

	/*	if server.cfg.Jaeger.Enable {
		tracer, closer, err := tracing.NewJaegerTracer(server.cfg.Jaeger)
		if err != nil {
			return err
		}
		defer closer.Close() // nolint: errcheck
		opentracing.SetGlobalTracer(tracer)
	}*/

	//server.metrics = metrics.NewESMicroserviceMetrics(server.cfg)
	server.interceptorManager = interceptors.NewInterceptorManager(server.log, server.getGrpcMetricsCb())
	server.mw = middlewares.NewMiddlewareManager(server.log, server.cfg, server.getHttpMetricsCb())

	db, err := eventstroredb.NewEventStoreDB(server.cfg.EventStoreConfig)
	if err != nil {
		return err
	}
	defer db.Close() // nolint: errcheck

	aggregateStore := store.NewAggregateStore(server.log, db)
	server.contactService = service.NewContactService(server.log, server.cfg, aggregateStore)

	//server.runMetrics(cancel)
	//server.runHealthCheck(ctx)

	//go func() {
	//	if err := server.runHttpServer(); err != nil {
	//		server.log.Errorf("(server.runHttpServer) err: {%validate}", err)
	//		cancel()
	//	}
	//}()
	//server.log.Infof("%server is listening on PORT: {%server}", GetMicroserviceName(server.cfg), server.cfg.Http.Port)

	closeGrpcServer, grpcServer, err := server.newEventProcessorGrpcServer()
	if err != nil {
		cancel()
		return err
	}
	defer closeGrpcServer() // nolint: errcheck

	<-ctx.Done()
	server.waitShootDown(waitShotDownDuration)

	grpcServer.GracefulStop()
	if err := server.shutDownHealthCheckServer(ctx); err != nil {
		server.log.Warnf("(shutDownHealthCheckServer) err: {%validate}", err)
	}
	if err := server.echo.Shutdown(ctx); err != nil {
		server.log.Warnf("(Shutdown) err: {%validate}", err)
	}

	<-server.doneCh
	server.log.Infof("%server server exited properly", GetMicroserviceName(server.cfg))
	return nil
}
