package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/service"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"sync"
	"time"
)

type taskQueue struct {
	taskFunctions []func()
	mutex         sync.Mutex
	waitGroup     sync.WaitGroup
}

func (t *taskQueue) AddTask(function func()) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.taskFunctions = append(t.taskFunctions, function)
}

func (t *taskQueue) RunTasks() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if len(t.taskFunctions) == 0 {
		logrus.Warn("No task found, exiting")
		return
	}
	for _, task := range t.taskFunctions {
		t.waitGroup.Add(1)
		go func(fn func()) {
			defer t.waitGroup.Done()
			fn()
		}(task)
	}
	t.taskFunctions = nil
	t.waitGroup.Wait()
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

	ctx := context.Background()
	// init openline neo4j db client
	neo4jDriver, errNeo4j := config.NewDriver(cfg)
	if errNeo4j != nil {
		logrus.Fatalf("failed opening connection to neo4j: %v", errNeo4j.Error())
	}
	defer (*neo4jDriver).Close(ctx)

	// init airbyte postgres db pool
	airbyteStoreDb := config.InitPoolManager(cfg)

	services := service.InitServices(neo4jDriver, gormDb, airbyteStoreDb)

	services.InitService.Init()

	var taskQueue = &taskQueue{}
	for {
		if errPostgres == nil && errNeo4j == nil {
			taskQueue.AddTask(func() {
				runId, _ := uuid.NewRandom()
				logrus.Infof("run id: %s syncing data into customer-os at %v", runId.String(), time.Now().UTC())
				services.SyncService.Sync(ctx, runId.String())
				logrus.Infof("run id: %s sync completed at %v", runId.String(), time.Now().UTC())

				timeout := time.Second * time.Duration(cfg.TimeoutAfterTaskRun)
				logrus.Infof("waiting %v seconds before next run", timeout.Seconds())
				time.Sleep(timeout)
			})
		}
		taskQueue.RunTasks()
	}
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
