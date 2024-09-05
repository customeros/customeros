package cron

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/service"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var jobLock1 sync.Mutex
var jobLock2 sync.Mutex

func StartCron(config *config.Config, services *service.Services) *cron.Cron {
	c := cron.New()

	err := c.AddFunc(config.SyncData.CronSync, func() {

		go func(jobLock *sync.Mutex) {
			lockAndRunJob(jobLock, services, syncEmails)
		}(&jobLock1)

		go func(jobLock *sync.Mutex) {
			lockAndRunJob(jobLock, services, syncCalendarEvents)
		}(&jobLock2)

	})
	if err != nil {
		logrus.Fatalf("Could not add cron job: %v", err.Error())
	}

	c.Start()

	return c
}

func lockAndRunJob(jobLock *sync.Mutex, services *service.Services, job func(services *service.Services)) {
	jobLock.Lock()
	defer jobLock.Unlock()

	job(services)
}

func StopCron(log logger.Logger, cron *cron.Cron) error {
	// Gracefully stop
	log.Info("Gracefully stopping cron")
	cron.Stop()
	return nil
}

func syncEmails(services *service.Services) {
	runId, _ := uuid.NewRandom()
	logrus.Infof("run id: %s syncing emails into customer-os at %v", runId.String(), time.Now().UTC())

	ctx := context.Background()

	distinctUsersForImport, err := services.Repositories.RawEmailRepository.GetDistinctUsersForImport()
	if err != nil {
		logrus.Errorf("failed to get distinct users for import: %v", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(distinctUsersForImport))

	for _, dt := range distinctUsersForImport {

		err = services.Repositories.Neo4jRepositories.ExternalSystemWriteRepository.CreateIfNotExists(ctx, dt.Tenant, neo4jenum.GMail.String(), neo4jenum.GMail.String())
		if err != nil {
			logrus.Errorf("failed to merge external system: %s", err.Error())
			return
		}

		err = services.Repositories.Neo4jRepositories.ExternalSystemWriteRepository.CreateIfNotExists(ctx, dt.Tenant, neo4jenum.Outlook.String(), neo4jenum.Outlook.String())
		if err != nil {
			logrus.Errorf("failed to merge external system: %s", err.Error())
			return
		}

		err = services.Repositories.Neo4jRepositories.ExternalSystemWriteRepository.CreateIfNotExists(ctx, dt.Tenant, neo4jenum.Mailstack.String(), neo4jenum.Mailstack.String())
		if err != nil {
			logrus.Errorf("failed to merge external system: %s", err.Error())
			return
		}

		go func(distinctUser entity.RawEmail) {
			defer wg.Done()

			logrus.Infof("syncing emails for %s in tenant %s", distinctUser.Tenant, distinctUser.Username)

			services.EmailService.SyncEmailsForUser(distinctUser.Tenant, distinctUser.Username)

			logrus.Infof("syncing emails for user: %s in tenant: %s completed", distinctUser.Tenant, distinctUser.Username)
		}(dt)
	}

	wg.Wait()
	logrus.Infof("syncing emails for all tenants completed")
	logrus.Infof("run id: %s sync completed at %v", runId.String(), time.Now().UTC())
}

func syncCalendarEvents(services *service.Services) {
	runId, _ := uuid.NewRandom()
	logrus.Infof("run id: %s syncing calendar events into customer-os at %v", runId.String(), time.Now().UTC())

	ctx := context.Background()

	tenants, err := services.CommonServices.TenantService.GetAllTenants(ctx)
	if err != nil {
		logrus.Errorf("failed to get tenants: %v", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(tenants))

	for _, tenant := range tenants {

		go func(tenant neo4jentity.TenantEntity) {
			defer wg.Done()

			logrus.Infof("syncing calendar events for tenant: %s", tenant.Name)

			services.MeetingService.SyncCalendarEvents("gcal", tenant.Name)

			logrus.Infof("syncing calendar events for tenant: %s completed", tenant.Name)
		}(*tenant)
	}

	wg.Wait()
	logrus.Infof("syncing calendar events for all tenants completed")
	logrus.Infof("run id: %s sync completed at %v", runId.String(), time.Now().UTC())
}
