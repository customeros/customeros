package main

import (
	"context"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/routes"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/service"
	"github.com/opentracing/opentracing-go"
	"io"
	"log"
)

func InitDB(cfg *config.Config, appLogger logger.Logger) (db *config.StorageDB, err error) {
	if db, err = config.NewDBConn(cfg); err != nil {
		appLogger.Fatalf("Coud not open db connection: %s", err.Error())
	}
	return
}

func main() {
	cfg := loadConfiguration()

	appLogger := logger.NewExtendedAppLogger(&cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("user-admin-api")

	// Initialize Tracing
	tracingCloser := initTracing(cfg, appLogger)
	if tracingCloser != nil {
		defer tracingCloser.Close()
	}

	db, _ := InitDB(cfg, appLogger)
	defer db.SqlDB.Close()

	neo4jDriver, err := config.NewDriver(cfg)
	if err != nil {
		appLogger.Fatalf("Could not establish connection with neo4j at: %v, error: %v", cfg.Neo4j.Target, err.Error())
	}
	ctx := context.Background()
	defer neo4jDriver.Close(ctx)

	// Setting up gRPC client
	df := grpc_client.NewDialFactory(&cfg.GrpcClientConfig)
	gRPCconn, err := df.GetEventsProcessingPlatformConn()
	if err != nil {
		panic(err)
	}
	defer df.Close(gRPCconn)
	grpcContainer := grpc_client.InitClients(gRPCconn)

	appCache := caches.NewCache()
	services := service.InitServices(cfg, db.GormDB, &neo4jDriver, grpcContainer, appCache, appLogger)

	//init app cache
	personalEmailProviderEntities, err := services.CommonServices.PostgresRepositories.PersonalEmailProviderRepository.GetPersonalEmailProviders()
	if err != nil {
		appLogger.Fatalf("Error getting personal email providers: %s", err.Error())
	}
	personalEmailProviders := make([]string, 0)
	for _, personalEmailProvider := range personalEmailProviderEntities {
		personalEmailProviders = append(personalEmailProviders, personalEmailProvider.ProviderDomain)
	}
	appCache.SetPersonalEmailProviders(personalEmailProviders)

	routes.Run(cfg, services)
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
