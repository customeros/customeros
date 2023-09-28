package cron

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/service"
	authEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/entity"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var jobLock1 sync.Mutex
var jobLock2 sync.Mutex
var jobLock3 sync.Mutex
var jobLock4 sync.Mutex

func StartCronJobs(config *config.Config, services *service.Services) *cron.Cron {
	c := cron.New()

	err := c.AddFunc(config.SyncData.CronSync, func() {

		go func(jobLock *sync.Mutex) {
			lockAndRunJob(jobLock, config, services, syncEmails)
		}(&jobLock1)

		//go func(jobLock *sync.Mutex) {
		//	lockAndRunJob(jobLock, config, services, syncCalendarEventsForAllTenantsWithServiceAccount)
		//}(&jobLock3)

		//go func(jobLock *sync.Mutex) {
		//	lockAndRunJob(jobLock, config, services, syncCalendarEventsForOauthTokens)
		//}(&jobLock4)

	})
	if err != nil {
		logrus.Fatalf("Could not add cron job: %v", err.Error())
	}

	c.Start()

	return c
}

func lockAndRunJob(jobLock *sync.Mutex, config *config.Config, services *service.Services, job func(config *config.Config, services *service.Services)) {
	jobLock.Lock()
	defer jobLock.Unlock()

	job(config, services)
}

func StopCron(log logger.Logger, cron *cron.Cron) error {
	// Gracefully stop
	log.Info("Gracefully stopping cron")
	cron.Stop()
	return nil
}

func syncEmails(config *config.Config, services *service.Services) {
	runId, _ := uuid.NewRandom()
	logrus.Infof("run id: %s syncing emails from gmail using service account into customer-os at %v", runId.String(), time.Now().UTC())

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

			logrus.Infof("syncing emails for tenant: %s", tenant.Name)

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

					emailForUser, err := services.EmailService.FindEmailForUser(tenant.Name, user.Id)
					if err != nil {
						logrus.Infof("failed to find email for user: %v", err)
						return
					}

					logrus.Infof("syncing emails for user with email: %s in tenant: %s", emailForUser.RawEmail, tenant.Name)

					gmailService, err := services.AuthServices.GoogleService.GetGmailService(emailForUser.RawEmail, tenant.Name)
					if err != nil {
						logrus.Errorf("failed to create gmail service: %v", err)
						return
					}

					if gmailService != nil {
						services.EmailService.ReadNewEmailsForUsername(gmailService, tenant.Name, emailForUser.RawEmail)
					}

					logrus.Infof("syncing emails for user with email: %s in tenant: %s completed", emailForUser.RawEmail, tenant.Name)
				}(*user)
			}

			wgTenant.Wait()
			logrus.Infof("syncing emails for tenant: %s completed", tenant.Name)

		}(*tenant)
	}

	wg.Wait()
	logrus.Infof("syncing emails for all tenants completed")
	logrus.Infof("run id: %s sync completed at %v", runId.String(), time.Now().UTC())
}

func syncCalendarEventsForAllTenantsWithServiceAccount(config *config.Config, services *service.Services) {
	runId, _ := uuid.NewRandom()
	logrus.Infof("run id: %s syncing emails from gmail using service account into customer-os at %v", runId.String(), time.Now().UTC())

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

			logrus.Infof("syncing calendar events for tenant: %s", tenant.Name)

			serviceAccountExistsForTenant, err := services.AuthServices.GoogleService.ServiceAccountCredentialsExistsForTenant(tenant.Name)
			if err != nil {
				logrus.Error(err)
				logrus.Infof("syncing calendar events for tenant: %s completed", tenant.Name)
				return
			}

			if !serviceAccountExistsForTenant {
				logrus.Infof("no service account credentials found for tenant: %s", tenant.Name)
				logrus.Infof("syncing calendar events for tenant: %s completed", tenant.Name)
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

					emailForUser, err := services.EmailService.FindEmailForUser(tenant.Name, user.Id)
					if err != nil {
						logrus.Infof("failed to find email for user: %v", err)
						return
					}

					logrus.Infof("syncing calendar events for user with email: %s in tenant: %s", emailForUser.RawEmail, tenant.Name)

					gCalService, err := services.AuthServices.GoogleService.GetGCalServiceWithServiceAccount(emailForUser.RawEmail, tenant.Name)
					if err != nil {
						logrus.Errorf("failed to create gmail service: %v", err)
						return
					}

					services.MeetingService.ReadNewCalendarEventsForUsername(gCalService, tenant.Name, emailForUser.RawEmail)

					logrus.Infof("syncing calendar events for user with email: %s in tenant: %s completed", emailForUser.RawEmail, tenant.Name)
				}(*user)
			}

			wgTenant.Wait()
			logrus.Infof("syncing calendar events for tenant: %s completed", tenant.Name)

		}(*tenant)
	}

	wg.Wait()
	logrus.Infof("syncing emails for all tenants completed")
	logrus.Infof("run id: %s sync completed at %v", runId.String(), time.Now().UTC())
}

func syncCalendarEventsForOauthTokens(config *config.Config, services *service.Services) {
	runId, _ := uuid.NewRandom()
	logrus.Infof("run id: %s syncing calendar events from google using oauth tokens into customer-os at %v", runId.String(), time.Now().UTC())

	tokenEntities, err := services.AuthServices.CommonAuthRepositories.OAuthTokenRepository.GetAll()
	if err != nil {
		logrus.Errorf("failed to get all oauth tokens: %v", err)
		return
	}

	var wg sync.WaitGroup

	for _, tokenEntity := range tokenEntities {

		if !tokenEntity.GoogleCalendarSyncEnabled {
			continue
		}

		wg.Add(1)

		go func(tokenEntity authEntity.OAuthTokenEntity) {
			defer wg.Done()

			serviceAccountExistsForTenant, err := services.AuthServices.GoogleService.ServiceAccountCredentialsExistsForTenant(tokenEntity.TenantName)
			if err != nil {
				logrus.Error(err)
				logrus.Errorf("syncing calendar events for oauth token with email address: %s error", tokenEntity.EmailAddress)
				return
			}

			if serviceAccountExistsForTenant {
				logrus.Infof("service account already exists for personal token for email address: %s. skipping oauth token import", tokenEntity.EmailAddress)
				return
			}

			gCalService, err := services.AuthServices.GoogleService.GetGCalServiceWithOauthToken(tokenEntity)
			if err != nil {
				logrus.Errorf("failed to create gmail service: %v", err)
				return
			}

			services.MeetingService.ReadNewCalendarEventsForUsername(gCalService, tokenEntity.TenantName, tokenEntity.EmailAddress)

			logrus.Infof("syncing calendar events with personal token for email address: %s completed", tokenEntity.EmailAddress)

		}(tokenEntity)
	}

	wg.Wait()
	logrus.Infof("syncing calendar events for all oauth tokens completed")
	logrus.Infof("run id: %s sync completed at %v", runId.String(), time.Now().UTC())
}
