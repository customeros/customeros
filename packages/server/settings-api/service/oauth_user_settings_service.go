package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository"
)

type OAuthUserSettingsService interface {
	GetTenantOAuthUserSettings(ctx context.Context, tenant string) ([]*model.OAuthUserSettingsResponse, error)
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

func (u oAuthUserSettingsService) GetTenantOAuthUserSettings(ctx context.Context, tenant string) ([]*model.OAuthUserSettingsResponse, error) {
	entities, err := u.repositories.AuthRepositories.OAuthTokenRepository.GetByTenant(ctx, tenant)
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
