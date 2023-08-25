package main

import (
	"database/sql"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/machinebox/graphql"
	commsApiConfig "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/routes"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/routes/ContactHub"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/service"
	commonRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	config := loadConfiguration()

	// Initialize logger
	appLogger := logger.NewExtendedAppLogger(&config.Logger)
	appLogger.InitLogger()
	appLogger.WithName("comms-api")

	// Tracing
	tracingCloser := initTracing(&config.Jaeger, appLogger)
	defer tracingCloser.Close()

	graphqlClient := graphql.NewClient(config.Service.CustomerOsAPI)
	redisUrl := fmt.Sprintf("%s://%s", config.Redis.Scheme, config.Redis.Host)
	log.Printf("redisUrl: %s", redisUrl)
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		log.Fatalf("unvalid redis redisUrl: %s %v", redisUrl, err)
	}
	redisClient := redis.NewClient(opt)

	db, err := InitDB(&config)
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}
	commonRepositoryContainer := commonRepository.InitRepositories(db.GormDB, nil)
	services := service.InitServices(graphqlClient, redisClient, &config, db)
	hub := ContactHub.NewContactHub()
	go hub.Run()
	routes.Run(&config, hub, services, commonRepositoryContainer) // run this as a background goroutine

}

func InitDB(cfg *commsApiConfig.Config) (db *commsApiConfig.StorageDB, err error) {
	connectString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s ", cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Db, cfg.Postgres.User, cfg.Postgres.Password)
	gormDb, err := gorm.Open(postgres.Open(connectString), initConfig(cfg))

	var sqlDb *sql.DB
	if err != nil {
		return nil, err
	}
	if sqlDb, err = gormDb.DB(); err != nil {
		return nil, err
	}
	if err = sqlDb.Ping(); err != nil {
		return nil, err
	}

	sqlDb.SetMaxIdleConns(cfg.Postgres.MaxIdleConn)
	sqlDb.SetMaxOpenConns(cfg.Postgres.MaxConn)
	sqlDb.SetConnMaxLifetime(time.Duration(cfg.Postgres.ConnMaxLifetime) * time.Second)

	return &commsApiConfig.StorageDB{
		SqlDB:  sqlDb,
		GormDB: gormDb,
	}, nil
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
