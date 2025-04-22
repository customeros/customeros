package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	commonConfig "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	services "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/customeros/customeros/packages/server/core-crm/internal/config"
	"github.com/customeros/customeros/packages/server/core-crm/internal/cron"
	"github.com/customeros/customeros/packages/server/core-crm/internal/repository"
	nats_internal "github.com/customeros/customeros/packages/server/core-crm/nats"
)

type Server struct {
	config                *config.Config
	logger                logger.Logger
	natsConn              *nats_internal.NATSConnections
	cronMgr               *cron.CronManager
	services              *services.CommonServices
	postgresRepositories  *postgres_repository.Repositories
	neo4jRepositories     *neo4j_repository.Repositories
	warehouseRepositories *repository.Repositories
}

func NewServer(cfg *config.Config, warehouseDB *gorm.DB) (*Server, error) {
	// Initialize logger
	appLogger := logger.NewAppLogger(&cfg.CommonConfig.Infrastructure.LoggerConfig)
	appLogger.InitLogger()
	appLogger.WithName("core-crm")

	// Initialize OpenTelemetry
	err := telemetry.InitOpenTelemetry(context.Background(), &cfg.CommonConfig.Infrastructure.OpenTelemetryConfig)
	if err != nil {
		log.Printf("Warning: Could not initialize OpenTelemetry: %s", err.Error())
	}

	// Initialize DBs
	openlineDB, err := commonConfig.InitPostgres(&commonConfig.CommonConfig{
		Infrastructure: commonConfig.InfrastructureConfig{
			PostgresConfig: cfg.CommonConfig.Infrastructure.PostgresConfig,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to postgres: %w", err)
	}

	neoDriver, err := commonConfig.NewNeo4jDriver(cfg.CommonConfig.Infrastructure.Neo4jConfig)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to neo4j: %w", err)
	}

	// Initialize repositories
	postgresRepos := postgres_repository.InitRepositories(openlineDB)
	if err != nil {
		return nil, err
	}

	neo4jRepos := neo4j_repository.InitNeo4jRepositories(&neoDriver, cfg.CommonConfig.Infrastructure.Neo4jConfig.Database)
	if err != nil {
		return nil, err
	}

	// Initialize NATS Streams
	natsConn, err := nats_internal.InitNats(cfg.NATSConfig, cfg.AppConfig.Environment)
	if err != nil {
		log.Fatalf("Failed to initialize NATS: %v", err)
	}

	// Initialize services
	services := services.InitCommonServices(appLogger, neo4jRepos, postgresRepos, cfg.CommonConfig, natsConn, nil)

	// Try to get Kubernetes config
	var k8sClient kubernetes.Interface
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Printf("Not running in Kubernetes cluster: %v", err)
	} else {
		k8sClient, err = kubernetes.NewForConfig(k8sConfig)
		if err != nil {
			log.Printf("Failed to create kubernetes client: %v", err)
		}
	}

	// Initialize and start cron manager
	cronManager := cron.NewCronManager(
		cfg,
		appLogger,
		k8sClient,
	)

	// If running in Kubernetes, use leader election
	if k8sClient != nil {
		podName := os.Getenv("POD_NAME")
		if podName == "" {
			log.Fatal("POD_NAME environment variable not set")
		}
		namespace := os.Getenv("POD_NAMESPACE")
		if namespace == "" {
			log.Fatal("POD_NAMESPACE environment variable not set")
		}

		go func() {
			if err := cronManager.Start(podName, namespace); err != nil {
				log.Fatalf("Failed to start cron manager: %v", err)
			}
		}()
	} else {
		// Local development - start cron manager directly
		log.Println("Running in local mode - starting cron manager without leader election")
		go func() {
			cronManager.StartCron()
		}()
	}

	return &Server{
		config:               cfg,
		natsConn:             natsConn,
		cronMgr:              cronManager,
		services:             services,
		postgresRepositories: postgresRepos,
		neo4jRepositories:    neo4jRepos,
		logger:               appLogger,
	}, nil
}

func (s *Server) Run() error {
	// Create root context for the application
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Starting services
	log.Println("Starting services...")
	if err := s.services.Start(ctx); err != nil {
		return fmt.Errorf("failed to start services: %w", err)
	}
	log.Println("✅ Services started successfully")

	return s.waitForShutdown()
}

func (s *Server) waitForShutdown() error {
	defer s.recoverWithTracing("shutdown")

	// Set up signal handling for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Wait for termination signal
	<-stop
	log.Println("Shutting down...")

	// Create a context with timeout for shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	// Stop cron manager
	s.cronMgr.Stop()
	log.Println("Shutdown complete")

	// Close NATS connection
	if s.natsConn != nil {
		log.Println("Closing NATS connection...")
		s.natsConn.Close()
		log.Println("✅ NATS connection closed")
	}

	// Stop services
	log.Println("Stopping services...")
	s.services.Stop(shutdownCtx)
	log.Println("✅ Services stopped")

	return nil
}

func (s *Server) recoverWithTracing(name string) {
	if r := recover(); r != nil {
		// Create a new span for the panic
		span := opentracing.GlobalTracer().StartSpan(
			fmt.Sprintf("panic.%s", name),
		)
		defer span.Finish()

		// Mark span as failed
		ext.Error.Set(span, true)

		// Log panic details
		span.LogKV(
			"event", "panic",
			"process", name,
			"error", fmt.Sprintf("%v", r),
			"stack", string(debug.Stack()),
		)

		log.Printf("❌ Panic in %s: %v\n%s", name, r, debug.Stack())
	}
}

func (s *Server) wrapGoroutine(name string, fn func()) {
	defer s.recoverWithTracing(name)
	fn()
}
