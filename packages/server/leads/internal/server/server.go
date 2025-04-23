package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/customeros/customeros/packages/server/leads/api/handlers"
	"github.com/customeros/customeros/packages/server/leads/internal/config"
	"github.com/customeros/customeros/packages/server/leads/internal/cron"
	"github.com/customeros/customeros/packages/server/leads/internal/logger"
	nats_internal "github.com/customeros/customeros/packages/server/leads/internal/nats"
	"github.com/customeros/customeros/packages/server/leads/internal/repository"
	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
	"github.com/customeros/customeros/packages/server/leads/services"
)

type Server struct {
	config       *config.Config
	logger       logger.Logger
	httpServer   *http.Server
	natsConn     *nats_internal.NATSConnections
	router       *gin.Engine
	cronMgr      *cron.CronManager
	services     *services.Services
	repositories *repository.Repositories
	apiHandlers  *handlers.APIHandlers
}

func NewServer(cfg *config.Config, leadsDB *gorm.DB, warehouseDB *gorm.DB) (*Server, error) {
	// Initialize logger
	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()

	// Initialize OpenTelemetry
	err := telemetry.InitOpenTelemetry(context.Background(), cfg.OpenTelemetry)
	if err != nil {
		log.Printf("Warning: Could not initialize OpenTelemetry: %s", err.Error())
	}

	// Initialize repositories
	repos := repository.InitRepositories(leadsDB, warehouseDB)
	if err != nil {
		return nil, err
	}

	// Initialize NATS Streams
	natsConn, err := nats_internal.InitNats(cfg.NATSConfig, cfg.AppConfig.Environment)
	if err != nil {
		log.Fatalf("Failed to initialize NATS: %v", err)
	}

	// Initialize services
	svcs := services.InitServices(natsConn, repos, cfg)

	// Initialize Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Initialize API Handlers
	handlers := handlers.InitHandlers(natsConn, repos)

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
		repos,
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
		cronMgr:      cronManager,
		services:     svcs,
		repositories: repos,
		logger:       appLogger,
		apiHandlers:  handlers,
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

	// Start HTTP server in a goroutine with panic recovery
	go s.wrapGoroutine("http_server", func() {
		log.Println("Starting HTTP server")
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("❌ HTTP server error: %v", err)
		}
	})
	log.Println("✅ HTTP server started successfully")
	log.Println("EventStraem is now running. Press Ctrl+C to exit.")

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

	// Shut down HTTP server
	log.Println("Shutting down HTTP server...")
	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("❌ HTTP server shutdown error: %v", err)
	} else {
		log.Println("✅ HTTP server shut down successfully")
	}

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
