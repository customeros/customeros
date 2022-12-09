package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracked-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracked-data/config/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracked-data/service"
	"log"
	"sync"
	"time"
)

var (
	taskFunctions      []func()
	taskFunctionsMutex sync.Mutex
)

func AddTask(function func()) {
	defer taskFunctionsMutex.Unlock()
	taskFunctionsMutex.Lock()

	taskFunctions = append(taskFunctions, function)
}

func RunTasks() {
	defer taskFunctionsMutex.Unlock()
	taskFunctionsMutex.Lock()

	for _, t := range taskFunctions {
		t()
	}
}

func init() {
	logger.Logger = logger.New(log.New(log.Default().Writer(), "", log.Ldate|log.Ltime|log.Lmicroseconds), logger.Config{
		Colorful: false,
		LogLevel: logger.Info,
	})
}

func main() {
	cfg := loadConfiguration()

	sqlDb, gormDb, errPostgres := config.NewPostgresClient(cfg)
	if errPostgres != nil {
		log.Fatalf("failed opening connection to postgres: %v", errPostgres.Error())
	}
	defer sqlDb.Close()

	neo4jDriver, errNeo4j := config.NewDriver(cfg)
	if errNeo4j != nil {
		log.Fatalf("failed opening connection to neo4j: %v", errNeo4j.Error())
	}
	defer (*neo4jDriver).Close()

	serviceContainer := service.InitServices(neo4jDriver, gormDb)

	if errPostgres == nil && errNeo4j == nil {
		AddTask(func() {
			runId, _ := uuid.NewRandom()
			logger.Logger.Info("run id: %s syncing tracked data into customer-os at %v", runId.String(), time.Now().UTC())
			result := serviceContainer.SyncService.Sync(runId.String(), cfg.PageViewsBucketSize)
			logger.Logger.Info("run id: %s sync completed at %v, processed %d records", runId.String(), time.Now().UTC(), result)

			if result == 0 {
				timeout := time.Second * time.Duration(cfg.TimeoutAfterTaskRun)
				logger.Logger.Info("waiting %v seconds before next run", timeout.Seconds())
				time.Sleep(timeout)
			}
		})
	}

	go func() {
		for {
			RunTasks()
		}
	}()

	select {}
}

func loadConfiguration() *config.Config {
	if err := godotenv.Load(); err != nil {
		logger.Logger.Warn("Failed loading .env file")
	}

	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		logger.Logger.Warn("%+v", err)
	}

	return &cfg
}
