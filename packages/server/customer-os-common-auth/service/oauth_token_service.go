package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/entity"
)

type OAuthTokenService interface {
	Save(tokenEntity entity.OAuthTokenEntity) (*entity.OAuthTokenEntity, error)
	GetByPlayerIdAndProvider(playerId string, provider string) (*entity.OAuthTokenEntity, error)
	DeleteByPlayerIdAndProvider(playerId string, provider string) error
}

type oAuthTokenService struct {
	repositories *repository.Repositories
}

func NewOAuthTokenService(repositories *repository.Repositories) OAuthTokenService {
	return &oAuthTokenService{
		repositories: repositories,
	}
}

func (o oAuthTokenService) Save(tokenEntity entity.OAuthTokenEntity) (*entity.OAuthTokenEntity, error) {
	result, err := o.repositories.OAuthTokenRepository.Save(tokenEntity)
	return result, err
}

func (o oAuthTokenService) GetByPlayerIdAndProvider(playerId string, provider string) (*entity.OAuthTokenEntity, error) {
	authTokenEntity, err := o.repositories.OAuthTokenRepository.GetByPlayerIdAndProvider(playerId, provider)

	if err != nil {
		return nil, err
	}

	return authTokenEntity, nil
}

func (o oAuthTokenService) DeleteByPlayerIdAndProvider(playerId string, provider string) error {
	return o.repositories.OAuthTokenRepository.DeleteByPlayerIdAndProvider(playerId, provider)
}
