package server

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	validator "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/eventbuffer"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/subscriber"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstoredb"
	"github.com/opentracing/opentracing-go"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	commonconf "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	graph_subscription "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/graph"
	invoice_subscription "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/invoice"
	location_validation_subscription "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/location"
	notifications_subscription "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/notifications"
	organization_subscription "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/organization"
	phone_number_validation_subscription "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/phone_number_validation"
	remindersubscription "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/reminder"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore/store"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

	// Setting up eventstore subscriptions
	err = subscriptions.NewSubscriptions(server.Log, esdb, server.Config).RefreshSubscriptions(ctx)
	if err != nil {
		server.Log.Errorf("(graphConsumer.Connect) err: {%v}", err)
		cancel()
	}

	// Initialize postgres db
	postgresDb, _ := InitPostgresDB(server.Config, server.Log)
	defer postgresDb.SqlDB.Close()

	// Setting up Neo4j
	neo4jDriver, err := commonconf.NewNeo4jDriver(server.Config.Neo4j)
	if err != nil {
		logrus.Fatalf("Could not establish connection with neo4j at: %v, error: %v", server.Config.Neo4j.Target, err.Error())
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

	server.Services = service.InitServices(server.Config, server.AggregateStore, server.Log, grpcClients, postgresDb.GormDB, &neo4jDriver)

	// Setting up cache
	industryMap, _ := server.Services.CommonServices.PostgresRepositories.IndustryMappingRepository.GetAllIndustryMappingsAsMap(ctx)
	server.caches = caches.InitCaches(industryMap)

	eventBufferWatcher := eventbuffer.NewEventBufferWatcher(server.Services.CommonServices.PostgresRepositories.EventBufferRepository, server.Log, server.AggregateStore)
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

func InitPostgresDB(cfg *config.Config, log logger.Logger) (db *commonconf.StorageDB, err error) {
	if db, err = commonconf.NewPostgresDBConn(cfg.Postgres); err != nil {
		log.Fatalf("Could not open db connection: %s", err.Error())
	}
	return
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

	if server.Config.Subscriptions.PhoneNumberValidationSubscription.Enabled {
		phoneNumberValidationSubscriber := phone_number_validation_subscription.NewPhoneNumberValidationSubscriber(server.Log, esdb, server.Config, server.Services, grpcClients)
		go func() {
			err := phoneNumberValidationSubscriber.Connect(ctx, phoneNumberValidationSubscriber.ProcessEvents)
			if err != nil {
				server.Log.Errorf("(phoneNumberValidationSubscriber.Connect) err: {%s}", err.Error())
				cancel()
			}
		}()
	}

	if server.Config.Subscriptions.LocationValidationSubscription.Enabled {
		locationValidationSubscriber := location_validation_subscription.NewLocationValidationSubscriber(server.Log, esdb, server.Config, server.Services, grpcClients)
		go func() {
			err := locationValidationSubscriber.Connect(ctx, locationValidationSubscriber.ProcessEvents)
			if err != nil {
				server.Log.Errorf("(locationValidationSubscriber.Connect) err: {%s}", err.Error())
				cancel()
			}
		}()
	}

	if server.Config.Subscriptions.OrganizationSubscription.Enabled {
		organizationSubscriber := organization_subscription.NewOrganizationSubscriber(server.Log, esdb, server.Config, server.Services, server.caches, grpcClients)
		go func() {
			err := organizationSubscriber.Connect(ctx, organizationSubscriber.ProcessEvents)
			if err != nil {
				server.Log.Errorf("(organizationSubscriber.Connect) err: {%s}", err.Error())
				cancel()
			}
		}()
	}

	if server.Config.Subscriptions.OrganizationWebscrapeSubscription.Enabled {
		organizationWebscrapeSubscriber := organization_subscription.NewOrganizationEnrichSubscriber(server.Log, esdb, server.Config, server.caches, grpcClients, server.Services)
		go func() {
			err := organizationWebscrapeSubscriber.Connect(ctx, organizationWebscrapeSubscriber.ProcessEvents)
			if err != nil {
				server.Log.Errorf("(organizationWebscrapeSubscriber.Connect) err: {%s}", err.Error())
				cancel()
			}
		}()
	}

	if server.Config.Subscriptions.NotificationsSubscription.Enabled {
		notificationsSubscriber := notifications_subscription.NewNotificationsSubscriber(server.Log, esdb, server.Services, server.Config)
		go func() {
			err := notificationsSubscriber.Connect(ctx, notificationsSubscriber.ProcessEvents)
			if err != nil {
				server.Log.Errorf("(notificationsSubscriber.Connect) err: {%s}", err.Error())
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

	if server.Config.Subscriptions.EnrichSubscription.Enabled {
		enrichSubscriber := subscriber.NewEnrichSubscriber(server.Log, esdb, server.Config, server.Services, server.caches, grpcClients)
		go func() {
			err := enrichSubscriber.Connect(ctx, enrichSubscriber.ProcessEvents)
			if err != nil {
				server.Log.Errorf("(enrichSubscriber.Connect) err: {%s}", err.Error())
				cancel()
			}
		}()
	}

	if server.Config.Subscriptions.ReminderSubscription.Enabled {
		reminderSubscriber := remindersubscription.NewReminderSubscriber(server.Log, esdb, server.Config, server.Services)
		go func() {
			err := reminderSubscriber.Connect(ctx, reminderSubscriber.ProcessEvents)
			if err != nil {
				server.Log.Errorf("(reminderSubscriber.Connect) err: {%s}", err.Error())
				cancel()
			}
		}()
	}
}
