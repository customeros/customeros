package repository

import (
	"errors"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/helper"
	"gorm.io/gorm"
)

type OAuthTokenRepository interface {
	GetByPlayerIdAndProvider(playerId string, provider string) helper.QueryResult
	Save(oAuthToken entity.OAuthTokenEntity) (*entity.OAuthTokenEntity, error)
}

type oAuthTokenRepository struct {
	db *gorm.DB
}

func NewOAuthTokenRepository(db *gorm.DB) OAuthTokenRepository {
	return &oAuthTokenRepository{
		db: db,
	}
}

func (oAuthTokenRepo oAuthTokenRepository) GetByPlayerIdAndProvider(playerId string, provider string) helper.QueryResult {
	var oAuthTokenEntity entity.OAuthTokenEntity

	err := oAuthTokenRepo.db.
		Where("player_identity_id = ?", playerId).
		Where("provider = ?", provider).
		First(&oAuthTokenEntity).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.QueryResult{Error: err}
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.QueryResult{Result: nil}
	}

	return helper.QueryResult{Result: oAuthTokenEntity}
}

func (oAuthTokenRepo oAuthTokenRepository) Save(oAuthToken entity.OAuthTokenEntity) (*entity.OAuthTokenEntity, error) {
	result := oAuthTokenRepo.db.Save(&oAuthToken)
	if result.Error != nil {
		return nil, fmt.Errorf("SaveKeys: %w", result.Error)
	}
	return &oAuthToken, nil
}
