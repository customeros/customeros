package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository"
)

type OAuthUserSettingsService interface {
	GetOAuthUserSettings(ctx context.Context, tenant, playerIdentityId, provider string) (*model.OAuthUserSettingsResponse, error)
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

func (u oAuthUserSettingsService) GetOAuthUserSettings(ctx context.Context, tenant, playerIdentityId, provider string) (*model.OAuthUserSettingsResponse, error) {
	authProvider, err := u.repositories.AuthRepositories.OAuthTokenRepository.GetByPlayerId(ctx, tenant, provider, playerIdentityId)
	if err != nil {
		return nil, err
	}

	if authProvider == nil {
		return &model.OAuthUserSettingsResponse{}, nil
	}

	var oAuthSettingsResponse = model.OAuthUserSettingsResponse{
		Email:              authProvider.EmailAddress,
		NeedsManualRefresh: authProvider.NeedsManualRefresh,
	}

	return &oAuthSettingsResponse, nil
}

func (u oAuthUserSettingsService) GetTenantOAuthUserSettings(ctx context.Context, tenant string) ([]*model.OAuthUserSettingsResponse, error) {
	entities, err := u.repositories.AuthRepositories.OAuthTokenRepository.GetAllByProvider(ctx, tenant, entity.ProviderGoogle)
	if err != nil {
		return nil, err
	}

	if entities == nil || len(entities) == 0 {
		return []*model.OAuthUserSettingsResponse{}, nil
	}

	var oAuthSettingsResponses = make([]*model.OAuthUserSettingsResponse, 0)

	for _, entity := range entities {
		oAuthSettingsResponse := model.OAuthUserSettingsResponse{
			Email:              entity.EmailAddress,
			UserId:             entity.UserId,
			NeedsManualRefresh: entity.NeedsManualRefresh,
		}
		oAuthSettingsResponses = append(oAuthSettingsResponses, &oAuthSettingsResponse)
	}

	return oAuthSettingsResponses, nil
}
