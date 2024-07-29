package cron

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracking/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracking/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/robfig/cron"
	"sync"
)

const (
	processNewRecordsGroup                  = "processNewRecordsGroup"
	processIPDataRequestsGroup              = "processIPDataRequestsGroup"
	processIPDataResponsesGroup             = "processIPDataResponsesGroup"
	identifyTrackingRecordsGroup            = "identifyTrackingRecordsGroup"
	createOrganizationsFromTrackedDataGroup = "createOrganizationsFromTrackedDataGroup"
	notifyOnSlackGroup                      = "notifyOnSlackGroup"
)

var jobLocks = struct {
	sync.Mutex
	locks map[string]*sync.Mutex
}{
	locks: map[string]*sync.Mutex{
		processNewRecordsGroup:                  {},
		processIPDataRequestsGroup:              {},
		processIPDataResponsesGroup:             {},
		identifyTrackingRecordsGroup:            {},
		createOrganizationsFromTrackedDataGroup: {},
		notifyOnSlackGroup:                      {},
	},
}

func StartCron(cfg *config.Config, services *service.Services) *cron.Cron {
	c := cron.New()

	// Add jobs
	err := c.AddFunc(cfg.Cron.CronScheduleProcessNewRecords, func() {
		lockAndRunJob(services, processNewRecordsGroup, processNewRecords) // 500 records processed
	})
	if err != nil {
		services.Logger.Fatalf("Could not add cron job %s: %v", "processNewRecords", err.Error())
	}
	err = c.AddFunc(cfg.Cron.CronScheduleProcessIPDataRequests, func() {
		lockAndRunJob(services, processIPDataRequestsGroup, processIPDataRequests) // sending 500 requests
	})
	if err != nil {
		services.Logger.Fatalf("Could not add cron job %s: %v", "processIPDataRequests", err.Error())
	}
	err = c.AddFunc(cfg.Cron.CronScheduleProcessIPDataResponses, func() {
		lockAndRunJob(services, processIPDataResponsesGroup, processIPDataResponses) // 500 tracking request processed
	})
	if err != nil {
		services.Logger.Fatalf("Could not add cron job %s: %v", "processIPDataResponses", err.Error())
	}
	err = c.AddFunc(cfg.Cron.CronScheduleIdentifyTrackingRecords, func() {
		lockAndRunJob(services, identifyTrackingRecordsGroup, identifyTrackingRecords)
	})
	if err != nil {
		services.Logger.Fatalf("Could not add cron job %s: %v", "identifyTrackingRecords", err.Error())
	}
	err = c.AddFunc(cfg.Cron.CronScheduleCreateOrganizationsFromTrackedDataRecords, func() {
		lockAndRunJob(services, createOrganizationsFromTrackedDataGroup, createOrganizationsFromTrackedData)
	})
	if err != nil {
		services.Logger.Fatalf("Could not add cron job %s: %v", "createOrganizationsFromTrackedData", err.Error())
	}
	err = c.AddFunc(cfg.Cron.CronScheduleNotifyOnSlack, func() {
		lockAndRunJob(services, notifyOnSlackGroup, notifyOnSlack)
	})
	if err != nil {
		services.Logger.Fatalf("Could not add cron job %s: %v", "notifyOnSlack", err.Error())
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

func processNewRecords(services *service.Services) {
	services.EnrichDetailsTrackingService.ProcessNewRecords(context.Background())
}

func processIPDataRequests(services *service.Services) {
	services.EnrichDetailsTrackingService.ProcessIPDataRequests(context.Background())
}

func processIPDataResponses(services *service.Services) {
	services.EnrichDetailsTrackingService.ProcessIPDataResponses(context.Background())
}

func identifyTrackingRecords(services *service.Services) {
	services.EnrichDetailsTrackingService.IdentifyTrackingRecords(context.Background())
}

func createOrganizationsFromTrackedData(services *service.Services) {
	services.EnrichDetailsTrackingService.CreateOrganizationsFromTrackedData(context.Background())
}

func notifyOnSlack(services *service.Services) {
	services.EnrichDetailsTrackingService.NotifyOnSlack(context.Background())
}
