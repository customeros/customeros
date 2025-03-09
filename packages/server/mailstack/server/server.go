package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/customeros/customeros/packages/server/mailstack/api"
	"github.com/customeros/customeros/packages/server/mailstack/config"
	"github.com/customeros/customeros/packages/server/mailstack/internal"
	"github.com/customeros/customeros/packages/server/mailstack/internal/database"
	"github.com/customeros/customeros/packages/server/mailstack/internal/repository"
	"github.com/customeros/customeros/packages/server/mailstack/services"
	"github.com/customeros/customeros/packages/server/mailstack/services/email_processor"
)

type Server struct {
	config         *config.Config
	httpServer     *http.Server
	router         *gin.Engine
	services       *services.Services
	repositories   *repository.Repositories
	emailProcessor *email_processor.Processor
}

const NEW_EMAIL_WEBHOOK = "https://webhook.site/9efaff8f-b23e-4874-9750-e0089cc092ab"

func NewServer() (*Server, error) {
	cfg, err := config.InitConfig()
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		err := errors.New("config is empty")
		return nil, err
	}

	// Setup the databases
	mailstackDB, err := database.InitMailstackDatabase(&database.DatabaseConfig{
		DBName:          cfg.MailstackDatabaseConfig.DBName,
		Host:            cfg.MailstackDatabaseConfig.Host,
		Port:            cfg.MailstackDatabaseConfig.Port,
		User:            cfg.MailstackDatabaseConfig.User,
		Password:        cfg.MailstackDatabaseConfig.Password,
		MaxConn:         cfg.MailstackDatabaseConfig.MaxConn,
		MaxIdleConn:     cfg.MailstackDatabaseConfig.MaxIdleConn,
		ConnMaxLifetime: cfg.MailstackDatabaseConfig.ConnMaxLifetime,
		LogLevel:        cfg.MailstackDatabaseConfig.LogLevel,
		SSLMode:         cfg.MailstackDatabaseConfig.SSLMode,
	})
	if err != nil {
		return nil, err
	}

	openlineDB, err := database.InitMailstackDatabase(&database.DatabaseConfig{
		DBName:          cfg.OpenlineDatabaseConfig.DBName,
		Host:            cfg.OpenlineDatabaseConfig.Host,
		Port:            cfg.OpenlineDatabaseConfig.Port,
		User:            cfg.OpenlineDatabaseConfig.User,
		Password:        cfg.OpenlineDatabaseConfig.Password,
		MaxConn:         cfg.OpenlineDatabaseConfig.MaxConn,
		MaxIdleConn:     cfg.OpenlineDatabaseConfig.MaxIdleConn,
		ConnMaxLifetime: cfg.OpenlineDatabaseConfig.ConnMaxLifetime,
		LogLevel:        cfg.OpenlineDatabaseConfig.LogLevel,
		SSLMode:         cfg.OpenlineDatabaseConfig.SSLMode,
	})
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	repos := repository.InitRepositories(mailstackDB, openlineDB)

	// Initialize services
	svcs := services.InitServices()

	// Set up webhook handler for email events
	emailProcessor := email_processor.NewProcessor(NEW_EMAIL_WEBHOOK)

	// Initialize Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	return &Server{
		config:         cfg,
		router:         router,
		services:       svcs,
		repositories:   repos,
		emailProcessor: emailProcessor,
		httpServer: &http.Server{
			Addr:    ":" + cfg.AppConfig.APIPort,
			Handler: router,
		},
	}, nil
}

func (s *Server) Initialize(ctx context.Context) error {
	// Register webhook handler
	log.Println("Registering event handler...")
	s.services.IMAPService.SetEventHandler(s.emailProcessor.ProcessMailEvent)

	// Setup mailboxes
	if err := internal.InitMailboxes(s.services); err != nil {
		return err
	}

	// Setup API routes
	api.RegisterRoutes(ctx, s.router, s.services, s.repositories)

	return nil
}

func (s *Server) Run() error {
	// Create root context for the application
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize server components
	if err := s.Initialize(ctx); err != nil {
		return err
	}

	// Start the IMAP service
	log.Println("Starting IMAP service...")
	if err := s.services.IMAPService.Start(ctx); err != nil {
		return err
	}
	log.Println("✅ IMAP service started successfully")

	// Start HTTP server in a goroutine
	go func() {
		log.Println("Starting HTTP server")
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("❌ HTTP server error: %v", err)
		}
	}()
	log.Println("✅ HTTP server started successfully")
	log.Println("MailStack is now running. Press Ctrl+C to exit.")

	return s.waitForShutdown()
}

func (s *Server) waitForShutdown() error {
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

	// Stop IMAP service with timeout
	log.Println("Stopping IMAP service...")
	stopDone := make(chan struct{})
	go func() {
		defer close(stopDone)
		if err := s.services.IMAPService.Stop(); err != nil {
			log.Printf("❌ IMAP service shutdown error: %v", err)
		} else {
			log.Println("✅ IMAP service stopped successfully")
		}
	}()

	// Wait for IMAP service to stop with timeout
	select {
	case <-stopDone:
		log.Println("IMAP service stopped gracefully")
	case <-time.After(10 * time.Second):
		log.Println("⚠️ IMAP service stop timed out, forcing exit")
	}

	return nil
}
