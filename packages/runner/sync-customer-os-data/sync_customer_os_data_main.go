package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/constants"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/service"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"sync"
	"time"
)

type TaskQueue struct {
	name          string
	taskFunctions []func()
	mutex         sync.Mutex
	waitGroup     sync.WaitGroup
}

func (t *TaskQueue) AddTask(function func()) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.taskFunctions = append(t.taskFunctions, function)
}

func (t *TaskQueue) RunTasks(log logger.Logger) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if len(t.taskFunctions) == 0 {
		log.Warn("No task found, exiting")
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
	cfg := config.Load()

	// Logging
	appLogger := initLogger(cfg)

	// Setting up tracing
	if cfg.Jaeger.Enabled {
		tracer, closer, err := tracing.NewJaegerTracer(&cfg.Jaeger, appLogger)
		if err != nil {
			appLogger.Fatalf("Could not initialize jaeger tracer: %s", err.Error())
		}
		defer closer.Close()
		opentracing.SetGlobalTracer(tracer)
	}

	db, errPostgres := config.NewDBConn(cfg)
	if errPostgres != nil {
		logrus.Fatalf("Coud not open db connection: %s", errPostgres.Error())
		return
	}
	defer db.SqlDB.Close()

	ctx := context.Background()

	// Neo4j DB
	neo4jDriver, errNeo4j := config.NewDriver(appLogger, cfg)
	if errNeo4j != nil {
		appLogger.Fatalf("failed opening connection to neo4j: %v", errNeo4j.Error())
	}
	defer (*neo4jDriver).Close(ctx)

	// Airbyte DB
	airbyteStoreDb := config.InitPoolManager(cfg)

	// gRPC
	var gRPCconn *grpc.ClientConn
	var err error
	if cfg.GrpcClientConfig.EventsProcessingPlatformEnabled {
		df := grpc_client.NewDialFactory(&cfg.GrpcClientConfig)
		gRPCconn, err = df.GetEventsProcessingPlatformConn()
		if err != nil {
			appLogger.Fatalf("Failed to connect: %v", err)
		}
		defer df.Close(gRPCconn)
	}

	// Services
	grpcContainer := grpc_client.InitClients(gRPCconn)
	services := service.InitServices(cfg, appLogger, neo4jDriver, db.GormDB, airbyteStoreDb, grpcContainer)

	services.InitService.Init()

	// Task queues
	var taskQueueSyncCustomerOsData = &TaskQueue{name: "Sync Customer OS Data"}

	go runTaskQueue(appLogger, taskQueueSyncCustomerOsData, cfg.SyncCustomerOsData.TimeoutAfterTaskRun, []func(){
		func() {
			runId, _ := uuid.NewRandom()
			appLogger.Infof("run id: %s syncing data into customer-os at %v", runId.String(), time.Now().UTC())
			services.SyncCustomerOsDataService.Sync(ctx, runId.String())
			appLogger.Infof("run id: %s sync completed at %v", runId.String(), time.Now().UTC())
		},
	})

	select {}

	// Flush logs and exit
	appLogger.Sync()
}

func runTaskQueue(log logger.Logger, taskQueue *TaskQueue, timeoutAfterTaskRun int, taskFuncs []func()) {
	for {
		for _, task := range taskFuncs {
			taskQueue.AddTask(task)
		}

		taskQueue.RunTasks(log)

		// Cooldown a fixed amount of time before running the tasks again
		timeout := time.Second * time.Duration(timeoutAfterTaskRun)
		log.Infof("waiting %v seconds before next run for %s", timeout.Seconds(), taskQueue.name)
		time.Sleep(timeout)
	}
}

func initLogger(cfg *config.Config) logger.Logger {
	appLogger := logger.NewExtendedAppLogger(&cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName(constants.ServiceName)
	return appLogger
}
