package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/repository"
	authEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/entity"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

type gmailService struct {
	cfg          *config.Config
	repositories *repository.Repositories
	services     *Services
}

type GmailService interface {
	ServiceAccountCredentialsExistsForTenant(tenant string) (bool, error)

	GetGmailServiceWithServiceAccount(username string, tenant string) (*gmail.Service, error)
	GetGCalServiceWithServiceAccount(username string, tenant string) (*calendar.Service, error)

	GetGmailServiceWithOauthToken(tokenEntity authEntity.OAuthTokenEntity) (*gmail.Service, error)
	GetGCalServiceWithOauthToken(tokenEntity authEntity.OAuthTokenEntity) (*calendar.Service, error)
}

func (s *gmailService) ServiceAccountCredentialsExistsForTenant(tenant string) (bool, error) {
	privateKey, err := s.repositories.ApiKeyRepository.GetApiKeyByTenantService(tenant, repository.GSUITE_SERVICE_PRIVATE_KEY)
	if err != nil {
		return false, nil
	}

	serviceEmail, err := s.repositories.ApiKeyRepository.GetApiKeyByTenantService(tenant, repository.GSUITE_SERVICE_EMAIL_ADDRESS)
	if err != nil {
		return false, nil
	}

	if privateKey == "" || serviceEmail == "" {
		return false, nil
	}

	return true, nil
}

func (s *gmailService) GetGmailServiceWithServiceAccount(username string, tenant string) (*gmail.Service, error) {
	tok, err := s.getGmailServiceAccountAuthToken(username, tenant)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve mail token for new gmail service: %v", err)
	}
	ctx := context.Background()
	client := tok.Client(ctx)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	return srv, err
}

func (s *gmailService) getGmailServiceAccountAuthToken(identityId, tenant string) (*jwt.Config, error) {
	privateKey, err := s.repositories.ApiKeyRepository.GetApiKeyByTenantService(tenant, repository.GSUITE_SERVICE_PRIVATE_KEY)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve private key for gmail service: %v", err)
	}

	serviceEmail, err := s.repositories.ApiKeyRepository.GetApiKeyByTenantService(tenant, repository.GSUITE_SERVICE_EMAIL_ADDRESS)
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

func (s *gmailService) GetGCalServiceWithServiceAccount(username string, tenant string) (*calendar.Service, error) {
	tok, err := s.getGCalServiceAccountAuthToken(username, tenant)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve mail token for new gmail service: %v", err)
	}
	ctx := context.Background()
	client := tok.Client(ctx)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	return srv, err
}

func (s *gmailService) getGCalServiceAccountAuthToken(identityId, tenant string) (*jwt.Config, error) {
	privateKey, err := s.repositories.ApiKeyRepository.GetApiKeyByTenantService(tenant, repository.GSUITE_SERVICE_PRIVATE_KEY)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve private key for gmail service: %v", err)
	}

	serviceEmail, err := s.repositories.ApiKeyRepository.GetApiKeyByTenantService(tenant, repository.GSUITE_SERVICE_EMAIL_ADDRESS)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve service email for gmail service: %v", err)
	}
	conf := &jwt.Config{
		Email:      serviceEmail,
		PrivateKey: []byte(privateKey),
		TokenURL:   google.JWTTokenURL,
		Scopes:     []string{"https://calendar.google.com/"},
		Subject:    identityId,
	}
	return conf, nil
}

