// main.go
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

	// Create IMAP service
	imapService := imap.NewIMAPService()

	// Create mail service
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
