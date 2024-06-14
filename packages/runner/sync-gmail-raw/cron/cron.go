package cron

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/service"
	authEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/entity"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var jobLock1 sync.Mutex
var jobLock2 sync.Mutex
var jobLock3 sync.Mutex

func StartCronJobs(config *config.Config, services *service.Services) *cron.Cron {
	c := cron.New()

	err := c.AddFunc(config.SyncData.CronSync, func() {

		go func(jobLock *sync.Mutex) {
			lockAndRunEmailsJob(jobLock, config, services, entity.REAL_TIME, syncEmailsInState)
		}(&jobLock1)

		go func(jobLock *sync.Mutex) {
			lockAndRunEmailsJob(jobLock, config, services, entity.HISTORY, syncEmailsInState)
		}(&jobLock2)

		go func(jobLock *sync.Mutex) {
			lockAndRunJob(jobLock, config, services, syncCalendarEventsForOauthTokens)
		}(&jobLock3)

	})
	if err != nil {
		logrus.Fatalf("Could not add cron job: %v", err.Error())
	}

	c.Start()

	return c
}

func lockAndRunEmailsJob(jobLock *sync.Mutex, config *config.Config, services *service.Services, state entity.GmailImportState, job func(config *config.Config, services *service.Services, state entity.GmailImportState)) {
	jobLock.Lock()
	defer jobLock.Unlock()

	job(config, services, state)
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

func syncEmailsInState(config *config.Config, services *service.Services, state entity.GmailImportState) {
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

			var wgTenant sync.WaitGroup

			privateKey, err := services.CommonServices.PostgresRepositories.GoogleServiceAccountKeyRepository.GetApiKeyByTenantService(ctx, tenant.Name, postgresEntity.GSUITE_SERVICE_PRIVATE_KEY)
			if err != nil {
				logrus.Errorf("failed to get private key: %v", err)
				return
			}
			serviceEmail, err := services.CommonServices.PostgresRepositories.GoogleServiceAccountKeyRepository.GetApiKeyByTenantService(ctx, tenant.Name, postgresEntity.GSUITE_SERVICE_EMAIL_ADDRESS)
			if err != nil {
				logrus.Errorf("failed to get service email: %v", err)
				return
			}

			if privateKey != "" && serviceEmail != "" {
				//import with service account

				usersForTenant, err := services.UserService.GetAllUsersForTenant(ctx, tenant.Name)
				if err != nil {
					logrus.Error(err)
					return
				}

				emailsToSync := make([]string, 0)
				for _, user := range usersForTenant {
					emailsForUser, err := services.EmailService.FindEmailsForUser(tenant.Name, user.Id)
					if err != nil {
						logrus.Infof("failed to find email for user: %v", err)
						return
					}

					if len(emailsForUser) > 0 {
						for _, emailForUser := range emailsForUser {
							if emailForUser.Email == serviceEmail {
								continue
							}
							emailsToSync = append(emailsToSync, emailForUser.Email)
						}
					}
				}

				wgTenant.Add(len(emailsToSync))

				for _, email := range emailsToSync {
					go func(tenant, email string) {
						defer wgTenant.Done()
						syncEmailsForEmailAddress(ctx, config, tenant, email, services, state)
					}(tenant.Name, email)
				}

			} else {
				//import with oauth token

				oAuthTokenEntities, err := services.AuthServices.CommonAuthRepositories.OAuthTokenRepository.GetByTenant(ctx, tenant.Name, "google")
				if err != nil {
					logrus.Errorf("failed to get all oauth tokens: %v", err)
					return
				}

				for _, oAuthTokenEntity := range oAuthTokenEntities {
					if oAuthTokenEntity.NeedsManualRefresh {
						continue
					}

					wgTenant.Add(1)
				}

				for _, oAuthTokenEntity := range oAuthTokenEntities {
					if oAuthTokenEntity.NeedsManualRefresh {
						continue
					}

					go func(oAuthTokenEntity authEntity.OAuthTokenEntity) {
						defer wgTenant.Done()
						syncEmailsForEmailAddress(ctx, config, oAuthTokenEntity.TenantName, oAuthTokenEntity.EmailAddress, services, state)
					}(oAuthTokenEntity)
				}
			}

			wgTenant.Wait()
			logrus.Infof("syncing emails for tenant: %s completed", tenant.Name)

		}(*tenant)
	}

	wg.Wait()
	logrus.Infof("syncing emails for all tenants completed")
	logrus.Infof("run id: %s sync completed at %v", runId.String(), time.Now().UTC())
}

