package server

import (
	"bytes"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	commonConfig "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	commonservice "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/validator"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/customeros/customeros/packages/server/customer-os-webhooks/caches"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/config"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/constants"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/repository"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/route"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/service"
)

type server struct {
	cfg    *config.Config
	log    logger.Logger
	doneCh chan struct{}
}

func NewServer(cfg *config.Config, log logger.Logger) *server {
	return &server{cfg: cfg, log: log, doneCh: make(chan struct{})}
}

func (server *server) Run(parentCtx context.Context) error {
	ctx, cancel := signal.NotifyContext(parentCtx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	if err := validator.GetValidator().Struct(server.cfg); err != nil {
		return errors.Wrap(err, "cfg validate")
	}

	// Setting up tracing
	tracer, closer, err := tracing.NewJaegerTracer(&server.cfg.Common.Infrastructure.JaegerConfig, server.log)
	if err != nil {
		server.log.Fatalf("Could not initialize jaeger tracer: %s", err.Error())
	}
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	registerPrometheusMetrics()

	// Initialize postgres db
	postgresDb, err := commonConfig.InitPostgres(&server.cfg.Common)
	if err != nil {
		logrus.Fatalf("failed opening connection to postgres: %v", err.Error())
	}
	defer postgresDb.Close()

	// Setting up Neo4j
	neo4jDriver, err := commonConfig.NewNeo4jDriver(server.cfg.Common.Infrastructure.Neo4jConfig)
	if err != nil {
		server.log.Fatalf("Could not establish connection with neo4j at: %v, error: %v", server.cfg.Common.Infrastructure.Neo4jConfig.Target, err.Error())
	}
	defer neo4jDriver.Close(ctx)

	// Setting up gRPC client
	df := grpc_client.NewDialFactory(&server.cfg.Common.Infrastructure.GrpcClientConfig)
	gRPCconn, err := df.GetEventsProcessingPlatformConn()
	if err != nil {
		server.log.Fatalf("Failed to connect: %v", err)
	}
	defer df.Close(gRPCconn)
	grpcContainer := grpc_client.InitClients(gRPCconn)

	// Setting up CommonServices & repositories
	repos := repository.InitRepos(&neo4jDriver, postgresDb, server.cfg.Common.Infrastructure.Neo4jConfig.Database)

	commonServices := commonservice.InitCommonServices(
		server.log,
		repos.Neo4jRepositories,
		repos.PostgresRepositories,
		&server.cfg.Common,
		grpcContainer,
		&commonservice.InitOptions{
			LoadPersonalEmailProviders: true,
			LoadEmailExclusionList:     true,
		},
	)

	// Setting up Gin
	r := gin.Default()

	// Setting up CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}

	r.Use(cors.New(corsConfig))
	r.Use(ginzap.GinzapWithConfig(server.log.Logger(), &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
		SkipPaths:  []string{"/metrics", "/health", "/readiness", "/"},
	}))
	r.Use(ginzap.RecoveryWithZap(server.log.Logger(), true))
	r.Use(tracing.RecoveryWithJaeger(opentracing.GlobalTracer()))
	r.Use(prometheusMiddleware())
	r.Use(bodyLoggerMiddleware)

	// Setting up caches
	appCache := caches.NewCache()

	// Setting up services
	serviceContainer := service.InitServices(
		server.log,
		repos,
		server.cfg,
		commonServices,
		grpcContainer,
		appCache,
	)
	route.AddExternalSystemRoutes(ctx, r, serviceContainer, server.log, serviceContainer.CommonServices.Cache)
	route.AddUserRoutes(ctx, r, serviceContainer, server.log, serviceContainer.CommonServices.Cache)
	route.AddOrganizationRoutes(ctx, r, serviceContainer, server.log, serviceContainer.CommonServices.Cache)
	route.AddLogEntryRoutes(ctx, r, serviceContainer, server.log, serviceContainer.CommonServices.Cache)
	route.AddContactRoutes(ctx, r, server.cfg, serviceContainer, server.log, serviceContainer.CommonServices.Cache)
	route.AddIssueRoutes(ctx, r, serviceContainer, server.log, serviceContainer.CommonServices.Cache)
	route.AddInteractionEventRoutes(ctx, r, serviceContainer, server.cfg, server.log, serviceContainer.CommonServices.Cache)
	route.AddCommentRoutes(ctx, r, serviceContainer, server.log, serviceContainer.CommonServices.Cache)
	route.AddInvoiceRoutes(ctx, r, serviceContainer, server.log, serviceContainer.CommonServices.Cache)
	route.AddSlackRoutes(ctx, r, serviceContainer, server.log, serviceContainer.CommonServices.Cache)
	route.AddEmailRoutes(ctx, r, server.cfg, serviceContainer)

	r.GET("/health", HealthCheckHandler)
	r.GET("/readiness", ReadinessHandler)
	r.GET("/", RootHandler)

	if server.cfg.App.ApiPort == server.cfg.App.MetricsPort {
		r.GET(server.cfg.App.Metrics.PrometheusPath, metricsHandler)
	} else {
		go func() {
			mr := gin.Default()
			mr.Use(prometheusMiddleware())
			mr.Use(bodyLoggerMiddleware)
			mr.GET(server.cfg.App.Metrics.PrometheusPath, metricsHandler)
			mr.Run(":" + server.cfg.App.MetricsPort)
		}()
	}

	r.Run(":" + server.cfg.App.ApiPort)

	<-server.doneCh
	server.log.Infof("Application %s exited properly", constants.ServiceName)
	return nil
}

func HealthCheckHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "OK"})
}

func ReadinessHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "READY"})
}

func RootHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Customer OS Webhooks",
	})
}

func metricsHandler(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}

func registerPrometheusMetrics() {
	// Implement metrics invocations here
}

func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func(start time.Time) {
			// TODO implement metrics COS-314 https://linear.app/customer-os/issue/COS-314/add-prometheus-metrics-on-success-and-failed-webhook-rest-api-calls
			// TODO count duration / success / failed requests
		}(time.Now())
		c.Next()
	}
}

func bodyLoggerMiddleware(c *gin.Context) {
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	c.Next()
	c.Set("bodyBytes", blw.body.Bytes())
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
