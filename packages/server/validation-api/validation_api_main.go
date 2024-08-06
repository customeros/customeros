package main

import (
	"context"
	commonconf "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/service"
	"github.com/opentracing/opentracing-go"
	"io"

	"github.com/caarlos0/env/v6"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/config"
	"log"
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

	ctx := context.Background()

	// Initialize postgres db
	postgresDb, _ := InitDB(cfg, appLogger)
	defer postgresDb.SqlDB.Close()

	// Setting up Neo4j
	neo4jDriver, err := commonconf.NewNeo4jDriver(cfg.Neo4j)
	if err != nil {
		appLogger.Fatalf("Could not establish connection with neo4j at: %v, error: %v", cfg.Neo4j.Target, err.Error())
	}
	defer neo4jDriver.Close(ctx)

	// Setting up Gin
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	services := service.InitServices(cfg, postgresDb.GormDB, &neo4jDriver, appLogger)

	rest.RegisterRoutes(ctx, r, services, cfg, appLogger)

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

func InitDB(cfg *config.Config, log logger.Logger) (db *commonconf.StorageDB, err error) {
	if db, err = commonconf.NewPostgresDBConn(cfg.Postgres); err != nil {
		log.Fatalf("Could not open db connection: %s", err.Error())
	}
	return
}
