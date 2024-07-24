package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/model"
)

type OAuthUserSettingsService interface {
	GetTenantOAuthUserSettings(ctx context.Context, tenant string) ([]*model.OAuthUserSettingsResponse, error)
}

type oAuthUserSettingsService struct {
	services *Services
	log      logger.Logger
}

func NewUserSettingsService(services *Services, log logger.Logger) OAuthUserSettingsService {
	return &oAuthUserSettingsService{
		services: services,
		log:      log,
	}
}

func (u oAuthUserSettingsService) GetTenantOAuthUserSettings(ctx context.Context, tenant string) ([]*model.OAuthUserSettingsResponse, error) {
	entities, err := u.services.CommonServices.PostgresRepositories.OAuthTokenRepository.GetByTenant(ctx, tenant)
	if err != nil {
		return nil, err
	}

	if entities == nil || len(entities) == 0 {
		return []*model.OAuthUserSettingsResponse{}, nil
	}

	var oAuthSettingsResponses = make([]*model.OAuthUserSettingsResponse, 0)

	for _, entity := range entities {
		oAuthSettingsResponse := model.OAuthUserSettingsResponse{
			Provider:           entity.Provider,
			Email:              entity.EmailAddress,
			NeedsManualRefresh: entity.NeedsManualRefresh,
			Type:               entity.Type,
		}
		oAuthSettingsResponses = append(oAuthSettingsResponses, &oAuthSettingsResponse)
	}

	return oAuthSettingsResponses, nil
}
