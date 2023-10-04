package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/constants"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/service"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
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

	// init openline postgres db client
	sqlDb, gormDb, errPostgres := config.NewPostgresClient(cfg)
	if errPostgres != nil {
		appLogger.Fatalf("failed opening connection to postgres: %v", errPostgres.Error())
	}
	defer sqlDb.Close()

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
	if cfg.Service.EventsProcessingPlatformEnabled {
		df := grpc_client.NewDialFactory(cfg, appLogger)
		gRPCconn, err = df.GetEventsProcessingPlatformConn()
		if err != nil {
			appLogger.Fatalf("Failed to connect: %v", err)
		}
		defer df.Close(gRPCconn)
	}

	// Services
	grpcContainer := grpc_client.InitClients(gRPCconn)
	services := service.InitServices(cfg, appLogger, neo4jDriver, gormDb, airbyteStoreDb, grpcContainer)

	services.InitService.Init()

	// Task queues
	var taskQueueSyncCustomerOsData = &TaskQueue{name: "Sync Customer OS Data"}
	var taskQueueSyncToEventStore = &TaskQueue{name: "Sync Neo4j Data to EventStore"}

	go runTaskQueue(appLogger, taskQueueSyncCustomerOsData, cfg.SyncCustomerOsData.TimeoutAfterTaskRun, []func(){
		func() {
			runId, _ := uuid.NewRandom()
			appLogger.Infof("run id: %s syncing data into customer-os at %v", runId.String(), time.Now().UTC())
			services.SyncCustomerOsDataService.Sync(ctx, runId.String())
			appLogger.Infof("run id: %s sync completed at %v", runId.String(), time.Now().UTC())
		},
	})

	if cfg.SyncToEventStore.Enabled {
		syncTasks := []func(){}
		if cfg.SyncToEventStore.Emails.Enabled {
			syncTasks = append(syncTasks, func() {
				ctxWithTimeout, cancel := utils.GetLongLivedContext(context.Background())
				defer cancel()
				services.SyncToEventStoreService.SyncEmails(ctxWithTimeout, syncToEventStoreBatchSize(cfg, cfg.SyncToEventStore.Emails.BatchSize))
				select {
				case <-ctxWithTimeout.Done():
					appLogger.Error("Timeout reached for syncing emails to event store")
				default:
				}
			})
		}
		if cfg.SyncToEventStore.PhoneNumbers.Enabled {
			syncTasks = append(syncTasks, func() {
				ctxWithTimeout, cancel := utils.GetLongLivedContext(context.Background())
				defer cancel()
				services.SyncToEventStoreService.SyncPhoneNumbers(ctxWithTimeout, syncToEventStoreBatchSize(cfg, cfg.SyncToEventStore.PhoneNumbers.BatchSize))
				select {
				case <-ctxWithTimeout.Done():
					appLogger.Error("Timeout reached for syncing phone numbers to event store")
				default:
				}
			})
		}
		if cfg.SyncToEventStore.Locations.Enabled {
			syncTasks = append(syncTasks, func() {
				ctxWithTimeout, cancel := utils.GetLongLivedContext(context.Background())
				defer cancel()
				services.SyncToEventStoreService.SyncLocations(ctxWithTimeout, syncToEventStoreBatchSize(cfg, cfg.SyncToEventStore.Locations.BatchSize))
				select {
				case <-ctxWithTimeout.Done():
					appLogger.Error("Timeout reached for syncing locations to event store")
				default:
				}
			})
		}
		if cfg.SyncToEventStore.Contacts.Enabled {
			syncTasks = append(syncTasks, func() {
				ctxWithTimeout, cancel := utils.GetLongLivedContext(context.Background())
				defer cancel()
				services.SyncToEventStoreService.SyncContacts(ctxWithTimeout, syncToEventStoreBatchSize(cfg, cfg.SyncToEventStore.Contacts.BatchSize))
				select {
				case <-ctxWithTimeout.Done():
					appLogger.Error("Timeout reached for syncing contacts to event store")
				default:
				}
			})
		}
		if cfg.SyncToEventStore.Organizations.Enabled {
			syncTasks = append(syncTasks, func() {
				ctxWithTimeout, cancel := utils.GetLongLivedContext(context.Background())
				defer cancel()
				services.SyncToEventStoreService.SyncOrganizations(ctxWithTimeout, syncToEventStoreBatchSize(cfg, cfg.SyncToEventStore.Organizations.BatchSize))
				services.SyncToEventStoreService.SyncOrganizationsLinksWithDomains(ctxWithTimeout, syncToEventStoreBatchSize(cfg, cfg.SyncToEventStore.Organizations.OrganizationDomainsBatchSize))
				select {
				case <-ctxWithTimeout.Done():
					appLogger.Error("Timeout reached for syncing organizations to event store")
				default:
				}
			})
		}
		go runTaskQueue(appLogger, taskQueueSyncToEventStore, cfg.SyncToEventStore.TimeoutAfterTaskRun, syncTasks)
	}

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

func syncToEventStoreBatchSize(cfg *config.Config, batchSize int) int {
	if batchSize == -1 {
		return cfg.SyncToEventStore.BatchSize
	}
	return batchSize
}
