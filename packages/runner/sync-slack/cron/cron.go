package cron

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/container"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/service"
	"github.com/robfig/cron"
	"sync"
)

var jobLock sync.Mutex

func StartCron(cont *container.Container) *cron.Cron {
	c := cron.New()

	err := c.AddFunc(cont.Cfg.Cron.CronScheduleSyncFromSlack, func() {
		lockAndRunJob(cont, syncSlackRawJob)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job: %v", err.Error())
	}

	c.Start()

	return c
}

func lockAndRunJob(cont *container.Container, job func(cont *container.Container)) {
	jobLock.Lock()
	defer jobLock.Unlock()

	job(cont)
}

func StopCron(log logger.Logger, cron *cron.Cron) error {
	// Gracefully stop
	log.Info("Gracefully stopping cron")
	cron.Stop()
	return nil
}

func syncSlackRawJob(cont *container.Container) {
	service.NewSlackRawService(cont.Cfg, cont.Log, cont.Repositories, cont.RawDataStoreDB).FetchAllTenantsDataFromSlack()
}
