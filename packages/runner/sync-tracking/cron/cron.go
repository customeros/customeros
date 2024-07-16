package cron

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracking/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/robfig/cron"
	"sync"
)

const (
	identifyTrackingRecordsGroup = "identifyTrackingRecordsGroup"
)

var jobLocks = struct {
	sync.Mutex
	locks map[string]*sync.Mutex
}{
	locks: map[string]*sync.Mutex{
		identifyTrackingRecordsGroup: {},
	},
}

func StartCron(services *service.Services) *cron.Cron {
	c := cron.New()

	// Add jobs
	err := c.AddFunc(CronSchedule, func() {
		lockAndRunJob(services, identifyTrackingRecordsGroup, identifyTrackingRecords)
	})
	if err != nil {
		services.Logger.Fatalf("Could not add cron job %s: %v", "identifyTrackingRecords", err.Error())
	}

	c.Start()

	return c
}

func lockAndRunJob(services *service.Services, groupName string, job func(services *service.Services)) {
	jobLocks.locks[groupName].Lock()
	defer jobLocks.locks[groupName].Unlock()

	job(services)
}

func StopCron(log logger.Logger, cron *cron.Cron) error {
	// Gracefully stop
	log.Info("Gracefully stopping cron")
	cron.Stop()
	return nil
}

func identifyTrackingRecords(services *service.Services) {
	services.EnrichDetailsTrackingService.IdentifyTrackingRecords(context.Background())
}
