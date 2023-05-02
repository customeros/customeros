package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain"
	contact_commands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/commands"
	email_commands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/commands"
	organization_commands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/commands"
	phone_number_commands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/commands"
	user_commands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/commands"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore/store"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstroredb"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	email_validation_projection "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/projection/email_validation"
	graph_projection "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/projection/graph"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

type server struct {
	cfg          *config.Config
	log          logger.Logger
	repositories *repository.Repositories
	commands     *domain.Commands
	echo         *echo.Echo
	doneCh       chan struct{}
	//validate           *validator.Validate
	//	metrics            *metrics.ESMicroserviceMetrics
}

func NewServer(cfg *config.Config, log logger.Logger) *server {
	return &server{cfg: cfg, log: log, echo: echo.New(), doneCh: make(chan struct{})}
}

func (server *server) Run(parentCtx context.Context) error {
	ctx, cancel := signal.NotifyContext(parentCtx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	validator.InitValidator()

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
	//server.interceptorManager = interceptors.NewInterceptorManager(server.log, server.getGrpcMetricsCb())
	//server.mw = middlewares.NewMiddlewareManager(server.log, server.cfg, server.getHttpMetricsCb())

	db, err := eventstroredb.NewEventStoreDB(server.cfg.EventStoreConfig)
	if err != nil {
		return err
	}
	defer db.Close() // nolint: errcheck

	// Setting up Neo4j
	neo4jDriver, err := config.NewDriver(server.cfg)
	if err != nil {
		logrus.Fatalf("Could not establish connection with neo4j at: %v, error: %v", server.cfg.Neo4j.Target, err.Error())
	}
	defer neo4jDriver.Close(ctx)
	server.repositories = repository.InitRepos(&neo4jDriver)

	aggregateStore := store.NewAggregateStore(server.log, db)
	server.commands = &domain.Commands{
		ContactCommands:      contact_commands.NewContactCommands(server.log, server.cfg, aggregateStore),
		OrganizationCommands: organization_commands.NewOrganizationCommands(server.log, server.cfg, aggregateStore),
		PhoneNumberCommands:  phone_number_commands.NewPhoneNumberCommands(server.log, server.cfg, aggregateStore),
		EmailCommands:        email_commands.NewEmailCommands(server.log, server.cfg, aggregateStore),
		UserCommands:         user_commands.NewUserCommands(server.log, server.cfg, aggregateStore),
	}

	graphProjection := graph_projection.NewGraphProjection(server.log, db, server.repositories, server.cfg)
	go func() {
		prefixes := []string{
			server.cfg.Subscriptions.ContactPrefix,
			server.cfg.Subscriptions.OrganizationPrefix,
			server.cfg.Subscriptions.PhoneNumberPrefix,
			server.cfg.Subscriptions.EmailPrefix,
			server.cfg.Subscriptions.UserPrefix,
		}
		err := graphProjection.Subscribe(ctx, prefixes, server.cfg.Subscriptions.PoolSize, graphProjection.ProcessEvents)
		if err != nil {
			server.log.Errorf("(graphProjection.Subscribe) err: {%v}", err)
			cancel()
		}
	}()

	dataValidationProjection := email_validation_projection.NewEmailValidationProjection(server.log, db, server.cfg, server.commands.EmailCommands)
	go func() {
		prefixes := []string{
			server.cfg.Subscriptions.EmailPrefix,
		}
		err := dataValidationProjection.Subscribe(ctx, prefixes, server.cfg.Subscriptions.PoolSize, dataValidationProjection.ProcessEvents)
		if err != nil {
			server.log.Errorf("(dataValidationProjection.Subscribe) err: {%v}", err)
			cancel()
		}
	}()

	//server.runMetrics(cancel)
	//server.runHealthCheck(ctx)

	closeGrpcServer, grpcServer, err := server.newEventProcessorGrpcServer()
	if err != nil {
		cancel()
		return err
	}
	defer closeGrpcServer() // nolint: errcheck

	<-ctx.Done()
	server.waitShootDown(waitShotDownDuration)

	grpcServer.GracefulStop()

	if err := server.echo.Shutdown(ctx); err != nil {
		server.log.Warnf("(Shutdown) err: {%validate}", err)
	}

	<-server.doneCh
	server.log.Infof("%server server exited properly", GetMicroserviceName(server.cfg))
	return nil
}
