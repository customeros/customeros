package main

import (
	"context"
	"github.com/caarlos0/env/v6"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/service"
	"github.com/sirupsen/logrus"
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

func (t *TaskQueue) RunTasks() {
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

func runTaskQueue(taskQueue *TaskQueue, timeoutAfterTaskRun int, taskFuncs []func()) {
	for {
		for _, task := range taskFuncs {
			taskQueue.AddTask(task)
		}

		taskQueue.RunTasks()

		// Cooldown a fixed amount of time before running the tasks again
		timeout := time.Second * time.Duration(timeoutAfterTaskRun)
		logrus.Infof("waiting %v seconds before next run for %s", timeout.Seconds(), taskQueue.name)
		time.Sleep(timeout)
	}
}

func main() {
	ctx := context.Background()

	cfg := loadConfiguration()
	config.InitLogger(cfg)

	sqlDb, gormDb, errPostgres := config.NewPostgresClient(cfg)
	if errPostgres != nil {
		logrus.Fatalf("failed opening connection to postgres: %v", errPostgres.Error())
	}
	defer sqlDb.Close()

	neo4jDriver, errNeo4j := config.NewDriver(cfg)
	if errNeo4j != nil {
		logrus.Fatalf("failed opening connection to neo4j: %v", errNeo4j.Error())
	}
	defer (*neo4jDriver).Close(ctx)

	services := service.InitServices(neo4jDriver, gormDb, cfg)

	var taskQueueSyncEmails = &TaskQueue{name: "Sync emails from gmail"}

	go runTaskQueue(taskQueueSyncEmails, cfg.SyncData.TimeoutAfterTaskRun, []func(){
		func() {
			runId, _ := uuid.NewRandom()
			logrus.Infof("run id: %s syncing emails from gmail into customer-os at %v", runId.String(), time.Now().UTC())

			//tenants, err := services.TenantService.GetAllTenants(ctx)
			//if err != nil {
			//	panic(err)
			//}
			//
			//for _, tenant := range tenants {
			//
			//	if tenant.Name != "openline" {
			//		continue
			//	}

			usersForTenant, err := services.UserService.GetAllUsersForTenant(ctx, "openline")
			if err != nil {
				panic(err)
			}

			for _, user := range usersForTenant {
				emailForUser, err := services.EmailService.FindEmailForUser("openline", user.Id)
				if err != nil {
					panic(err)
				}
				logrus.Infof("user: %v", user)
				services.EmailService.ReadNewEmailsForUsername("openline", emailForUser.RawEmail)
			}

			//}

			logrus.Infof("run id: %s sync completed at %v", runId.String(), time.Now().UTC())
		},
	})

	select {}

	//job - read all users and trigger email sync per user ( 5 mintues )
	//job - read all new emails for a user and sync them ( 1 minute )
	//1 job per user in thread pool

	//job 1
	//get tenants
	//get users in tenant
	//sync emails

	//get all users from tenant with access enabled to gmail ( all users except blacklisted )
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
