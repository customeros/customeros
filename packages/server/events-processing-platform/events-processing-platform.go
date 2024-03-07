package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/server"
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

	// Start a heartbeat
	done := make(chan interface{})
	defer close(done)
	const timeout = time.Second
	heartbeat := Heartbeat(done, timeout)
	go logHeartbeat(heartbeat, appLogger)

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

func Heartbeat(
	done <-chan interface{},
	pulseInterval time.Duration,
	nums ...int,
) <-chan int {
	heartbeat := make(chan int, 1)
	go func() {
		defer close(heartbeat)

		time.Sleep(2 * time.Second)

		pulse := time.Tick(pulseInterval)
		for {
			select {
			case <-done:
				return
			case <-pulse:
				select {
				case heartbeat <- 1:
				default:
				}
			}
		}
	}()

	return heartbeat
}

func logHeartbeat(heartbeat <-chan int, logger *logger.ExtendedLogger) {
	for {
		if _, ok := <-heartbeat; !ok {
			return
		}
		logger.Debug("pulse")
	}
}

func initLogger(cfg *config.Config) *logger.ExtendedLogger {
	appLogger := logger.NewExtendedAppLogger(&cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName(server.GetMicroserviceName(cfg))
	return appLogger
}
