package main

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/machinebox/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	commsApiConfig "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/routes"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/routes/ContactHub"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	cfg := loadConfiguration()

	// Initialize logger
	appLogger := logger.NewExtendedAppLogger(&cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("comms-api")

	// Tracing
	tracingCloser := initTracing(&cfg.Jaeger, appLogger)
	defer tracingCloser.Close()

	graphqlClient := graphql.NewClient(cfg.Service.CustomerOsAPI)
	redisUrl := fmt.Sprintf("%s://%s", cfg.Redis.Scheme, cfg.Redis.Host)
	log.Printf("redisUrl: %s", redisUrl)
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		log.Fatalf("unvalid redis redisUrl: %s %v", redisUrl, err)
	}
	redisClient := redis.NewClient(opt)

	db, _ := InitDB(&cfg, appLogger)
	defer db.SqlDB.Close()

	neo4jDriver, err := config.NewDriver(&cfg)
	if err != nil {
		appLogger.Fatalf("Could not establish connection with neo4j at: %v, error: %v", cfg.Neo4jConfig.Target, err.Error())
	}
	ctx := context.Background()
	defer neo4jDriver.Close(ctx)

	services := service.InitServices(graphqlClient, redisClient, &cfg, db.GormDB, &neo4jDriver, cfg.Neo4jConfig.Database)

	hub := ContactHub.NewContactHub()
	go hub.Run()

	routes.Run(&cfg, hub, services)

}

func InitDB(cfg *config.Config, appLogger logger.Logger) (db *config.StorageDB, err error) {
	if db, err = config.NewDBConn(cfg); err != nil {
		appLogger.Fatalf("Coud not open db connection: %s", err.Error())
	}
	return
}

func loadConfiguration() commsApiConfig.Config {
	if err := godotenv.Load(); err != nil {
		log.Println("[WARNING] Error loading .env file")
	}

	cfg := commsApiConfig.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}

	return cfg
}

func initTracing(cfg *tracing.JaegerConfig, appLogger logger.Logger) io.Closer {
	tracer, closer, err := tracing.NewJaegerTracer(cfg, appLogger)
	if err != nil {
		appLogger.Fatalf("Could not initialize jaeger tracer: %v", err.Error())
	}
	opentracing.SetGlobalTracer(tracer)
	return closer
}

func initConfig(cfg *commsApiConfig.Config) *gorm.Config {
	return &gorm.Config{
		AllowGlobalUpdate: true,
		Logger:            initLog(cfg),
	}
}

// initLog Connection Log Configuration
func initLog(cfg *commsApiConfig.Config) gormLogger.Interface {
	var logLevel = gormLogger.Silent
	switch cfg.Postgres.LogLevel {
	case "ERROR":
		logLevel = gormLogger.Error
	case "WARN":
		logLevel = gormLogger.Warn
	case "INFO":
		logLevel = gormLogger.Info
	}
	newLogger := gormLogger.New(log.New(io.MultiWriter(os.Stdout), "\r\n", log.LstdFlags), gormLogger.Config{
		Colorful:      true,
		LogLevel:      logLevel,
		SlowThreshold: time.Second,
	})
	return newLogger
}