func syncEmailsForEmailAddress(ctx context.Context, config *config.Config, tenant, email string, services *service.Services, state entity.GmailImportState) {
	logrus.Infof("syncing emails for user with email: %s in tenant: %s", email, tenant)

	gmailService, err := services.AuthServices.GoogleService.GetGmailService(ctx, email, tenant)
	if err != nil {
		logrus.Errorf("failed to create gmail service: %v", err)
		return
	}

	if gmailService != nil {
		err := InitializeGmailImportState(services, tenant, email)
		if err != nil {
			logrus.Errorf("failed to initialize gmail import state: %v", err)
			return
		}

		var gmailImportState *entity.UserGmailImportState
		if state == entity.REAL_TIME {
			gmailImportStateLastWeek, err := services.Repositories.UserGmailImportPageTokenRepository.GetGmailImportState(tenant, email, entity.LAST_WEEK)
			if err != nil {
				logrus.Errorf("failed to get gmail import state: %v", err)
				return
			}
			if gmailImportStateLastWeek.Active == true {
				logrus.Infof("gmail import state for tenant: %s and username: %s is active for last week. skipping real time import", tenant, email)
				return
			}

			gmailImportState, err = services.Repositories.UserGmailImportPageTokenRepository.GetGmailImportState(tenant, email, entity.REAL_TIME)
			if err != nil {
				logrus.Errorf("failed to get gmail import state: %v", err)
				return
			}

			if gmailImportState.Active == false {
				importedEmails, err := services.Repositories.RawEmailRepository.CountForUsername("gmail", tenant, email)
				if err != nil {
					logrus.Errorf("failed to count imported emails: %v", err)
					return
				}

				if importedEmails > config.SyncData.BatchSize {
					err = services.Repositories.UserGmailImportPageTokenRepository.ActivateGmailImportState(tenant, email, entity.REAL_TIME)
					if err != nil {
						logrus.Errorf("failed to update gmail import state: %v", err)
						return
					}

					gmailImportState, err = services.Repositories.UserGmailImportPageTokenRepository.GetGmailImportState(tenant, email, entity.REAL_TIME)
					if err != nil {
						logrus.Errorf("failed to get gmail import state: %v", err)
						return
					}
				} else {
					logrus.Infof("gmail import state for tenant: %s and username: %s is not active for real time. skipping real time import", tenant, email)
					return
				}
			}
		} else {
			gmailImportState, err = getHistoryImportGmailImportState(services, tenant, email, entity.LAST_WEEK)
			if err != nil {
				logrus.Errorf("failed to get gmail import state: %v", err)
				return
			}
			if gmailImportState == nil {
				logrus.Infof("no gmail import state found for tenant: %s and username: %s", tenant, email)
				return
			}
		}

		gmailImportState, err = services.EmailService.ReadEmailsForUsername(gmailService, gmailImportState)
		if err != nil {
			logrus.Errorf("failed to read emails for username: %v", err)
			return
		}

		if state == entity.HISTORY && gmailImportState.Cursor == "" {
			err = services.Repositories.UserGmailImportPageTokenRepository.DeactivateGmailImportState(tenant, email, gmailImportState.State)
			if err != nil {
				logrus.Errorf("failed to update gmail import state: %v", err)
				return
			}
		}
	}

	logrus.Infof("syncing emails for user with email: %s in tenant: %s completed", email, tenant)
}

func getHistoryImportGmailImportState(services *service.Services, tenant string, username string, state entity.GmailImportState) (*entity.UserGmailImportState, error) {
	gmailImportState, err := services.Repositories.UserGmailImportPageTokenRepository.GetGmailImportState(tenant, username, state)
	if err != nil {
		return nil, err
	}
	if gmailImportState == nil {
		return nil, fmt.Errorf("failed to get gmail import state for tenant: %s and username: %s and week: %s", tenant, username, state)
	}

	if gmailImportState.Active {
		return gmailImportState, nil
	}

	if state == entity.OLDER_THAN_ONE_YEAR {
		return nil, nil
	}
	importState, err := getNextGmailImportState(state)
	if err != nil {
		return nil, err
	}
	return getHistoryImportGmailImportState(services, tenant, username, importState)
}

