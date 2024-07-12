package server

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventbuffer"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstoredb"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/service"
	"google.golang.org/grpc"

	commonconf "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore/store"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
	return &Server{Config: cfg,
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

	//Server.metrics = metrics.NewESMicroserviceMetrics(Server.cfg)
	//Server.interceptorManager = interceptors.NewInterceptorManager(Server.log, Server.getGrpcMetricsCb())
	//Server.mw = middlewares.NewMiddlewareManager(Server.log, Server.cfg, Server.getHttpMetricsCb())

	esdb, err := eventstoredb.NewEventStoreDB(server.Config.EventStoreConfig, server.Log)
	if err != nil {
		return err
	}
	defer esdb.Close() // nolint: errcheck

	// Initialize postgres db
	postgresDb, _ := InitPostgresDB(server.Config, server.Log)
	defer postgresDb.SqlDB.Close()

	repository.Migration(postgresDb.GormDB)

	// Setting up Neo4j
	neo4jDriver, err := commonconf.NewNeo4jDriver(server.Config.Neo4j)
	if err != nil {
		logrus.Fatalf("Could not establish connection with neo4j at: %v, error: %v", server.Config.Neo4j.Target, err.Error())
	}
	defer neo4jDriver.Close(ctx)
	server.Repositories = repository.InitRepos(&neo4jDriver, server.Config.Neo4j.Database, postgresDb.GormDB)

	server.AggregateStore = store.NewAggregateStore(server.Log, esdb)

	bufferService := eventbuffer.NewEventBufferStoreService(server.Repositories.PostgresRepositories.EventBufferRepository, server.Log)
	server.CommandHandlers = command.NewCommandHandlers(server.Log, server.Config, server.AggregateStore, bufferService)

	//Server.runMetrics(cancel)
	//Server.runHealthCheck(ctx)

	server.Services = service.InitServices(server.Config, server.Repositories, server.AggregateStore, server.CommandHandlers, server.Log, bufferService)

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

func InitPostgresDB(cfg *config.Config, log logger.Logger) (db *commonconf.StorageDB, err error) {
	if db, err = commonconf.NewPostgresDBConn(cfg.Postgres); err != nil {
		log.Fatalf("Could not open db connection: %s", err.Error())
	}
	return
}
