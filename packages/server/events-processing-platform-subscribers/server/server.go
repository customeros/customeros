package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	commonconf "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/validator"
	"github.com/customeros/customeros/packages/server/events/eventstore"
	"github.com/customeros/customeros/packages/server/events/eventstore/store"
	"github.com/customeros/customeros/packages/server/events/eventstoredb"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/caches"
	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/config"
	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/eventbuffer"
	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/logger"
	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/repository"
	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/service"
	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/subscriptions"
	graph_subscription "github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/subscriptions/graph"
	invoice_subscription "github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/subscriptions/invoice"
)

const (
	waitShotDownDuration = 3 * time.Second
)

type Server struct {
	Config         *config.Config
	Log            logger.Logger
	Services       *service.Services
	AggregateStore eventstore.AggregateStore

	doneCh chan struct{}
	caches caches.Cache
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

	// Setting up eventstore subscriptions
	err = subscriptions.NewSubscriptions(server.Log, esdb, server.Config).RefreshSubscriptions(ctx)
	if err != nil {
		server.Log.Errorf("(graphConsumer.Connect) err: {%v}", err)
		cancel()
	}

	// Initialize postgres db
	postgresDb, err := commonconf.InitPostgres(&server.Config.CommonServices)
	if err != nil {
		logrus.Fatalf("failed opening connection to postgres: %v", err.Error())
	}
	defer postgresDb.Close()

	// Setting up Neo4j
	neo4jDriver, err := commonconf.NewNeo4jDriver(server.Config.CommonServices.Infrastructure.Neo4jConfig)
	if err != nil {
		logrus.Fatalf("Could not establish connection with neo4j at: %v, error: %v", server.Config.CommonServices.Infrastructure.Neo4jConfig.Target, err.Error())
	}
	defer neo4jDriver.Close(ctx)

	server.AggregateStore = store.NewAggregateStore(server.Log, esdb)

	// Setting up gRPC client
	df := grpc_client.NewDialFactory(&server.Config.GrpcClientConfig)
	gRPCconn, err := df.GetEventsProcessingPlatformConn()
	if err != nil {
		server.Log.Fatalf("Failed to connect: %v", err)
	}
	defer df.Close(gRPCconn)
	grpcClients := grpc_client.InitClients(gRPCconn)

	// setting up services
	repositories := repository.InitRepos(&neo4jDriver, server.Config.CommonServices.Infrastructure.Neo4jConfig.Database, postgresDb)

	server.Services = service.InitServices(
		server.Log,
		repositories.Neo4jRepositories,
		repositories.PostgresRepositories,
		server.Config.CommonServices,
		grpcClients,
		server.AggregateStore,
	)

	// Setting up cache
	server.caches = caches.InitCaches()

	eventBufferWatcher := eventbuffer.NewEventBufferWatcher(server.Services.PostgresRepositories.EventBufferRepository, server.Log, server.AggregateStore)
	eventBufferWatcher.Start(ctx)
	defer eventBufferWatcher.Stop()

	server.InitSubscribers(ctx, grpcClients, esdb, cancel)

	<-ctx.Done()
	server.waitShootDown(waitShotDownDuration)

	<-server.doneCh

	server.Log.Infof("%Server Server exited properly")
	return nil
}

func (server *Server) waitShootDown(duration time.Duration) {
	go func() {
		time.Sleep(duration)
		server.doneCh <- struct{}{}
	}()
}

func (server *Server) InitSubscribers(ctx context.Context, grpcClients *grpc_client.Clients, esdb *esdb.Client, cancel context.CancelFunc) {
	if server.Config.Subscriptions.GraphSubscription.Enabled {
		graphSubscriber := graph_subscription.NewGraphSubscriber(server.Log, esdb, server.Services, grpcClients, server.Config, server.caches)
		go func() {
			err := graphSubscriber.Connect(ctx, graphSubscriber.ProcessEvents)
			if err != nil {
				server.Log.Errorf("(graphSubscriber.Connect) err: {%s}", err.Error())
				cancel()
			}
		}()
	}

	if server.Config.Subscriptions.InvoiceSubscription.Enabled {
		invoiceSubscriber := invoice_subscription.NewInvoiceSubscriber(server.Log, esdb, server.Config, server.Services, grpcClients)
		go func() {
			err := invoiceSubscriber.Connect(ctx, invoiceSubscriber.ProcessEvents)
			if err != nil {
				server.Log.Errorf("(invoiceSubscriber.Connect) err: {%s}", err.Error())
				cancel()
			}
		}()
	}
}
