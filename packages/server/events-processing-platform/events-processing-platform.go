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
	"sync"
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
	appLogger := initLogger(cfg)

	// Create context and add cancel capability
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	var waitGroup sync.WaitGroup

	// Launch server goroutine
	waitGroup.Add(1)
	go startServer(ctx, cfg, appLogger, &waitGroup)

	// Propagate cancel signal
	go handleSignals(cancel, appLogger)

	// Wait for everything to exit
	waitGroup.Wait()

	// Flush logs and exit
	appLogger.Sync()
}

func startServer(ctx context.Context, cfg *config.Config, logger *logger.ExtendedLogger, waitGroup *sync.WaitGroup) {
	// Create the server
	srv := server.NewServer(cfg, logger)

	// Start it in the background
	go func() {
		if err := srv.Start(ctx); err != nil {
			logger.Fatal(err)
		}
		waitGroup.Done()
	}()

	// Return so main can continue
	return
}

func handleSignals(cancel context.CancelFunc, appLogger *logger.ExtendedLogger) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-sigChan:
			appLogger.Info("Interrupt signal received. Shutting down...")
			cancel()
		}
	}()
}

func initLogger(cfg *config.Config) *logger.ExtendedLogger {
	appLogger := logger.NewExtendedAppLogger(&cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName(server.GetMicroserviceName(cfg))
	return appLogger
}
