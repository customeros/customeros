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

var jobLock sync.Mutex

func StartCron(config *config.Config, services *service.Services) *cron.Cron {
	c := cron.New()

	err := c.AddFunc(config.SyncData.CronSync, func() {
		lockAndRunJob(services, syncEmails)
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

func syncEmails(services *service.Services) {
	runId, _ := uuid.NewRandom()
	logrus.Infof("run id: %s syncing emails from gmail into customer-os at %v", runId.String(), time.Now().UTC())

	ctx := context.Background()

	tenants, err := services.TenantService.GetAllTenants(ctx)
	if err != nil {
		panic(err) //todo handle error
	}

	var wg sync.WaitGroup
	wg.Add(len(tenants))
	done := make(chan struct{})

	for _, tenant := range tenants {

		go func(tenant entity.TenantEntity) {
			defer wg.Done()

			logrus.Infof("syncing emails for tenant: %s", tenant.Name)

			externalSystemId, err := services.Repositories.ExternalSystemRepository.Merge(ctx, tenant.Name, "gmail")
			if err != nil {
				logrus.Errorf("failed to merge external system: %v", err)
				panic(err) //todo handle error
			}

			users, err := services.UserService.GetAllUsersForTenant(ctx, tenant.Name)
			if err != nil {
				logrus.Errorf("failed to get users for tenant: %v", err)
				panic(err) //todo handle error
			}

			personalEmailProviderList, err := services.Repositories.PersonalEmailProviderRepository.GetPersonalEmailProviderList()
			if err != nil {
				logrus.Errorf("failed to get personal email provider list: %v", err)
				panic(err) //todo handle error
			}

			organizationAllowedForImport, err := services.Repositories.CommonRepositories.ImportAllowedOrganizationRepository.GetOrganizationsAllowedForImport(tenant.Name)
			if err != nil {
				logrus.Errorf("failed to check if organization is allowed for import: %v", err)
				panic(err) //todo handle error
			}

			var wgTenant sync.WaitGroup
			wgTenant.Add(len(users))
			doneTenant := make(chan struct{})

			for _, user := range users {

				go func(user entity.UserEntity) {
					defer wgTenant.Done()

					logrus.Infof("syncing emails in tenant: %s for user: %s", tenant.Name, user.Id)
					time.Sleep(10 * time.Second)

					email, err := services.EmailService.FindEmailForUser(tenant.Name, user.Id)
					if err != nil {
						logrus.Errorf("failed to find email in tenant: %s for user: %s: %v ", tenant.Name, user.Id, err)
						panic(err) //todo handle error
					}

					logrus.Infof("syncing emails in tenant FINISHED: %s for user: %s", tenant.Name, user.Id)

					services.EmailService.SyncEmailsForUser(externalSystemId, "openline", email.RawEmail, personalEmailProviderList, organizationAllowedForImport)

				}(*user)
			}

			go func() {
				wgTenant.Wait()
				close(doneTenant)
			}()

		}(*tenant)
	}
	// Wait for goroutines to finish
	go func() {
		wg.Wait()
		close(done)
	}()
	go func() {
		<-done
	}()

	logrus.Infof("run id: %s sync completed at %v", runId.String(), time.Now().UTC())
}
