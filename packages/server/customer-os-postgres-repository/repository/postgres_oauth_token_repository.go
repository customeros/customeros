package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
	"time"
)

type OAuthTokenRepository interface {
	GetAll(ctx context.Context) ([]entity.OAuthTokenEntity, error)
	GetByTenant(ctx context.Context, tenant string) ([]entity.OAuthTokenEntity, error)
	GetByProvider(ctx context.Context, tenant string, provider string) ([]entity.OAuthTokenEntity, error)
	GetByEmailOnly(ctx context.Context, tenant, email string) (*entity.OAuthTokenEntity, error)
	GetByEmail(ctx context.Context, tenant, provider, email string) (*entity.OAuthTokenEntity, error)
	GetByPlayerId(ctx context.Context, tenant, provider, playerId string) (*entity.OAuthTokenEntity, error)
	Save(ctx context.Context, oAuthToken entity.OAuthTokenEntity) (*entity.OAuthTokenEntity, error)
	Update(ctx context.Context, tenant, playerId, provider, accessToken, refreshToken string, expiresAt time.Time) (*entity.OAuthTokenEntity, error)
	MarkForManualRefresh(ctx context.Context, tenant, playerId, provider string) error
	DeleteByEmail(ctx context.Context, tenant, provider, email string) error
}

type oAuthTokenRepository struct {
	db *gorm.DB
}

func NewOAuthTokenRepository(db *gorm.DB) OAuthTokenRepository {
	return &oAuthTokenRepository{
		db: db,
	}
}

func (repo oAuthTokenRepository) GetAll(ctx context.Context) ([]entity.OAuthTokenEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OAuthTokenRepository.GetAll")
	defer span.Finish()

	var entities []entity.OAuthTokenEntity

	err := repo.db.Where("needs_manual_refresh = ?", false).Find(&entities).Error

	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (repo oAuthTokenRepository) GetByTenant(ctx context.Context, tenant string) ([]entity.OAuthTokenEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OAuthTokenRepository.GetByTenant")
	defer span.Finish()
	span.LogFields(log.String("tenant", tenant))

	var entities []entity.OAuthTokenEntity

	err := repo.db.
		Where("tenant_name = ?", tenant).
		Find(&entities).Error

	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (repo oAuthTokenRepository) GetByProvider(ctx context.Context, tenant string, provider string) ([]entity.OAuthTokenEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OAuthTokenRepository.GetAllByProvider")
	defer span.Finish()
	span.LogFields(log.String("tenant", tenant), log.String("provider", provider))

	var entities []entity.OAuthTokenEntity

	err := repo.db.
		Where("tenant_name = ?", tenant).
		Where("provider = ?", provider).
		Find(&entities).Error

	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (repo oAuthTokenRepository) GetByEmailOnly(ctx context.Context, tenant, email string) (*entity.OAuthTokenEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OAuthTokenRepository.GetByEmail")
	defer span.Finish()
	span.LogFields(log.String("tenant", tenant), log.String("email", email))

	var oAuthTokenEntity entity.OAuthTokenEntity

	err := repo.db.
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

func (repo oAuthTokenRepository) GetByEmail(ctx context.Context, tenant, provider, email string) (*entity.OAuthTokenEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OAuthTokenRepository.GetByEmail")
	defer span.Finish()
	span.LogFields(log.String("provider", provider), log.String("tenant", tenant), log.String("email", email))

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

func (repo oAuthTokenRepository) GetByPlayerId(ctx context.Context, tenant, provider, playerId string) (*entity.OAuthTokenEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OAuthTokenRepository.GetByPlayerId")
	defer span.Finish()
	span.LogFields(log.String("tenant", tenant), log.String("playerId", playerId), log.String("provider", provider))

	var oAuthTokenEntity entity.OAuthTokenEntity

	err := repo.db.
		Where("tenant_name = ?", tenant).
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

func (repo oAuthTokenRepository) Save(ctx context.Context, oAuthToken entity.OAuthTokenEntity) (*entity.OAuthTokenEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OAuthTokenRepository.Save")
	defer span.Finish()

	result := repo.db.Save(&oAuthToken)
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return nil, fmt.Errorf("saving oauth token failed: %w", result.Error)
	}
	return &oAuthToken, nil
}

func (repo oAuthTokenRepository) Update(ctx context.Context, tenant, playerId, provider, accessToken, refreshToken string, expiresAt time.Time) (*entity.OAuthTokenEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OAuthTokenRepository.Update")
	defer span.Finish()
	span.LogFields(log.String("tenant", tenant), log.String("playerId", playerId), log.String("provider", provider), log.String("expiresAt", expiresAt.String()))

	existing, err := repo.GetByPlayerId(ctx, tenant, provider, playerId)
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

func (repo oAuthTokenRepository) MarkForManualRefresh(ctx context.Context, tenant, playerId, provider string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OAuthTokenRepository.MarkForManualRefresh")
	defer span.Finish()
	span.LogFields(log.String("playerId", playerId), log.String("provider", provider))

	existing, err := repo.GetByPlayerId(ctx, tenant, provider, playerId)
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

func (repo oAuthTokenRepository) DeleteByEmail(ctx context.Context, tenant, provider, email string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OAuthTokenRepository.DeleteByEmail")
	defer span.Finish()
	span.LogFields(log.String("tenant", tenant), log.String("provider", provider), log.String("email", email))

	existing, err := repo.GetByEmail(ctx, tenant, provider, email)
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
