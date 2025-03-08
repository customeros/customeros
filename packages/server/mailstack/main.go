package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/customeros/customeros/packages/server/mailstack/internal/database"
	"github.com/customeros/customeros/packages/server/mailstack/internal/models"
	"github.com/customeros/customeros/packages/server/mailstack/internal/repository"
	"github.com/customeros/customeros/packages/server/mailstack/services/imap"
	"github.com/customeros/customeros/packages/server/mailstack/services/mailbox"
)

func main() {
	// Parse command-line flags
	listenAddr := flag.String("addr", ":8080", "HTTP listen address")
	flag.Parse()

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup the database
	dbConfig := database.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "mailstack_user",
		Password: "your_secure_password",
		DBName:   "mailstack",
		SSLMode:  "disable", // Use "require" in production
	}

	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the db schema
	if err := db.AutoMigrate(&models.Mailbox{}, &models.MessageState{}); err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	// Create repositories
	mailboxRepo := repository.NewMailboxRepository(db)
	messageStateRepo := repository.NewMessageStateRepository(db)

	// Create Services
	imapService := imap.NewIMAPService()
	mailService := mailbox.NewMailService(imapService)

	// Set up HTTP server
	mux := http.NewServeMux()

	// Create and set up API server
	mailServer := api.NewMailServer(mailService)
	mailServer.SetupRoutes(mux)

	// Create HTTP server
	server := &http.Server{
		Addr:    *listenAddr,
		Handler: mux,
	}

	// Start mail service
	if err := mailService.Start(ctx); err != nil {
		log.Fatalf("Failed to start mail service: %v", err)
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("Starting HTTP server on %s", *listenAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Set up signal handling for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signal
	<-stop
	log.Println("Shutting down...")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Shut down HTTP server
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	// Stop mail service
	if err := mailService.Stop(); err != nil {
		log.Printf("Mail service shutdown error: %v", err)
	}

	log.Println("Shutdown complete")
}
