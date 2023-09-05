package cron

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/service"
	authEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/entity"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"sync"
	"time"
)

var jobLock1 sync.Mutex
var jobLock2 sync.Mutex

func StartCronJobs(config *config.Config, services *service.Services) *cron.Cron {
	c := cron.New()

	err := c.AddFunc(config.SyncData.CronSync, func() {

		go func(jobLock *sync.Mutex) {
			lockAndRunJob(jobLock, config, services, syncEmailsForAllTenantsWithServiceAccount)
		}(&jobLock1)

		go func(jobLock *sync.Mutex) {
			lockAndRunJob(jobLock, config, services, syncEmailsForOauthTokens)
		}(&jobLock2)

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

func syncEmailsForAllTenantsWithServiceAccount(config *config.Config, services *service.Services) {
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

					gmailService, err := getGmailServiceWithServiceAccount(services, emailForUser.RawEmail, tenant.Name)
					if err != nil {
						logrus.Errorf("failed to create gmail service: %v", err)
						return
					}

					services.EmailService.ReadNewEmailsForUsername(gmailService, tenant.Name, emailForUser.RawEmail)

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

func syncEmailsForOauthTokens(config *config.Config, services *service.Services) {
	runId, _ := uuid.NewRandom()
	logrus.Infof("run id: %s syncing emails from gmail using oauth tokens into customer-os at %v", runId.String(), time.Now().UTC())

	tokenEntities, err := services.Repositories.OAuthRepositories.OAuthTokenRepository.GetAll()
	if err != nil {
		logrus.Errorf("failed to get all oauth tokens: %v", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(tokenEntities))

	for _, tokenEntity := range tokenEntities {

		go func(tokenEntity authEntity.OAuthTokenEntity) {
			defer wg.Done()

			serviceAccountExistsForTenant, err := services.EmailService.ServiceAccountCredentialsExistsForTenant(tokenEntity.TenantName)
			if err != nil {
				logrus.Error(err)
				logrus.Infof("syncing emails for tenant: %s completed", tokenEntity.TenantName)
				return
			}

			if serviceAccountExistsForTenant {
				logrus.Infof("service account already exists for tenant: %s. skipping personal access import", tokenEntity.TenantName)
				return
			}

			gmailService, err := getGmailServiceWithOauthToken(config, services, tokenEntity)
			if err != nil {
				logrus.Errorf("failed to create gmail service: %v", err)
				return
			}

			services.EmailService.ReadNewEmailsForUsername(gmailService, tokenEntity.TenantName, tokenEntity.EmailAddress)
		}(tokenEntity)
	}

	wg.Wait()
	logrus.Infof("syncing emails for all oauth tokens completed")
	logrus.Infof("run id: %s sync completed at %v", runId.String(), time.Now().UTC())
}

func getGmailServiceWithServiceAccount(services *service.Services, username string, tenant string) (*gmail.Service, error) {
	tok, err := getServiceAccountAuthToken(services, username, tenant)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve mail token for new gmail service: %v", err)
	}
	ctx := context.Background()
	client := tok.Client(ctx)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	return srv, err
}

func getServiceAccountAuthToken(services *service.Services, identityId, tenant string) (*jwt.Config, error) {
	privateKey, err := services.Repositories.ApiKeyRepository.GetApiKeyByTenantService(tenant, repository.GSUITE_SERVICE_PRIVATE_KEY)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve private key for gmail service: %v", err)
	}

	serviceEmail, err := services.Repositories.ApiKeyRepository.GetApiKeyByTenantService(tenant, repository.GSUITE_SERVICE_EMAIL_ADDRESS)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve service email for gmail service: %v", err)
	}
	conf := &jwt.Config{
		Email:      serviceEmail,
		PrivateKey: []byte(privateKey),
		TokenURL:   google.JWTTokenURL,
		Scopes:     []string{"https://mail.google.com/"},
		Subject:    identityId,
	}
	return conf, nil
}

func getGmailServiceWithOauthToken(config *config.Config, services *service.Services, tokenEntity authEntity.OAuthTokenEntity) (*gmail.Service, error) {
	oauth2Config := &oauth2.Config{
		ClientID:     config.GoogleOAuth.ClientId,
		ClientSecret: config.GoogleOAuth.ClientSecret,
		Endpoint:     google.Endpoint,
	}

	token := oauth2.Token{
		AccessToken:  tokenEntity.AccessToken,
		RefreshToken: tokenEntity.RefreshToken,
		Expiry:       tokenEntity.ExpiresAt,
	}

	tokenSource := oauth2Config.TokenSource(context.TODO(), &token)
	reuseTokenSource := oauth2.ReuseTokenSource(&token, tokenSource)

	if !token.Valid() {
		newToken, err := reuseTokenSource.Token()
		if err != nil {
			logrus.Errorf("failed to get new token: %v", err)
			return nil, err
		}

		if newToken.AccessToken != tokenEntity.AccessToken {

			_, err := services.Repositories.OAuthRepositories.OAuthTokenRepository.Update(tokenEntity.PlayerIdentityId, tokenEntity.Provider, newToken.AccessToken, newToken.RefreshToken, newToken.Expiry)
			if err != nil {
				logrus.Errorf("failed to update token: %v", err)
				return nil, err
			}
		}

	}

	gmailService, err := gmail.NewService(context.TODO(), option.WithTokenSource(reuseTokenSource))
	if err != nil {
		logrus.Errorf("failed to create gmail service for token: %v", err)
		return nil, err
	}

	return gmailService, nil
}
