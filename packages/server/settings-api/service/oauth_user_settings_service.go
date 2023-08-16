package service

import (
	"fmt"
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
	qrGoogleProvider := u.repositories.OAuthTokenRepository.GetByPlayerIdAndProvider(playerIdentityId, entity.ProviderGoogle)
	var oAuthToken entity.OAuthTokenEntity

	var ok bool
	if qrGoogleProvider.Error != nil {
		return nil, qrGoogleProvider.Error
	} else if qrGoogleProvider.Result == nil {
		var oAuthSettingsResponse = model.OAuthUserSettingsResponse{
			GoogleCalendarSyncEnabled: false,
			GmailSyncEnabled:          false,
		}
		return &oAuthSettingsResponse, nil
	} else {
		oAuthToken, ok = qrGoogleProvider.Result.(entity.OAuthTokenEntity)
		if !ok {
			return nil, fmt.Errorf("GetForTenant: unexpected type %T", qrGoogleProvider.Result)
		}
	}

	var oAuthSettingsResponse = model.OAuthUserSettingsResponse{
		GoogleCalendarSyncEnabled: oAuthToken.GoogleCalendarSyncEnabled,
		GmailSyncEnabled:          oAuthToken.GmailSyncEnabled,
	}

	return &oAuthSettingsResponse, nil
}
