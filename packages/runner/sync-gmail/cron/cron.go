package cron

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/service"
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

			externalSystemId, err := services.Repositories.ExternalSystemRepository.Merge(ctx, tenant.Name, "gmail")
			if err != nil {
				logrus.Errorf("failed to merge external system: %v", err)
				return
			}

			usersForTenant, err := services.UserService.GetAllUsersForTenant(ctx, tenant.Name)
			if err != nil {
				logrus.Errorf("failed to get users for tenant: %v", err)
				return
			}

			personalEmailProviderList, err := services.Repositories.CommonRepositories.PersonalEmailProviderRepository.GetPersonalEmailProviders()
			if err != nil {
				logrus.Errorf("failed to get personal email provider list: %v", err)
				return
			}

			organizationAllowedForImport, err := services.Repositories.CommonRepositories.WhitelistDomainRepository.GetWhitelistDomains(tenant.Name)
			if err != nil {
				logrus.Errorf("failed to check if organization is allowed for import: %v", err)
				return
			}

			var wgTenant sync.WaitGroup
			wgTenant.Add(len(usersForTenant))

			for _, user := range usersForTenant {
				go func(user entity.UserEntity) {
					defer wgTenant.Done()
					logrus.Infof("syncing emails for user: %s in tenant: %s", user, tenant.Name)

					email, err := services.EmailService.FindEmailForUser(tenant.Name, user.Id)
					if err != nil {
						logrus.Errorf("failed to find email in tenant: %s for user: %s: %v ", tenant.Name, user.Id, err)
						return
					}

					services.EmailService.SyncEmailsForUser(externalSystemId, tenant.Name, email.RawEmail, personalEmailProviderList, organizationAllowedForImport)

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

func syncCalendarEvents(services *service.Services) {
	runId, _ := uuid.NewRandom()
	logrus.Infof("run id: %s syncing calendar events into customer-os at %v", runId.String(), time.Now().UTC())

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

			logrus.Infof("syncing calendar events for tenant: %s", tenant)

			externalSystemId, err := services.Repositories.ExternalSystemRepository.Merge(ctx, tenant.Name, "gcal")
			if err != nil {
				logrus.Errorf("failed to merge external system: %v", err)
				return
			}

			personalEmailProviderList, err := services.Repositories.CommonRepositories.PersonalEmailProviderRepository.GetPersonalEmailProviders()
			if err != nil {
				logrus.Errorf("failed to get personal email provider list: %v", err)
				return
			}

			organizationAllowedForImport, err := services.Repositories.CommonRepositories.WhitelistDomainRepository.GetWhitelistDomains(tenant.Name)
			if err != nil {
				logrus.Errorf("failed to check if organization is allowed for import: %v", err)
				return
			}

			services.MeetingService.SyncCalendarEvents(externalSystemId, tenant.Name, personalEmailProviderList, organizationAllowedForImport)

			logrus.Infof("syncing calendar events for tenant: %s completed", tenant)
		}(*tenant)
	}

	wg.Wait()
	logrus.Infof("syncing calendar events for all tenants completed")
	logrus.Infof("run id: %s sync completed at %v", runId.String(), time.Now().UTC())
}
