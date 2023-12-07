package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository"
	authEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

type googleService struct {
	cfg          *config.Config
	repositories *repository.Repositories
	services     *Services
}

type GoogleService interface {
	ServiceAccountCredentialsExistsForTenant(ctx context.Context, tenant string) (bool, error)

	GetGmailService(ctx context.Context, username, tenant string) (*gmail.Service, error)

	GetGmailServiceWithServiceAccount(ctx context.Context, username string, tenant string) (*gmail.Service, error)
	GetGCalServiceWithServiceAccount(ctx context.Context, username string, tenant string) (*calendar.Service, error)

	GetGmailServiceWithOauthToken(ctx context.Context, tokenEntity authEntity.OAuthTokenEntity) (*gmail.Service, error)
	GetGCalServiceWithOauthToken(ctx context.Context, tokenEntity authEntity.OAuthTokenEntity) (*calendar.Service, error)
}

func NewGoogleService(cfg *config.Config, repositories *repository.Repositories, services *Services) GoogleService {
	return &googleService{
		cfg:          cfg,
		repositories: repositories,
		services:     services,
	}
}

func (s *googleService) ServiceAccountCredentialsExistsForTenant(ctx context.Context, tenant string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GoogleService.ServiceAccountCredentialsExistsForTenant")
	defer span.Finish()
	span.LogFields(log.String("tenant", tenant))

	privateKey, err := s.repositories.ApiKeyRepository.GetApiKeyByTenantService(ctx, tenant, authEntity.GSUITE_SERVICE_PRIVATE_KEY)
	if err != nil {
		return false, nil
	}

	serviceEmail, err := s.repositories.ApiKeyRepository.GetApiKeyByTenantService(ctx, tenant, authEntity.GSUITE_SERVICE_EMAIL_ADDRESS)
	if err != nil {
		return false, nil
	}

	if privateKey == "" || serviceEmail == "" {
		return false, nil
	}

	return true, nil
}

func (s *googleService) GetGmailService(ctx context.Context, username, tenant string) (*gmail.Service, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GoogleService.GetGmailService")
	defer span.Finish()
	span.LogFields(log.String("username", username), log.String("tenant", tenant))

	serviceAccountExistsForTenant, err := s.ServiceAccountCredentialsExistsForTenant(ctx, tenant)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	if serviceAccountExistsForTenant {
		gmailService, err := s.GetGmailServiceWithServiceAccount(ctx, username, tenant)
		if err != nil {
			logrus.Errorf("failed to create gmail service: %v", err)
			return nil, err
		}

		return gmailService, nil
	} else {
		tokenEntity, err := s.repositories.OAuthTokenRepository.GetForEmail(ctx, "google", tenant, username)
		if err != nil {
			return nil, err
		}
		if tokenEntity != nil && tokenEntity.NeedsManualRefresh {
			return nil, nil
		} else if tokenEntity != nil {
			if tokenEntity.RefreshToken == "" {
				err := s.repositories.OAuthTokenRepository.MarkForManualRefresh(ctx, tokenEntity.PlayerIdentityId, tokenEntity.Provider)
				return nil, err
			} else {
				gmailService, err := s.GetGmailServiceWithOauthToken(ctx, *tokenEntity)
				if err != nil {
					logrus.Errorf("failed to create gmail service: %v", err)
					return nil, err
				}
				return gmailService, nil
			}
		} else {
			return nil, nil
		}
	}
}

func (s *googleService) GetGmailServiceWithServiceAccount(ctx context.Context, username, tenant string) (*gmail.Service, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GoogleService.GetGmailServiceWithServiceAccount")
	defer span.Finish()
	span.LogFields(log.String("username", username), log.String("tenant", tenant))

	tok, err := s.getGmailServiceAccountAuthToken(ctx, username, tenant)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve mail token for new gmail service: %v", err)
	}
	client := tok.Client(ctx)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	return srv, err
}

func (s *googleService) getGmailServiceAccountAuthToken(ctx context.Context, identityId, tenant string) (*jwt.Config, error) {
	privateKey, err := s.repositories.ApiKeyRepository.GetApiKeyByTenantService(ctx, tenant, authEntity.GSUITE_SERVICE_PRIVATE_KEY)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve private key for gmail service: %v", err)
	}

	serviceEmail, err := s.repositories.ApiKeyRepository.GetApiKeyByTenantService(ctx, tenant, authEntity.GSUITE_SERVICE_EMAIL_ADDRESS)
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

