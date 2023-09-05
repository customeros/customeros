package cron

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/service"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var jobLock sync.Mutex

func StartCron(config *config.Config, services *service.Services) *cron.Cron {
	c := cron.New()

	err := c.AddFunc(config.SyncData.CronSync, func() {
		lockAndRunJob(services, syncEmailsForAllTenantsWithServiceAccount)
	})
	if err != nil {
		logrus.Fatalf("Could not add cron job: %v", err.Error())
	}

	c.Start()

	return c
}

func lockAndRunJob(services *service.Services, job func(services *service.Services)) {
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

func syncEmailsForAllTenantsWithServiceAccount(services *service.Services) {
	runId, _ := uuid.NewRandom()
	logrus.Infof("run id: %s syncing emails from gmail into customer-os at %v", runId.String(), time.Now().UTC())

	ctx := context.Background()

	tenants, err := services.TenantService.GetAllTenants(ctx)
	if err != nil {
		logrus.Errorf("failed to get tenants: %v", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(tenants))

	for _, tenant := range tenants {

		go func(tenant entity.TenantEntity) {
			defer wg.Done()

			logrus.Infof("syncing emails for tenant: %s", tenant)

			serviceAccountExistsForTenant, err := services.EmailService.ServiceAccountCredentialsExistsForTenant(tenant.Name)
			if err != nil {
				logrus.Error(err)
				logrus.Infof("syncing emails for tenant: %s completed", tenant)
				return
			}

			if !serviceAccountExistsForTenant {
				logrus.Infof("no service account credentials found for tenant: %s", tenant.Name)
				logrus.Infof("syncing emails for tenant: %s completed", tenant)
				return
			}

			if tenant.Name != "openline" {
				logrus.Infof("syncing emails for tenant: %s skip", tenant)
				return
			}

			usersForTenant, err := services.UserService.GetAllUsersForTenant(ctx, tenant.Name)
			if err != nil {
				logrus.Error(err)
				return
			}

			var wgTenant sync.WaitGroup
			wgTenant.Add(len(usersForTenant))

			for _, user := range usersForTenant {
				go func(user entity.UserEntity) {
					defer wgTenant.Done()
					logrus.Infof("syncing emails for user: %s in tenant: %s", user, tenant)

					emailForUser, err := services.EmailService.FindEmailForUser(tenant.Name, user.Id)
					if err != nil {
						logrus.Infof("failed to find email for user: %v", err)
						return
					}
					services.EmailService.ReadNewEmailsForUsername(tenant.Name, emailForUser.RawEmail)

					logrus.Infof("syncing emails for user: %s in tenant: %s completed", user, tenant)
				}(*user)
			}

			wgTenant.Wait()
			logrus.Infof("syncing emails for tenant: %s completed", tenant)

		}(*tenant)
	}

	wg.Wait()
	logrus.Infof("syncing emails for all tenants completed")
	logrus.Infof("run id: %s sync completed at %v", runId.String(), time.Now().UTC())
}
