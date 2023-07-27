package cron

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/integrity-checker/container"
	"github.com/openline-ai/openline-customer-os/packages/runner/integrity-checker/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/integrity-checker/service"
	"github.com/robfig/cron"
	"sync"
)

var jobLock sync.Mutex

func StartCron(cont *container.Container) *cron.Cron {
	c := cron.New()

	err := c.AddFunc(cont.Cfg.Cron.CronScheduleNeo4jIntegrityChecker, func() {
		lockAndRunJob(cont, neo4jIntegrityCheckerJob)
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

func neo4jIntegrityCheckerJob(cont *container.Container) {
	service.NewNeo4jIntegrityCheckerService(cont.Cfg, cont.Log, cont.Repositories).RunIntegrityCheckerQueries()
}
