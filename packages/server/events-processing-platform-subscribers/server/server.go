package server

import (
	"context"
	validator "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/labstack/echo/v4"
	commonconf "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	email_validation_subscription "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/email_validation"
	graph_subscription "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/graph"
	graph_low_prio_subscription "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/graph_low_prio"
	interaction_event_subscription "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/interaction_event"
	invoice_subscription "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/invoice"
	location_validation_subscription "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/location_validation"
	notifications_subscription "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/notifications"
	organization_subscription "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/organization"
	phone_number_validation_subscription "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/phone_number_validation"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore/store"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstroredb"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	waitShotDownDuration = 3 * time.Second
)

type Server struct {
	Config         *config.Config
	Log            logger.Logger
	Repositories   *repository.Repositories
	Services       *service.Services
	AggregateStore eventstore.AggregateStore

	echo   *echo.Echo
	doneCh chan struct{}
	caches caches.Cache
	//	metrics            *metrics.ESMicroserviceMetrics
}

func NewServer(cfg *config.Config, log logger.Logger) *Server {
	return &Server{Config: cfg,
		Log:    log,
		echo:   echo.New(),
		doneCh: make(chan struct{}),
		caches: caches.InitCaches(),
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

	esdb, err := eventstroredb.NewEventStoreDB(server.Config.EventStoreConfig, server.Log)
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

	repository.Migration(postgresDb.GormDB)

	// Setting up Neo4j
	neo4jDriver, err := commonconf.NewNeo4jDriver(server.Config.Neo4j)
	if err != nil {
		logrus.Fatalf("Could not establish connection with neo4j at: %v, error: %v", server.Config.Neo4j.Target, err.Error())
	}
	defer neo4jDriver.Close(ctx)
	server.Repositories = repository.InitRepos(&neo4jDriver, server.Config.Neo4j.Database, postgresDb.GormDB)

	server.AggregateStore = store.NewAggregateStore(server.Log, esdb)

	//Server.runMetrics(cancel)
	//Server.runHealthCheck(ctx)

	server.Services = service.InitServices(server.Config, server.Repositories, server.Log)

	// Setting up gRPC client
	df := grpc_client.NewDialFactory(server.Config)
	gRPCconn, err := df.GetEventsProcessingPlatformConn()
	if err != nil {
		server.Log.Fatalf("Failed to connect: %v", err)
	}
	defer df.Close(gRPCconn)
	grpcClients := grpc_client.InitGrpcClients(gRPCconn)

	InitSubscribers(server, ctx, grpcClients, esdb, cancel, server.Services)

	<-ctx.Done()
	server.waitShootDown(waitShotDownDuration)

	if err := server.echo.Shutdown(ctx); err != nil {
		server.Log.Warnf("(Shutdown) err: {%validate}", err)
	}

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

func InitSubscribers(server *Server, ctx context.Context, grpcClients *grpc_client.Clients, esdb *esdb.Client, cancel context.CancelFunc, services *service.Services) {
	if server.Config.Subscriptions.GraphSubscription.Enabled {
		graphSubscriber := graph_subscription.NewGraphSubscriber(server.Log, esdb, server.Repositories, grpcClients, server.Config)
		go func() {
			err := graphSubscriber.Connect(ctx, graphSubscriber.ProcessEvents)
			if err != nil {
				server.Log.Errorf("(graphSubscriber.Connect) err: {%v}", err)
				cancel()
			}
		}()
	}

	if server.Config.Subscriptions.GraphLowPrioritySubscription.Enabled {
		subscriber := graph_low_prio_subscription.NewGraphLowPrioSubscriber(server.Log, esdb, server.Repositories, grpcClients, server.Config)
		go func() {
			err := subscriber.Connect(ctx, subscriber.ProcessEvents)
			if err != nil {
				server.Log.Errorf("(graphLowPrioSubscriber.Connect) err: {%v}", err)
				cancel()
			}
		}()
	}

	if server.Config.Subscriptions.EmailValidationSubscription.Enabled {
		emailValidationSubscriber := email_validation_subscription.NewEmailValidationSubscriber(server.Log, esdb, server.Config, grpcClients)
		go func() {
			err := emailValidationSubscriber.Connect(ctx, emailValidationSubscriber.ProcessEvents)
			if err != nil {
				server.Log.Errorf("(emailValidationSubscriber.Connect) err: {%v}", err)
				cancel()
			}
		}()
	}

	if server.Config.Subscriptions.PhoneNumberValidationSubscription.Enabled {
		phoneNumberValidationSubscriber := phone_number_validation_subscription.NewPhoneNumberValidationSubscriber(server.Log, esdb, server.Config, server.Repositories, grpcClients)
		go func() {
			err := phoneNumberValidationSubscriber.Connect(ctx, phoneNumberValidationSubscriber.ProcessEvents)
			if err != nil {
				server.Log.Errorf("(phoneNumberValidationSubscriber.Connect) err: {%v}", err)
				cancel()
			}
		}()
	}

	if server.Config.Subscriptions.LocationValidationSubscription.Enabled {
		locationValidationSubscriber := location_validation_subscription.NewLocationValidationSubscriber(server.Log, esdb, server.Config, server.Repositories, grpcClients)
		go func() {
			err := locationValidationSubscriber.Connect(ctx, locationValidationSubscriber.ProcessEvents)
			if err != nil {
				server.Log.Errorf("(locationValidationSubscriber.Connect) err: {%v}", err)
				cancel()
			}
		}()
	}

	if server.Config.Subscriptions.OrganizationSubscription.Enabled {
		organizationSubscriber := organization_subscription.NewOrganizationSubscriber(server.Log, esdb, server.Config, server.Repositories, server.caches, grpcClients)
		go func() {
			err := organizationSubscriber.Connect(ctx, organizationSubscriber.ProcessEvents)
			if err != nil {
				server.Log.Errorf("(organizationSubscriber.Connect) err: {%v}", err)
				cancel()
			}
		}()
	}

	if server.Config.Subscriptions.OrganizationWebscrapeSubscription.Enabled {
		organizationWebscrapeSubscriber := organization_subscription.NewOrganizationWebscrapeSubscriber(server.Log, esdb, server.Config, server.Repositories, server.caches, grpcClients)
		go func() {
			err := organizationWebscrapeSubscriber.Connect(ctx, organizationWebscrapeSubscriber.ProcessEvents)
			if err != nil {
				server.Log.Errorf("(organizationWebscrapeSubscriber.Connect) err: {%v}", err)
				cancel()
			}
		}()
	}

	if server.Config.Subscriptions.InteractionEventSubscription.Enabled {
		interactionEventSubscriber := interaction_event_subscription.NewInteractionEventSubscriber(server.Log, esdb, server.Config, server.Repositories, grpcClients)
		go func() {
			err := interactionEventSubscriber.Connect(ctx, interactionEventSubscriber.ProcessEvents)
			if err != nil {
				server.Log.Errorf("(interactionEventSubscriber.Connect) err: {%v}", err)
				cancel()
			}
		}()
	}

	if server.Config.Subscriptions.NotificationsSubscription.Enabled {
		notificationsSubscriber := notifications_subscription.NewNotificationsSubscriber(server.Log, esdb, server.Repositories, grpcClients, server.Config)
		go func() {
			err := notificationsSubscriber.Connect(ctx, notificationsSubscriber.ProcessEvents)
			if err != nil {
				server.Log.Errorf("(notificationsSubscriber.Connect) err: {%v}", err)
				cancel()
			}
		}()
	}

	if server.Config.Subscriptions.InvoiceSubscription.Enabled {
		invoiceSubscriber := invoice_subscription.NewInvoiceSubscriber(server.Log, esdb, server.Config, server.Repositories, grpcClients, services.FileStoreApiService, services.PostmarkProvider)
		go func() {
			err := invoiceSubscriber.Connect(ctx, invoiceSubscriber.ProcessEvents)
			if err != nil {
				server.Log.Errorf("(invoiceSubscriber.Connect) err: {%v}", err)
				cancel()
			}
		}()
	}
}
