package main

import (
	"context"
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
	processor "github.com/customeros/customeros/packages/server/mailstack/services/email_processor"
)

func main() {
	// Configure logging to include timestamps and file information
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("MailStack starting up...")

	// Initialize configuration
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	if cfg == nil {
		log.Fatal("config is empty")
	}

	// Set up context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup the database
	db, err := database.InitDatabase(cfg.MailstackDatabaseConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize repositories
	repos := repository.InitRepositories(db)

	// Initialize services
	svcs := services.InitServices()

	// Set up webhook handler for email events
	webhookURL := "https://webhook.site/9efaff8f-b23e-4874-9750-e0089cc092ab"
	emailProcessor := processor.NewProcessor(webhookURL)

	// Register webhook handler
	log.Println("Registering event handler...")
	svcs.IMAPService.SetEventHandler(emailProcessor.ProcessMailEvent)

	// Setup mailboxes
	if err := internal.InitMailboxes(svcs); err != nil {
		log.Printf("Failed to initialize mailboxes: %v", err)
	}

	// Start the IMAP service
	log.Println("Starting IMAP service...")
	if err := svcs.IMAPService.Start(ctx); err != nil {
		log.Fatalf("❌ Failed to start IMAP service: %v", err)
	}
	log.Println("✅ IMAP service started successfully")

	// Initialize Gin in release mode for production
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Setup API routes
	api.RegisterRoutes(router, svcs, repos)

	// Create HTTP server with Gin handler
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Println("Starting HTTP server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("❌ HTTP server error: %v", err)
		}
	}()
	log.Println("✅ HTTP server started successfully")

	log.Println("MailStack is now running. Press Ctrl+C to exit.")

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
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("❌ HTTP server shutdown error: %v", err)
	} else {
		log.Println("✅ HTTP server shut down successfully")
	}

	// Stop IMAP service with timeout
	log.Println("Stopping IMAP service...")
	stopDone := make(chan struct{})
	go func() {
		defer close(stopDone)
		if err := svcs.IMAPService.Stop(); err != nil {
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

	log.Println("Shutdown complete")
}
