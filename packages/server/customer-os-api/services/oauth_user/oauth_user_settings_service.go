package api_oauthuser

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	cosapi_interfaces "github.com/customeros/customeros/packages/server/customer-os-api/interfaces"
)

type oAuthUserSettingsService struct {
	log      logger.Logger
	postgres *postgres_repository.Repositories
}

func NewUserSettingsService(log logger.Logger, postgres *postgres_repository.Repositories) cosapi_interfaces.OAuthUserSettingsService {
	return &oAuthUserSettingsService{
		log:      log,
		postgres: postgres,
	}
}

func (u oAuthUserSettingsService) GetTenantOAuthUserSettings(ctx context.Context, tenant string) ([]*cosapi_interfaces.OAuthUserSettingsResponse, error) {
	entities, err := u.postgres.OAuthTokenRepository.GetByTenant(ctx, tenant)
	if err != nil {
		return nil, err
	}

	if entities == nil || len(entities) == 0 {
		return []*cosapi_interfaces.OAuthUserSettingsResponse{}, nil
	}

	oAuthSettingsResponses := make([]*cosapi_interfaces.OAuthUserSettingsResponse, 0)

	for _, entity := range entities {
		oAuthSettingsResponse := cosapi_interfaces.OAuthUserSettingsResponse{
			Provider:           entity.Provider,
			Email:              entity.EmailAddress,
			NeedsManualRefresh: entity.NeedsManualRefresh,
			Type:               entity.Type,
		}
		oAuthSettingsResponses = append(oAuthSettingsResponses, &oAuthSettingsResponse)
	}

	return oAuthSettingsResponses, nil
}
