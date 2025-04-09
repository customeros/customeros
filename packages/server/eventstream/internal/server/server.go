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

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/inbox/internal/config"
	"github.com/customeros/customeros/packages/server/inbox/internal/logger"
	nats_internal "github.com/customeros/customeros/packages/server/inbox/internal/nats"
	"github.com/customeros/customeros/packages/server/inbox/internal/repository"
	"github.com/customeros/customeros/packages/server/inbox/internal/telemetry"
	"github.com/customeros/customeros/packages/server/inbox/services"
)

type Server struct {
	config       *config.Config
	logger       logger.Logger
	natsConn     *nats_internal.NATSConnections
	services     *services.Services
	repositories *repository.Repositories
}

func NewServer(cfg *config.Config, mailstackDB *gorm.DB, warehouseDB *gorm.DB) (*Server, error) {
	// Initialize logger
	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()

	// Initialize OpenTelemetry
	err := telemetry.InitOpenTelemetry(context.Background(), cfg.OpenTelemetry)
	if err != nil {
		log.Printf("Warning: Could not initialize OpenTelemetry: %s", err.Error())
	}

	// Initialize repositories
	repos := repository.InitRepositories(mailstackDB, warehouseDB)
	if err != nil {
		return nil, err
	}

	// Initialize NATS Streams
	natsConn, err := nats_internal.InitNats(cfg.NATSConfig, cfg.AppConfig.Environment)
	if err != nil {
		log.Fatalf("Failed to initialize NATS: %v", err)
	}

	// Initialize services
	svcs := services.InitServices(natsConn, appLogger, repos, cfg)

	return &Server{
		config:       cfg,
		natsConn:     natsConn,
		services:     svcs,
		repositories: repos,
		logger:       appLogger,
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

	log.Println("Inbox is now running. Press Ctrl+C to exit.")

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

	// Close NATS connection
	if s.natsConn != nil {
		log.Println("Closing NATS connection...")
		s.natsConn.Close()
		log.Println("✅ NATS connection closed")
	}

	// Stop services
	log.Println("Stopping services...")
	if err := s.services.Stop(shutdownCtx); err != nil {
		log.Printf("⚠️ Services shutdown error: %v", err)
	} else {
		log.Println("✅ Services stopped successfully")
	}

	return nil
}

func (s *Server) Logger() logger.Logger {
	return s.logger
}

func (s *Server) Services() *services.Services {
	return s.services
}

func (s *Server) Repositories() *repository.Repositories {
	return s.repositories
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
