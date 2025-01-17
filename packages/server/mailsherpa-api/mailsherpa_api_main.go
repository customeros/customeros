package main

import (
	"context"
	"io"

	commonconf "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/route"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/service"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"

	"log"

	"github.com/caarlos0/env/v6"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/config"
)

func main() {
	cfg := loadConfiguration()

	// Initialize Logging
	appLogger := initLogger(cfg)

	// Initialize Tracing
	tracingCloser := initTracing(cfg, appLogger)
	if tracingCloser != nil {
		defer tracingCloser.Close()
	}
	defer tracing.RecoverAndLogToJaeger(appLogger)

	ctx := context.Background()

	// Initialize postgres db
	postgresDb, err := commonconf.InitPostgres(&commonconf.CommonConfig{
		PostgresConfig:      &cfg.PostgresConfig,
		PostgresAsyncConfig: &cfg.PostgresAsyncConfig,
	})
	if err != nil {
		logrus.Fatalf("failed opening connection to postgres: %v", err.Error())
	}
	defer postgresDb.Close()

	// Setting up Neo4j
	neo4jDriver, err := commonconf.NewNeo4jDriver(*cfg.CommonConfig.Neo4jConfig)
	if err != nil {
		appLogger.Fatalf("Could not establish connection with neo4j at: %v, error: %v", cfg.CommonConfig.Neo4jConfig.Target, err.Error())
	}
	defer neo4jDriver.Close(ctx)

	// Setting up Gin
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	postgresRepositories := repository.InitRepositories(postgresDb)

	services := service.InitServices(
		appLogger,
		postgresRepositories,
		cfg.CommonConfig,
	)

	route.RegisterRoutes(ctx, r, services, cfg, appLogger)

	r.Run(":" + cfg.ApiPort)
}

func loadConfiguration() *config.Config {
	if err := godotenv.Load(); err != nil {
		log.Println("[WARNING] Error loading .env file")
	}

	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}

	return &cfg
}

func initLogger(cfg *config.Config) logger.Logger {
	appLogger := logger.NewExtendedAppLogger(&cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("VALIDATION-API")
	return appLogger
}

func initTracing(cfg *config.Config, appLogger logger.Logger) io.Closer {
	if cfg.Jaeger.Enabled {
		tracer, closer, err := tracing.NewJaegerTracer(&cfg.Jaeger, appLogger)
		if err != nil {
			appLogger.Fatalf("Could not initialize jaeger tracer: %v", err.Error())
		}
		opentracing.SetGlobalTracer(tracer)
		return closer
	}
	return nil
}
