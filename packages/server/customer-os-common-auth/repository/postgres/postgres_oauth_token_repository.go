package repository

import (
	"errors"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/entity"
	"gorm.io/gorm"
	"time"
)

type OAuthTokenRepository interface {
	GetAll() ([]entity.OAuthTokenEntity, error)
	GetByPlayerIdAndProvider(playerId string, provider string) (*entity.OAuthTokenEntity, error)
	GetForEmail(provider, tenant, email string) (*entity.OAuthTokenEntity, error)

	Save(oAuthToken entity.OAuthTokenEntity) (*entity.OAuthTokenEntity, error)
	Update(playerId, provider, accessToken, refreshToken string, expiresAt time.Time) (*entity.OAuthTokenEntity, error)

	MarkForManualRefresh(playerId, provider string) error

	DeleteByPlayerIdAndProvider(playerId string, provider string) error
}

type oAuthTokenRepository struct {
	db *gorm.DB
}

func NewOAuthTokenRepository(db *gorm.DB) OAuthTokenRepository {
	return &oAuthTokenRepository{
		db: db,
	}
}

func (repo oAuthTokenRepository) GetAll() ([]entity.OAuthTokenEntity, error) {
	var entities []entity.OAuthTokenEntity

	err := repo.db.Where("needs_manual_refresh = ?", false).Find(&entities).Error

	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (repo oAuthTokenRepository) GetByPlayerIdAndProvider(playerId, provider string) (*entity.OAuthTokenEntity, error) {
	var oAuthTokenEntity entity.OAuthTokenEntity

	err := repo.db.
		Where("player_identity_id = ?", playerId).
		Where("provider = ?", provider).
		First(&oAuthTokenEntity).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &oAuthTokenEntity, nil
}

func (repo oAuthTokenRepository) GetForEmail(provider, tenant, email string) (*entity.OAuthTokenEntity, error) {
	var oAuthTokenEntity entity.OAuthTokenEntity

	err := repo.db.
		Where("provider = ?", provider).
		Where("tenant_name = ?", tenant).
		Where("email_address = ?", email).
		First(&oAuthTokenEntity).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &oAuthTokenEntity, nil
}

func (repo oAuthTokenRepository) Save(oAuthToken entity.OAuthTokenEntity) (*entity.OAuthTokenEntity, error) {
	result := repo.db.Save(&oAuthToken)
	if result.Error != nil {
		return nil, fmt.Errorf("saving oauth token failed: %w", result.Error)
	}
	return &oAuthToken, nil
}

func (repo oAuthTokenRepository) Update(playerId, provider, accessToken, refreshToken string, expiresAt time.Time) (*entity.OAuthTokenEntity, error) {
	existing, err := repo.GetByPlayerIdAndProvider(playerId, provider)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("oauth token not found")
	}

	existing.AccessToken = accessToken
	existing.RefreshToken = refreshToken
	existing.ExpiresAt = expiresAt

	result := repo.db.Save(&existing)
	if result.Error != nil {
		return nil, fmt.Errorf("updating oauth token failed: %w", result.Error)
	}

	return existing, nil
}

func (repo oAuthTokenRepository) MarkForManualRefresh(playerId, provider string) error {
	existing, err := repo.GetByPlayerIdAndProvider(playerId, provider)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("oauth token not found")
	}

	existing.NeedsManualRefresh = true

	result := repo.db.Save(&existing)
	if result.Error != nil {
		return fmt.Errorf("updating oauth token failed: %w", result.Error)
	}

	return nil
}

func (repo oAuthTokenRepository) DeleteByPlayerIdAndProvider(playerId, provider string) error {
	existing, err := repo.GetByPlayerIdAndProvider(playerId, provider)
	if err != nil {
		return err
	}
	if existing == nil {
		return nil
	}

	err = repo.db.Delete(&existing).Error
	if err != nil {
		return fmt.Errorf("deleting oauth token failed: %w", err)
	}

	return nil
}
