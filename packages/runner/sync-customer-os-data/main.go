package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/service"
	"log"
	"sync"
	"time"
)

var (
	taskFunctions      []func()
	taskFunctionsMutex sync.Mutex
	services           *service.Services
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

func main() {
	cfg := loadConfiguration()

	// init openline postgres db client
	sqlDb, gormDb, errPostgres := config.NewPostgresClient(cfg)
	if errPostgres != nil {
		log.Fatalf("failed opening connection to postgres: %v", errPostgres.Error())
	}
	defer sqlDb.Close()

	// init openline neo4j db client
	neo4jDriver, errNeo4j := config.NewDriver(cfg)
	if errNeo4j != nil {
		log.Fatalf("failed opening connection to neo4j: %v", errNeo4j.Error())
	}
	defer (*neo4jDriver).Close()

	// init airbyte postgres db pool
	airbyteStoreDb := config.InitPoolManager(cfg)

	services = service.InitServices(neo4jDriver, gormDb, airbyteStoreDb)

	services.InitService.Init()

	if errPostgres == nil && errNeo4j == nil {
		AddTask(func() {
			runId, _ := uuid.NewRandom()
			log.Printf("run id: %s syncing data into customer-os at %v", runId.String(), time.Now().UTC())
			services.SyncService.Sync()
			log.Printf("run id: %s sync completed at %v", runId.String(), time.Now().UTC())

			timeout := time.Second * time.Duration(cfg.TimeoutAfterTaskRun)
			log.Printf("waiting %v seconds before next run", timeout.Seconds())
			time.Sleep(timeout)
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
		log.Print("Failed loading .env file")
	}

	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v", err)
	}

	return &cfg
}
