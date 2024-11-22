package server

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	commonconf "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	commonlogger "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	commontracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/routes"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/service"
)

type server struct {
	cfg      *config.Config
	log      commonlogger.Logger
	doneCh   chan struct{}
	engine   *gin.Engine
	services *service.Services
	cache    *caches.Cache
}

func NewServer(cfg *config.Config, log commonlogger.Logger) *server {
	return &server{
		cfg:    cfg,
		log:    log,
		doneCh: make(chan struct{}),
		engine: gin.Default(),
		cache:  caches.NewCache(),
	}
}

func (s *server) Run(parentCtx context.Context) error {
	if err := s.validateConfig(); err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(parentCtx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	_, closer, err := s.setupTracing()
	if err != nil {
		return err
	}
	defer closer.Close()

	if err := s.setupDependencies(ctx); err != nil {
		return err
	}

	s.setupMiddleware()

	routeRegistry := routes.NewRegistry(
		s.engine,
		s.services,
		s.cfg,
		s.log,
		s.cache,
	)
	routeRegistry.RegisterRoutes(ctx)

	s.setupMetrics()

	go s.startMetricsServer()

	s.log.Infof("Starting %s on port %s", constants.ServiceName, s.cfg.ApiPort)
	return s.startMainServer()
}

func (s *server) validateConfig() error {
	if err := validator.GetValidator().Struct(s.cfg); err != nil {
		return errors.Wrap(err, "cfg validate")
	}
	return nil
}

func (s *server) setupTracing() (opentracing.Tracer, io.Closer, error) {
	tracer, closer, err := commontracing.NewJaegerTracer(&s.cfg.Jaeger, s.log)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not initialize jaeger tracer")
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer, nil
}

func (s *server) initPostgresDB() (*commonconf.StorageDB, error) {
	db, err := commonconf.NewPostgresDBConn(s.cfg.Postgres)
	if err != nil {
		return nil, errors.Wrap(err, "could not open db connection")
	}
	if db == nil {
		return nil, errors.New("postgres connection is nil")
	}
	return db, nil
}

func (s *server) initNeo4j() (neo4j.Driver, error) {
	driver, err := commonconf.NewNeo4jDriver(s.cfg.Neo4j)
	if err != nil {
		return nil, errors.Wrap(err, "could not establish connection with neo4j")
	}
	return driver, nil
}

func (s *server) initGRPC() (*grpc_client.Container, error) {
	df := grpc_client.NewDialFactory(&s.cfg.GrpcClientConfig)
	conn, err := df.GetEventsProcessingPlatformConn()
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to gRPC")
	}
	return grpc_client.InitClients(conn), nil
}

func (s *server) initializeServices(postgresDb *commonconf.StorageDB, neo4jDriver *neo4j.Driver, grpcContainer *grpc_client.Container) *service.Services {
	commonServices := commonservice.InitServices(&commonconf.GlobalConfig{
		RabbitMQConfig: &s.cfg.RabbitMQConfig,
	}, postgresDb.GormDB, neo4jDriver, s.cfg.Neo4j.Database, grpcContainer, s.log)

	return service.InitServices(
		s.log,
		neo4jDriver,
		postgresDb.GormDB,
		s.cfg,
		commonServices,
		grpcContainer,
		s.cache,
	)
}

func (s *server) setupDependencies(ctx context.Context) error {
	// Initialize postgres db
	postgresDb, err := s.initPostgresDB()
	if err != nil {
		return err
	}

	// Initialize Neo4j
	neo4jDriver, err := s.initNeo4j()
	if err != nil {
		return err
	}
	defer neo4jDriver.Close(ctx)

	// Initialize gRPC
	grpcContainer, err := s.initGRPC()
	if err != nil {
		return errors.Wrap(err, "failed to connect to gRPC")
	}

	// Initialize services
	s.services = s.initializeServices(postgresDb, &neo4jDriver, grpcContainer)

	return nil
}

func (s *server) setupMetrics() {
	registerPrometheusMetrics()
}

func (s *server) startMetricsServer() {
	if s.cfg.ApiPort != s.cfg.MetricsPort {
		mr := gin.Default()
		mr.Use(s.prometheusMiddleware())
		mr.Use(s.bodyLoggerMiddleware)
		mr.GET(s.cfg.Metrics.PrometheusPath, metricsHandler)
		if err := mr.Run(":" + s.cfg.MetricsPort); err != nil {
			s.log.Errorf("Failed to start metrics server: %v", err)
		}
	}
}

func (s *server) setupMiddleware() {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}

	s.engine.Use(cors.New(corsConfig))
	s.engine.Use(ginzap.GinzapWithConfig(s.log.Logger(), &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
		SkipPaths:  []string{"/metrics", "/health", "/readiness", "/"},
	}))
	s.engine.Use(ginzap.RecoveryWithZap(s.log.Logger(), true))
	s.engine.Use(s.prometheusMiddleware())
	s.engine.Use(s.bodyLoggerMiddleware)
}

func (s *server) prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func(start time.Time) {
			// TODO implement metrics COS-314
			// TODO count duration / success / failed requests
		}(time.Now())
		c.Next()
	}
}

func (s *server) bodyLoggerMiddleware(c *gin.Context) {
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	c.Next()
	c.Set("bodyBytes", blw.body.Bytes())
}

func (s *server) startMainServer() error {
	return s.engine.Run(":" + s.cfg.ApiPort)
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func metricsHandler(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}

func registerPrometheusMetrics() {
	// TODO: Implement metrics registration
}