func (s *googleService) GetGCalServiceWithServiceAccount(ctx context.Context, username, tenant string) (*calendar.Service, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GoogleService.GetGCalServiceWithServiceAccount")
	defer span.Finish()

	tok, err := s.getGCalServiceAccountAuthToken(ctx, username, tenant)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve mail token for new gmail service: %v", err)
	}
	client := tok.Client(ctx)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	return srv, err
}

func (s *googleService) getGCalServiceAccountAuthToken(ctx context.Context, identityId, tenant string) (*jwt.Config, error) {
	privateKey, err := s.repositories.ApiKeyRepository.GetApiKeyByTenantService(ctx, tenant, authEntity.GSUITE_SERVICE_PRIVATE_KEY)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve private key for gmail service: %v", err)
	}

	serviceEmail, err := s.repositories.ApiKeyRepository.GetApiKeyByTenantService(ctx, tenant, authEntity.GSUITE_SERVICE_EMAIL_ADDRESS)
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

func (s *googleService) GetGmailServiceWithOauthToken(ctx context.Context, tokenEntity authEntity.OAuthTokenEntity) (*gmail.Service, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GoogleService.GetGmailServiceWithOauthToken")
	defer span.Finish()

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

	tokenSource := oauth2Config.TokenSource(ctx, &token)
	reuseTokenSource := oauth2.ReuseTokenSource(&token, tokenSource)

	if !token.Valid() {
		newToken, err := reuseTokenSource.Token()
		if err != nil && err.(*oauth2.RetrieveError) != nil && err.(*oauth2.RetrieveError).ErrorCode == "invalid_grant" {
			err := s.repositories.OAuthTokenRepository.MarkForManualRefresh(ctx, tokenEntity.PlayerIdentityId, tokenEntity.Provider)
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

			_, err := s.repositories.OAuthTokenRepository.Update(ctx, tokenEntity.PlayerIdentityId, tokenEntity.Provider, newToken.AccessToken, newToken.RefreshToken, newToken.Expiry)
			if err != nil {
				logrus.Errorf("failed to update token: %v", err)
				return nil, err
			}
		}

	}

	gmailService, err := gmail.NewService(ctx, option.WithTokenSource(reuseTokenSource))
	if err != nil && err.(*oauth2.RetrieveError) != nil && err.(*oauth2.RetrieveError).ErrorCode == "invalid_grant" {
		err := s.repositories.OAuthTokenRepository.MarkForManualRefresh(ctx, tokenEntity.PlayerIdentityId, tokenEntity.Provider)
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
		err3 := s.repositories.OAuthTokenRepository.MarkForManualRefresh(ctx, tokenEntity.PlayerIdentityId, tokenEntity.Provider)
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

func (s *googleService) GetGCalServiceWithOauthToken(ctx context.Context, tokenEntity authEntity.OAuthTokenEntity) (*calendar.Service, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GoogleService.GetGCalServiceWithOauthToken")
	defer span.Finish()

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

	tokenSource := oauth2Config.TokenSource(ctx, &token)
	reuseTokenSource := oauth2.ReuseTokenSource(&token, tokenSource)

	if !token.Valid() {
		newToken, err := reuseTokenSource.Token()
		if err != nil && err.(*oauth2.RetrieveError) != nil && (err.(*oauth2.RetrieveError).ErrorCode == "invalid_grant" || err.(*oauth2.RetrieveError).ErrorCode == "unauthorized_client") {
			err := s.repositories.OAuthTokenRepository.MarkForManualRefresh(ctx, tokenEntity.PlayerIdentityId, tokenEntity.Provider)
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

			_, err := s.repositories.OAuthTokenRepository.Update(ctx, tokenEntity.PlayerIdentityId, tokenEntity.Provider, newToken.AccessToken, newToken.RefreshToken, newToken.Expiry)
			if err != nil {
				logrus.Errorf("failed to update token: %v", err)
				return nil, err
			}
		}

	}

	gCalService, err := calendar.NewService(ctx, option.WithTokenSource(reuseTokenSource))
	if err != nil && err.(*oauth2.RetrieveError) != nil && err.(*oauth2.RetrieveError).ErrorCode == "invalid_grant" {
		err := s.repositories.OAuthTokenRepository.MarkForManualRefresh(ctx, tokenEntity.PlayerIdentityId, tokenEntity.Provider)
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
		err3 := s.repositories.OAuthTokenRepository.MarkForManualRefresh(ctx, tokenEntity.PlayerIdentityId, tokenEntity.Provider)
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
