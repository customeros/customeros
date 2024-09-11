package main

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// @title CustomerOS API
// @version 1.0
// @description CustomerOS API for multiple services (Verify, Enrich, Orgs)
// @host api.customeros.ai
// @schemes https
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-CUSTOMER-OS-API-KEY
func init() {
	// This empty function ensures that swagger_info.go is included in the build
}

func main() {
	// Initialize configuration
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize logger
	appLogger := logger.NewExtendedAppLogger(&cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName(constants.ServiceName)

	// Create a context with a cancel function
	ctx, cancel := context.WithCancel(context.Background())

	// Set up signal handler to cancel context on interrupt
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-shutdown:
			appLogger.Warn("Interrupt signal received. Shutting down...")
			cancel()
		case <-ctx.Done():
			appLogger.Warn("Context canceled. Shutting down...")
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
