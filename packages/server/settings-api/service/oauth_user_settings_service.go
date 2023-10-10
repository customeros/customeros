package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository"
)

type OAuthUserSettingsService interface {
	GetOAuthUserSettings(playerIdentityId string) (*model.OAuthUserSettingsResponse, error)
}

type oAuthUserSettingsService struct {
	repositories *repository.PostgresRepositories
	log          logger.Logger
}

func NewUserSettingsService(repositories *repository.PostgresRepositories, log logger.Logger) OAuthUserSettingsService {
	return &oAuthUserSettingsService{
		repositories: repositories,
		log:          log,
	}
}

func (u oAuthUserSettingsService) GetOAuthUserSettings(playerIdentityId string) (*model.OAuthUserSettingsResponse, error) {
	authProvider, err := u.repositories.AuthRepositories.OAuthTokenRepository.GetByPlayerIdAndProvider(playerIdentityId, entity.ProviderGoogle)
	if err != nil {
		return nil, err
	}

	if authProvider == nil {
		return &model.OAuthUserSettingsResponse{
			GoogleCalendarSyncEnabled: false,
			GmailSyncEnabled:          false,
		}, nil
	}

	var oAuthSettingsResponse = model.OAuthUserSettingsResponse{
		GoogleCalendarSyncEnabled: authProvider.GoogleCalendarSyncEnabled,
		GmailSyncEnabled:          authProvider.GmailSyncEnabled,
	}

	return &oAuthSettingsResponse, nil
}
