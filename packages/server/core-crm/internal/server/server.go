package server

import (
	"context"
	"fmt"
	nats_common "github.com/customeros/customeros/packages/server/customer-os-common-module/nats"
	"github.com/customeros/customeros/packages/server/enums"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	commonConfig "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_db "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/database"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/customeros/customeros/packages/server/core-crm/api"
	"github.com/customeros/customeros/packages/server/core-crm/internal/config"
	"github.com/customeros/customeros/packages/server/core-crm/internal/cron"
	"github.com/customeros/customeros/packages/server/core-crm/internal/database"
	"github.com/customeros/customeros/packages/server/core-crm/internal/telemetry"
	"github.com/customeros/customeros/packages/server/core-crm/service"
)

type Server struct {
	config                *config.Config
	logger                logger.Logger
	natsConn              *nats_common.NATSConnections
	httpServer            *http.Server
	router                *gin.Engine
	cronMgr               *cron.CronManager
	services              *service.Services
	postgresRepositories  *postgres_repository.Repositories
	neo4jRepositories     *neo4j_repository.Repositories
	warehouseRepositories *postgres_repository.WarehouseRepositories
}

func NewServer(cfg *config.Config, warehouseDB *database.DatabaseConnection) (*Server, error) {
	// Initialize logger
	appLogger := logger.NewAppLogger(&cfg.CommonConfig.Infrastructure.LoggerConfig)
	appLogger.InitLogger()
	appLogger.WithName("core-crm")

	// Initialize OpenTelemetry
	err := telemetry.InitOpenTelemetry(context.Background(), cfg.Telemetry)
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
	neo4jRepos := neo4j_repository.InitNeo4jRepositories(&neoDriver, cfg.CommonConfig.Infrastructure.Neo4jConfig.Database)
	warehouseRepos := postgres_repository.InitWarehouseRepositories(&postgres_db.DbConnections{
		ReadDB:  warehouseDB.ReadDB,
		WriteDB: warehouseDB.WriteDB,
	})

	// Initialize NATS Streams
	natsConn, err := nats_common.InitNats(&cfg.CommonConfig.Infrastructure.NatsConfig, cfg.AppConfig.Environment, enums.GetAllStreams())
	if err != nil {
		log.Fatalf("Failed to initialize NATS: %v", err)
	}

	// Initialize common services
	services := service.InitServices(appLogger, neo4jRepos, postgresRepos, warehouseRepos, cfg.CommonConfig, natsConn)

	// Initialize Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// register API Routes
	api.RegisterRoutes(context.Background(), router, services.CommonServices, cfg.AppConfig)

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
		services,
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
		config:   cfg,
		natsConn: natsConn,
		router:   router,
		httpServer: &http.Server{
			Addr:    ":" + cfg.AppConfig.APIPort,
			Handler: router,
		},
		cronMgr:               cronManager,
		services:              services,
		postgresRepositories:  postgresRepos,
		neo4jRepositories:     neo4jRepos,
		warehouseRepositories: warehouseRepos,
		logger:                appLogger,
	}, nil
}

func (s *Server) Run() error {
	// Create root context for the application
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Starting services
	log.Println("Starting services...")
	if err := s.services.CommonServices.Start(ctx); err != nil {
		return fmt.Errorf("failed to start services: %w", err)
	}
	log.Println("✅ Services started successfully")

	// Start HTTP server in a goroutine with panic recovery
	go s.wrapGoroutine("http_server", func() {
		err := s.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("❌ HTTP server error: %v", err)
		}
	})
	log.Println("✅ HTTP server started successfully")
	log.Printf("Core CRM is now running and listening on port %s. Press Ctrl+C to exit.", s.httpServer.Addr)
	fmt.Println("")

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

	// Shut down HTTP server
	log.Println("Shutting down HTTP server...")
	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("❌ HTTP server shutdown error: %v", err)
	} else {
		log.Println("✅ HTTP server shut down successfully")
	}

	// Close NATS connection
	if s.natsConn != nil {
		log.Println("Closing NATS connection...")
		s.natsConn.Close()
		log.Println("✅ NATS connection closed")
	}

	// Stop services
	log.Println("Stopping services...")
	s.services.CommonServices.Stop(shutdownCtx)
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
