package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/service"
	"github.com/sirupsen/logrus"
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
	config.InitLogger(cfg)

	// init openline postgres db client
	sqlDb, gormDb, errPostgres := config.NewPostgresClient(cfg)
	if errPostgres != nil {
		logrus.Fatalf("failed opening connection to postgres: %v", errPostgres.Error())
	}
	defer sqlDb.Close()

	// init openline neo4j db client
	neo4jDriver, errNeo4j := config.NewDriver(cfg)
	if errNeo4j != nil {
		logrus.Fatalf("failed opening connection to neo4j: %v", errNeo4j.Error())
	}
	defer (*neo4jDriver).Close()

	// init airbyte postgres db pool
	airbyteStoreDb := config.InitPoolManager(cfg)

	services = service.InitServices(neo4jDriver, gormDb, airbyteStoreDb)

	services.InitService.Init()

	if errPostgres == nil && errNeo4j == nil {
		AddTask(func() {
			runId, _ := uuid.NewRandom()
			logrus.Infof("run id: %s syncing data into customer-os at %v", runId.String(), time.Now().UTC())
			services.SyncService.Sync()
			logrus.Infof("run id: %s sync completed at %v", runId.String(), time.Now().UTC())

			timeout := time.Second * time.Duration(cfg.TimeoutAfterTaskRun)
			logrus.Infof("waiting %v seconds before next run", timeout.Seconds())
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
		logrus.Warn("Failed loading .env file")
	}

	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		logrus.Errorf("%+v", err)
	}

	return &cfg
}
