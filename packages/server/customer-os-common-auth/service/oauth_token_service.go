package service

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/entity"
)

type OAuthTokenService interface {
	Save(tokenEntity entity.OAuthTokenEntity) (*entity.OAuthTokenEntity, error)
	GetByPlayerIdAndProvider(playerId string, provider string) (*entity.OAuthTokenEntity, error)
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
	qr := o.repositories.OAuthTokenRepository.GetByPlayerIdAndProvider(playerId, provider)
	var oAuthToken entity.OAuthTokenEntity
	var ok bool
	if qr.Error != nil {
		return nil, qr.Error
	} else if qr.Result == nil {
		return nil, nil
	} else {
		oAuthToken, ok = qr.Result.(entity.OAuthTokenEntity)
		if !ok {
			return nil, fmt.Errorf("GetForTenant: unexpected type %T", qr.Result)
		}
	}
	return &oAuthToken, nil
}