func getNextGmailImportState(state entity.GmailImportState) (entity.GmailImportState, error) {
	switch state {
	case entity.LAST_WEEK:
		return entity.LAST_3_MONTHS, nil
	case entity.LAST_3_MONTHS:
		return entity.LAST_YEAR, nil
	case entity.LAST_YEAR:
		return entity.OLDER_THAN_ONE_YEAR, nil
	case entity.OLDER_THAN_ONE_YEAR:
		return entity.OLDER_THAN_ONE_YEAR, fmt.Errorf("invalid state: %s", state)
	default:
		return entity.OLDER_THAN_ONE_YEAR, fmt.Errorf("invalid state: %s", state)
	}
}

func InitializeGmailImportState(services *service.Services, tenantName, userEmail string) error {
	now := time.Now()

	gmailImportState, err := services.Repositories.UserGmailImportPageTokenRepository.GetGmailImportState(tenantName, userEmail, entity.REAL_TIME)
	if err != nil {
		return err
	}
	if gmailImportState == nil {
		_, err = services.Repositories.UserGmailImportPageTokenRepository.CreateGmailImportState(tenantName, userEmail, entity.REAL_TIME, nil, nil, false, "")
		if err != nil {
			return err
		}
	}

	gmailImportState, err = services.Repositories.UserGmailImportPageTokenRepository.GetGmailImportState(tenantName, userEmail, entity.LAST_WEEK)
	if err != nil {
		return err
	}
	if gmailImportState == nil {
		stop := now.AddDate(0, 0, -7)
		_, err = services.Repositories.UserGmailImportPageTokenRepository.CreateGmailImportState(tenantName, userEmail, entity.LAST_WEEK, &now, &stop, true, "")
		if err != nil {
			return err
		}
	}

	gmailImportState, err = services.Repositories.UserGmailImportPageTokenRepository.GetGmailImportState(tenantName, userEmail, entity.LAST_3_MONTHS)
	if err != nil {
		return err
	}
	if gmailImportState == nil {
		stop := now.AddDate(0, -3, 0)
		_, err = services.Repositories.UserGmailImportPageTokenRepository.CreateGmailImportState(tenantName, userEmail, entity.LAST_3_MONTHS, &now, &stop, true, "")
		if err != nil {
			return err
		}
	}

	gmailImportState, err = services.Repositories.UserGmailImportPageTokenRepository.GetGmailImportState(tenantName, userEmail, entity.LAST_YEAR)
	if err != nil {
		return err
	}
	if gmailImportState == nil {
		stop := now.AddDate(-1, 0, 0)
		_, err = services.Repositories.UserGmailImportPageTokenRepository.CreateGmailImportState(tenantName, userEmail, entity.LAST_YEAR, &now, &stop, true, "")
		if err != nil {
			return err
		}
	}

	gmailImportState, err = services.Repositories.UserGmailImportPageTokenRepository.GetGmailImportState(tenantName, userEmail, entity.OLDER_THAN_ONE_YEAR)
	if err != nil {
		return err
	}
	if gmailImportState == nil {
		stop := now.AddDate(-50, 0, 0)
		_, err = services.Repositories.UserGmailImportPageTokenRepository.CreateGmailImportState(tenantName, userEmail, entity.OLDER_THAN_ONE_YEAR, &now, &stop, true, "")
		if err != nil {
			return err
		}
	}

	return nil
}

func syncCalendarEventsForOauthTokens(config *config.Config, services *service.Services) {
	ctx := context.Background()

	runId, _ := uuid.NewRandom()
	logrus.Infof("run id: %s syncing calendar events from google using oauth tokens into customer-os at %v", runId.String(), time.Now().UTC())

	tokenEntities, err := services.AuthServices.CommonAuthRepositories.OAuthTokenRepository.GetAll(ctx)
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

			serviceAccountExistsForTenant, err := services.AuthServices.GoogleService.ServiceAccountCredentialsExistsForTenant(ctx, tokenEntity.TenantName)
			if err != nil {
				logrus.Error(err)
				logrus.Errorf("syncing calendar events for oauth token with email address: %s error", tokenEntity.EmailAddress)
				return
			}

			if serviceAccountExistsForTenant {
				logrus.Infof("service account already exists for personal token for email address: %s. skipping oauth token import", tokenEntity.EmailAddress)
				return
			}

			gCalService, err := services.AuthServices.GoogleService.GetGCalServiceWithOauthToken(ctx, tokenEntity)
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
