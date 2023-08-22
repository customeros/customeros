package main

import (
	"context"
	"flag"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	flag.Parse()

	// Initialize configuration
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize logger
	appLogger := logger.NewExtendedAppLogger(&cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName(server.GetMicroserviceName(cfg))

	// Create a context with a cancel function
	ctx, cancel := context.WithCancel(context.Background())

	// Set up signal handler to cancel context on interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-sigChan:
			appLogger.Info("Interrupt signal received. Shutting down...")
			cancel()
		case <-ctx.Done():
			// Do nothing
		}
	}()

	// Init server
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.NewServer(cfg, appLogger).Run(ctx)
	}()

	// Wait for server to exit or context to be canceled
	select {
	case err := <-errChan:
		appLogger.Fatalf("Server error: %v", err)
	case <-ctx.Done():
		// Do nothing
	}

	// Flush logs and exit
	appLogger.Sync()
}
