package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	commonConfig "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/validator"
	"github.com/customeros/customeros/packages/server/events/eventstore"
	"github.com/customeros/customeros/packages/server/events/eventstore/store"
	"github.com/customeros/customeros/packages/server/events/eventstoredb"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/customeros/customeros/packages/server/events-processing-platform/config"
	"github.com/customeros/customeros/packages/server/events-processing-platform/domain/common/command"
	"github.com/customeros/customeros/packages/server/events-processing-platform/logger"
	"github.com/customeros/customeros/packages/server/events-processing-platform/repository"
	"github.com/customeros/customeros/packages/server/events-processing-platform/service"
)

const (
	waitShotDownDuration = 3 * time.Second
)

type Server struct {
	Config          *config.Config
	Log             logger.Logger
	Repositories    *repository.Repositories
	Services        *service.Services
	CommandHandlers *command.CommandHandlers
	AggregateStore  eventstore.AggregateStore
	GrpcServer      *grpc.Server

	doneCh chan struct{}
	//	metrics            *metrics.ESMicroserviceMetrics
}

func NewServer(cfg *config.Config, log logger.Logger) *Server {
	return &Server{
		Config: cfg,
		Log:    log,
		doneCh: make(chan struct{}),
	}
}

func (server *Server) Start(parentCtx context.Context) error {
	ctx, cancel := signal.NotifyContext(parentCtx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	if err := validator.GetValidator().Struct(server.Config); err != nil {
		return errors.Wrap(err, "cfg validate")
	}

	// Setting up tracing
	tracer, closer, err := tracing.NewJaegerTracer(&server.Config.Jaeger, server.Log)
	if err != nil {
		server.Log.Fatalf("Could not initialize jaeger tracer: %s", err.Error())
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// Server.metrics = metrics.NewESMicroserviceMetrics(Server.cfg)
	// Server.interceptorManager = interceptors.NewInterceptorManager(Server.log, Server.getGrpcMetricsCb())
	// Server.mw = middlewares.NewMiddlewareManager(Server.log, Server.cfg, Server.getHttpMetricsCb())

	esdb, err := eventstoredb.NewEventStoreDB(server.Config.EventStoreConfig, server.Log)
	if err != nil {
		return err
	}
	defer esdb.Close() // nolint: errcheck

	// Initialize postgres db
	postgresDb, err := commonConfig.InitPostgres(&server.Config.CommonServices)
	if err != nil {
		logrus.Fatalf("failed opening connection to postgres: %v", err.Error())
	}
	defer postgresDb.Close()

	repository.Migration(postgresDb.GormDB)

	// Setting up Neo4j
	neo4jDriver, err := commonConfig.NewNeo4jDriver(server.Config.CommonServices.Infrastructure.Neo4jConfig)
	if err != nil {
		logrus.Fatalf("Could not establish connection with neo4j at: %v, error: %v", server.Config.CommonServices.Infrastructure.Neo4jConfig.Target, err.Error())
	}
	defer neo4jDriver.Close(ctx)
	server.Repositories = repository.InitRepos(&neo4jDriver, server.Config.CommonServices.Infrastructure.Neo4jConfig.Database, postgresDb)

	server.AggregateStore = store.NewAggregateStore(server.Log, esdb)

	server.CommandHandlers = command.NewCommandHandlers(server.Log, server.Config, server.AggregateStore)

	// Server.runMetrics(cancel)
	// Server.runHealthCheck(ctx)

	server.Services = service.InitServices(
		server.Config,
		server.Repositories,
		server.AggregateStore,
		server.CommandHandlers,
		server.Log,
	)

	closeGrpcServer, grpcServer, err := server.NewEventProcessorGrpcServer()
	if err != nil {
		cancel()
		return err
	}
	defer closeGrpcServer()
	server.GrpcServer = grpcServer

	<-ctx.Done()
	server.waitShootDown(waitShotDownDuration)

	grpcServer.GracefulStop()

	<-server.doneCh

	server.Log.Infof("%Server Server exited properly", GetMicroserviceName(server.Config))
	return nil
}

func (server *Server) waitShootDown(duration time.Duration) {
	go func() {
		time.Sleep(duration)
		server.doneCh <- struct{}{}
	}()
}