func (s *gmailService) GetGmailServiceWithOauthToken(tokenEntity authEntity.OAuthTokenEntity) (*gmail.Service, error) {
	oauth2Config := &oauth2.Config{
		ClientID:     s.cfg.GoogleOAuth.ClientId,
		ClientSecret: s.cfg.GoogleOAuth.ClientSecret,
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
		if err != nil && err.(*oauth2.RetrieveError) != nil && err.(*oauth2.RetrieveError).ErrorCode == "invalid_grant" {
			err := s.repositories.OAuthRepositories.OAuthTokenRepository.MarkForManualRefresh(tokenEntity.PlayerIdentityId, tokenEntity.Provider)
			if err != nil {
				logrus.Errorf("failed to mark token for manual refresh: %v", err)
				return nil, err
			}
			return nil, fmt.Errorf("token is invalid and marked for manual refresh")
		} else if err != nil {
			logrus.Errorf("failed to get new token: %v", err)
			return nil, err
		}

		if newToken.AccessToken != tokenEntity.AccessToken {

			_, err := s.repositories.OAuthRepositories.OAuthTokenRepository.Update(tokenEntity.PlayerIdentityId, tokenEntity.Provider, newToken.AccessToken, newToken.RefreshToken, newToken.Expiry)
			if err != nil {
				logrus.Errorf("failed to update token: %v", err)
				return nil, err
			}
		}

	}

	gmailService, err := gmail.NewService(context.TODO(), option.WithTokenSource(reuseTokenSource))
	if err != nil && err.(*oauth2.RetrieveError) != nil && err.(*oauth2.RetrieveError).ErrorCode == "invalid_grant" {
		err := s.repositories.OAuthRepositories.OAuthTokenRepository.MarkForManualRefresh(tokenEntity.PlayerIdentityId, tokenEntity.Provider)
		if err != nil {
			logrus.Errorf("failed to mark token for manual refresh: %v", err)
			return nil, err
		}
		return nil, fmt.Errorf("token is invalid and marked for manual refresh")
	} else if err != nil {
		logrus.Errorf("failed to create gmail service for token: %v", err)
		return nil, err
	}

	//Request had invalid authentication credentials. Expected OAuth 2 access token, login cookie or other valid authentication credential.
	//See https://developers.google.com/identity/sign-in/web/devconsole-project.
	_, err2 := gmailService.Users.GetProfile("me").Do()
	if err2 != nil && err2.(*googleapi.Error) != nil && err2.(*googleapi.Error).Code == 401 {
		err3 := s.repositories.OAuthRepositories.OAuthTokenRepository.MarkForManualRefresh(tokenEntity.PlayerIdentityId, tokenEntity.Provider)
		if err3 != nil {
			logrus.Errorf("failed to mark token for manual refresh: %v", err)
			return nil, err3
		}
		return nil, fmt.Errorf("token is invalid and marked for manual refresh")
	} else if err2 != nil {
		logrus.Errorf("failed to get new token: %v", err)
		return nil, err2
	}

	return gmailService, nil
}

func (s *gmailService) GetGCalServiceWithOauthToken(tokenEntity authEntity.OAuthTokenEntity) (*calendar.Service, error) {
	oauth2Config := &oauth2.Config{
		ClientID:     s.cfg.GoogleOAuth.ClientId,
		ClientSecret: s.cfg.GoogleOAuth.ClientSecret,
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
		if err != nil && err.(*oauth2.RetrieveError) != nil && (err.(*oauth2.RetrieveError).ErrorCode == "invalid_grant" || err.(*oauth2.RetrieveError).ErrorCode == "unauthorized_client") {
			err := s.repositories.OAuthRepositories.OAuthTokenRepository.MarkForManualRefresh(tokenEntity.PlayerIdentityId, tokenEntity.Provider)
			if err != nil {
				logrus.Errorf("failed to mark token for manual refresh: %v", err)
				return nil, err
			}
			return nil, fmt.Errorf("token is invalid and marked for manual refresh")
		} else if err != nil {
			logrus.Errorf("failed to get new token: %v", err)
			return nil, err
		}

		if newToken.AccessToken != tokenEntity.AccessToken {

			_, err := s.repositories.OAuthRepositories.OAuthTokenRepository.Update(tokenEntity.PlayerIdentityId, tokenEntity.Provider, newToken.AccessToken, newToken.RefreshToken, newToken.Expiry)
			if err != nil {
				logrus.Errorf("failed to update token: %v", err)
				return nil, err
			}
		}

	}

	gCalService, err := calendar.NewService(context.TODO(), option.WithTokenSource(reuseTokenSource))
	if err != nil && err.(*oauth2.RetrieveError) != nil && err.(*oauth2.RetrieveError).ErrorCode == "invalid_grant" {
		err := s.repositories.OAuthRepositories.OAuthTokenRepository.MarkForManualRefresh(tokenEntity.PlayerIdentityId, tokenEntity.Provider)
		if err != nil {
			logrus.Errorf("failed to mark token for manual refresh: %v", err)
			return nil, err
		}
		return nil, fmt.Errorf("token is invalid and marked for manual refresh")
	} else if err != nil {
		logrus.Errorf("failed to create gmail service for token: %v", err)
		return nil, err
	}

	//Request had invalid authentication credentials. Expected OAuth 2 access token, login cookie or other valid authentication credential.
	//See https://developers.google.com/identity/sign-in/web/devconsole-project.
	_, err2 := gCalService.CalendarList.Get("primary").Do()
	if err2 != nil && err2.(*googleapi.Error) != nil && err2.(*googleapi.Error).Code == 401 {
		err3 := s.repositories.OAuthRepositories.OAuthTokenRepository.MarkForManualRefresh(tokenEntity.PlayerIdentityId, tokenEntity.Provider)
		if err3 != nil {
			logrus.Errorf("failed to mark token for manual refresh: %v", err)
			return nil, err3
		}
		return nil, fmt.Errorf("token is invalid and marked for manual refresh")
	} else if err2 != nil {
		logrus.Errorf("failed to get new token: %v", err)
		return nil, err2
	}

	return gCalService, nil
}

func NewGmailService(cfg *config.Config, repositories *repository.Repositories, services *Services) GmailService {
	return &gmailService{
		cfg:          cfg,
		repositories: repositories,
		services:     services,
	}
}
