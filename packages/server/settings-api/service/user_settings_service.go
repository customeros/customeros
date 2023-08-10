package service

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository"
)

type UserSettingsService interface {
	GetOAuthUserSettings(playerIdentityId string, tenant string) (*model.OAuthUserSettingsResponse, error)
}

type userSettingsService struct {
	repositories *repository.PostgresRepositories
	log          logger.Logger
}

func NewUserSettingsService(repositories *repository.PostgresRepositories, log logger.Logger) UserSettingsService {
	return &userSettingsService{
		repositories: repositories,
		log:          log,
	}
}

func (u userSettingsService) GetOAuthUserSettings(playerIdentityId string, tenant string) (*model.OAuthUserSettingsResponse, error) {
	qrGoogleProvider := u.repositories.OAuthTokenRepository.GetByPlayerIdAndTenantAndProvider(playerIdentityId, tenant, entity.ProviderGoogle)
	var oAuthToken entity.OAuthTokenEntity

	var ok bool
	if qrGoogleProvider.Error != nil {
		return nil, qrGoogleProvider.Error
	} else if qrGoogleProvider.Result == nil {
		return nil, nil
	} else {
		oAuthToken, ok = qrGoogleProvider.Result.(entity.OAuthTokenEntity)
		if !ok {
			return nil, fmt.Errorf("GetForTenant: unexpected type %T", qrGoogleProvider.Result)
		}
	}

	var oAuthSettingsResponse = model.OAuthUserSettingsResponse{
		TenantName:             tenant,
		EmailAddress:           oAuthToken.EmailAddress,
		GoogleOAuthSyncEnabled: oAuthToken.EnabledForSync,
	}

	return &oAuthSettingsResponse, nil
}
