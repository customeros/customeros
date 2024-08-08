package server

import (
	"bytes"
	"context"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	commonconf "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/route"
	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/service"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	tracer, closer, err := tracing.NewJaegerTracer(&server.cfg.Jaeger, server.log)
	if err != nil {
		server.log.Fatalf("Could not initialize jaeger tracer: %s", err.Error())
	}
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	registerPrometheusMetrics()

	// Initialize postgres db
	postgresDb, _ := InitDB(server.cfg, server.log)
	defer postgresDb.SqlDB.Close()

	// Setting up Neo4j
	neo4jDriver, err := commonconf.NewNeo4jDriver(server.cfg.Neo4j)
	if err != nil {
		server.log.Fatalf("Could not establish connection with neo4j at: %v, error: %v", server.cfg.Neo4j.Target, err.Error())
	}
	defer neo4jDriver.Close(ctx)

	// Setting up common services
	services := service.InitServices(server.cfg, postgresDb.GormDB, &neo4jDriver, server.log)

	// Setting up Gin
	r := gin.Default()

	// Setting up CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}

	r.Use(cors.New(corsConfig))
	r.Use(ginzap.GinzapWithConfig(server.log.Logger(), &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
		SkipPaths:  []string{"/metrics", "/health", "/"},
	}))
	r.Use(ginzap.RecoveryWithZap(server.log.Logger(), true))
	r.Use(prometheusMiddleware())
	r.Use(bodyLoggerMiddleware)

	// set routes
	route.RegisterRoutes(ctx, r, services)

	r.GET("/health", HealthCheckHandler)
	r.GET("/", RootHandler)

	if server.cfg.ApiPort == server.cfg.MetricsPort {
		r.GET(server.cfg.Metrics.PrometheusPath, metricsHandler)
	} else {
		go func() {
			mr := gin.Default()
			mr.Use(prometheusMiddleware())
			mr.Use(bodyLoggerMiddleware)
			mr.GET(server.cfg.Metrics.PrometheusPath, metricsHandler)
			mr.Run(":" + server.cfg.MetricsPort)
		}()
	}

	r.Run(":" + server.cfg.ApiPort)

	<-server.doneCh
	server.log.Infof("Application %s exited properly", constants.ServiceName)
	return nil
}

func InitDB(cfg *config.Config, log logger.Logger) (db *commonconf.StorageDB, err error) {
	if db, err = commonconf.NewPostgresDBConn(cfg.Postgres); err != nil {
		log.Fatalf("Could not open db connection: %s", err.Error())
	}
	return
}

func HealthCheckHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "OK"})
}

func RootHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Up and running!",
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
			//TODO implement custom metrics
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
