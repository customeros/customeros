package main

import (
	"context"
	"io"
	"log"

	"github.com/caarlos0/env/v6"
	commonconf "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"

	"github.com/customeros/customeros/packages/server/mailsherpa-api/config"
	"github.com/customeros/customeros/packages/server/mailsherpa-api/logger"
	"github.com/customeros/customeros/packages/server/mailsherpa-api/route"
	"github.com/customeros/customeros/packages/server/mailsherpa-api/service"
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
		Infrastructure: commonconf.InfrastructureConfig{
			PostgresConfig:      cfg.PostgresConfig,
			PostgresAsyncConfig: cfg.PostgresAsyncConfig,
		},
	})
	if err != nil {
		logrus.Fatalf("failed opening connection to postgres: %v", err.Error())
	}
	defer postgresDb.Close()

	// Setting up Gin
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	postgresRepositories := postgres_repository.InitRepositories(postgresDb)

	services := service.InitServices(
		appLogger,
		postgresRepositories,
		cfg,
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
