package cron

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/service"
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
			lockAndRunEmailsJob(jobLock, config, services, postgresEntity.REAL_TIME, syncEmailsInState)
		}(&jobLock1)

		go func(jobLock *sync.Mutex) {
			lockAndRunEmailsJob(jobLock, config, services, postgresEntity.HISTORY, syncEmailsInState)
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

func lockAndRunEmailsJob(jobLock *sync.Mutex, config *config.Config, services *service.Services, state postgresEntity.EmailImportState, job func(config *config.Config, services *service.Services, state postgresEntity.EmailImportState)) {
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

func syncEmailsInState(config *config.Config, services *service.Services, state postgresEntity.EmailImportState) {
	runId, _ := uuid.NewRandom()
	logrus.Infof("run id: %s syncing emails using service account into customer-os at %v", runId.String(), time.Now().UTC())

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
						syncEmailsForEmailAddress(ctx, config, tenant, "google", email, services, state)
					}(tenant.Name, email)
				}

			} else {
				//import with oauth token

				oAuthTokenEntities, err := services.CommonServices.PostgresRepositories.OAuthTokenRepository.GetByTenant(ctx, tenant.Name)
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

					go func(oAuthTokenEntity postgresEntity.OAuthTokenEntity) {
						defer wgTenant.Done()
						syncEmailsForEmailAddress(ctx, config, oAuthTokenEntity.TenantName, oAuthTokenEntity.Provider, oAuthTokenEntity.EmailAddress, services, state)
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

func syncEmailsForEmailAddress(ctx context.Context, config *config.Config, tenant, provider, email string, services *service.Services, state postgresEntity.EmailImportState) {
	logrus.Infof("syncing emails for user with email: %s in tenant: %s", email, tenant)

	var emailImportState *postgresEntity.UserEmailImportState
	if state == postgresEntity.REAL_TIME {

		//activate real time sync only if there are more than Batch Size emails imported, to not overlap with the history sync
		//if there are no emails imported, activate real time sync only if there are no other active syncs

		emailImportStateLastWeek, err := services.CommonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(tenant, provider, email, postgresEntity.LAST_WEEK)
		if err != nil {
			logrus.Errorf("failed to get email import state: %v", err)
			return
		}

		if emailImportStateLastWeek == nil {
			logrus.Infof("no email import state found for tenant: %s and username: %s", tenant, email)
			return
		}

		if emailImportStateLastWeek.Active == true {
			logrus.Infof("gmail import state for tenant: %s and username: %s is active for last week. skipping real time import", tenant, email)
			return
		}

		emailImportState, err = services.CommonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(tenant, provider, email, postgresEntity.REAL_TIME)
		if err != nil {
			logrus.Errorf("failed to get gmail import state: %v", err)
			return
		}

		if emailImportState.Active == false {
			importedEmails, err := services.CommonServices.PostgresRepositories.RawEmailRepository.CountForUsername("gmail", tenant, email)
			if err != nil {
				logrus.Errorf("failed to count imported emails: %v", err)
				return
			}

			if importedEmails > config.SyncData.BatchSize {
				err = services.CommonServices.PostgresRepositories.UserEmailImportPageTokenRepository.ActivateEmailImportState(tenant, provider, email, postgresEntity.REAL_TIME)
				if err != nil {
					logrus.Errorf("failed to update gmail import state: %v", err)
					return
				}

				emailImportState, err = services.CommonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(tenant, provider, email, postgresEntity.REAL_TIME)
				if err != nil {
					logrus.Errorf("failed to get gmail import state: %v", err)
					return
				}
			} else {

				lastWeek, err := services.CommonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(tenant, provider, email, postgresEntity.LAST_WEEK)
				if err != nil {
					logrus.Errorf("failed to get gmail import state: %v", err)
					return
				}

				last3Months, err := services.CommonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(tenant, provider, email, postgresEntity.LAST_3_MONTHS)
				if err != nil {
					logrus.Errorf("failed to get gmail import state: %v", err)
					return
				}

				lastYear, err := services.CommonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(tenant, provider, email, postgresEntity.LAST_YEAR)
				if err != nil {
					logrus.Errorf("failed to get gmail import state: %v", err)
					return
				}

				olderThanOneYear, err := services.CommonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(tenant, provider, email, postgresEntity.OLDER_THAN_ONE_YEAR)
				if err != nil {
					logrus.Errorf("failed to get gmail import state: %v", err)
					return
				}

				if lastWeek.Active == false && last3Months.Active == false && lastYear.Active == false && olderThanOneYear.Active == false {
					err = services.CommonServices.PostgresRepositories.UserEmailImportPageTokenRepository.ActivateEmailImportState(tenant, provider, email, postgresEntity.REAL_TIME)
					if err != nil {
						logrus.Errorf("failed to update gmail import state: %v", err)
						return
					}

					emailImportState, err = services.CommonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(tenant, provider, email, postgresEntity.REAL_TIME)
					if err != nil {
						logrus.Errorf("failed to get gmail import state: %v", err)
						return
					}
				}

				logrus.Infof("gmail import state for tenant: %s and username: %s is not active for real time. skipping real time import", tenant, email)
				return
			}
		}
	} else {

		err := InitializeEmailImportState(services, tenant, provider, email)
		if err != nil {
			logrus.Errorf("failed to initialize email import state: %v", err)
			return
		}

		emailImportState, err = getHistoryImportState(services, tenant, provider, email, postgresEntity.LAST_WEEK)
		if err != nil {
			logrus.Errorf("failed to get email import state: %v", err)
			return
		}
		if emailImportState == nil {
			logrus.Infof("no gmail import state found for tenant: %s and username: %s", tenant, email)
			return
		}
	}

	emailImportState, err := services.EmailService.SyncEmailsForState(ctx, emailImportState)
	if err != nil {
		logrus.Errorf("failed to read emails for username: %v", err)
		return
	}

	if state == postgresEntity.HISTORY && emailImportState.Cursor == "" {
		err = services.CommonServices.PostgresRepositories.UserEmailImportPageTokenRepository.DeactivateEmailImportState(tenant, provider, email, emailImportState.State)
		if err != nil {
			logrus.Errorf("failed to update gmail import state: %v", err)
			return
		}
	}

	logrus.Infof("syncing emails for user with email: %s in tenant: %s completed", email, tenant)
}

func getHistoryImportState(services *service.Services, tenant string, provider, username string, state postgresEntity.EmailImportState) (*postgresEntity.UserEmailImportState, error) {
	emailImportState, err := services.CommonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(tenant, provider, username, state)
	if err != nil {
		return nil, err
	}
	if emailImportState == nil {
		return nil, fmt.Errorf("failed to get gmail import state for tenant: %s and username: %s and week: %s", tenant, username, state)
	}

	if emailImportState.Active {
		return emailImportState, nil
	}

	if state == postgresEntity.OLDER_THAN_ONE_YEAR {
		return nil, nil
	}
	importState, err := getNextEmailImportState(state)
	if err != nil {
		return nil, err
	}
	return getHistoryImportState(services, tenant, provider, username, importState)
}

func getNextEmailImportState(state postgresEntity.EmailImportState) (postgresEntity.EmailImportState, error) {
	switch state {
	case postgresEntity.LAST_WEEK:
		return postgresEntity.LAST_3_MONTHS, nil
	case postgresEntity.LAST_3_MONTHS:
		return postgresEntity.LAST_YEAR, nil
	case postgresEntity.LAST_YEAR:
		return postgresEntity.OLDER_THAN_ONE_YEAR, nil
	case postgresEntity.OLDER_THAN_ONE_YEAR:
		return postgresEntity.OLDER_THAN_ONE_YEAR, fmt.Errorf("invalid state: %s", state)
	default:
		return postgresEntity.OLDER_THAN_ONE_YEAR, fmt.Errorf("invalid state: %s", state)
	}
}

func InitializeEmailImportState(services *service.Services, tenantName, provider, userEmail string) error {
	now := time.Now()

	emailImportState, err := services.CommonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(tenantName, provider, userEmail, postgresEntity.REAL_TIME)
	if err != nil {
		return err
	}
	if emailImportState == nil {
		_, err = services.CommonServices.PostgresRepositories.UserEmailImportPageTokenRepository.CreateEmailImportState(tenantName, provider, userEmail, postgresEntity.REAL_TIME, nil, nil, false, "")
		if err != nil {
			return err
		}
	}

	emailImportState, err = services.CommonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(tenantName, provider, userEmail, postgresEntity.LAST_WEEK)
	if err != nil {
		return err
	}
	if emailImportState == nil {
		stop := now.AddDate(0, 0, -7)
		_, err = services.CommonServices.PostgresRepositories.UserEmailImportPageTokenRepository.CreateEmailImportState(tenantName, provider, userEmail, postgresEntity.LAST_WEEK, &now, &stop, true, "")
		if err != nil {
			return err
		}
	}

	emailImportState, err = services.CommonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(tenantName, provider, userEmail, postgresEntity.LAST_3_MONTHS)
	if err != nil {
		return err
	}
	if emailImportState == nil {
		stop := now.AddDate(0, -3, 0)
		_, err = services.CommonServices.PostgresRepositories.UserEmailImportPageTokenRepository.CreateEmailImportState(tenantName, provider, userEmail, postgresEntity.LAST_3_MONTHS, &now, &stop, true, "")
		if err != nil {
			return err
		}
	}

	emailImportState, err = services.CommonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(tenantName, provider, userEmail, postgresEntity.LAST_YEAR)
	if err != nil {
		return err
	}
	if emailImportState == nil {
		stop := now.AddDate(-1, 0, 0)
		_, err = services.CommonServices.PostgresRepositories.UserEmailImportPageTokenRepository.CreateEmailImportState(tenantName, provider, userEmail, postgresEntity.LAST_YEAR, &now, &stop, true, "")
		if err != nil {
			return err
		}
	}

	emailImportState, err = services.CommonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(tenantName, provider, userEmail, postgresEntity.OLDER_THAN_ONE_YEAR)
	if err != nil {
		return err
	}
	if emailImportState == nil {
		stop := now.AddDate(-50, 0, 0)
		_, err = services.CommonServices.PostgresRepositories.UserEmailImportPageTokenRepository.CreateEmailImportState(tenantName, provider, userEmail, postgresEntity.OLDER_THAN_ONE_YEAR, &now, &stop, true, "")
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

	tokenEntities, err := services.CommonServices.PostgresRepositories.OAuthTokenRepository.GetAll(ctx)
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

		go func(tokenEntity postgresEntity.OAuthTokenEntity) {
			defer wg.Done()

			serviceAccountExistsForTenant, err := services.CommonServices.GoogleService.ServiceAccountCredentialsExistsForTenant(ctx, tokenEntity.TenantName)
			if err != nil {
				logrus.Error(err)
				logrus.Errorf("syncing calendar events for oauth token with email address: %s error", tokenEntity.EmailAddress)
				return
			}

			if serviceAccountExistsForTenant {
				logrus.Infof("service account already exists for personal token for email address: %s. skipping oauth token import", tokenEntity.EmailAddress)
				return
			}

			gCalService, err := services.CommonServices.GoogleService.GetGCalServiceWithOauthToken(ctx, tokenEntity)
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
