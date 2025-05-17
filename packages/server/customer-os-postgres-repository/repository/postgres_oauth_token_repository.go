package postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"gorm.io/gorm"
)

type OAuthTokenRepository interface {
	GetAll(ctx context.Context) ([]postgres_entity.OAuthTokenEntity, error)
	GetByTenant(ctx context.Context, tenant string) ([]postgres_entity.OAuthTokenEntity, error)
	GetByEmailAndProvider(ctx context.Context, tenant, provider, email string) (*postgres_entity.OAuthTokenEntity, error)
	GetByEmail(ctx context.Context, tenant, email string) (*postgres_entity.OAuthTokenEntity, error)
	GetByPlayerId(ctx context.Context, tenant, provider, playerId string) (*postgres_entity.OAuthTokenEntity, error)
	Save(ctx context.Context, oAuthToken postgres_entity.OAuthTokenEntity) (*postgres_entity.OAuthTokenEntity, error)
	Update(ctx context.Context, tenant, playerId, provider, accessToken, refreshToken string, expiresAt time.Time) (*postgres_entity.OAuthTokenEntity, error)
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

func (repo oAuthTokenRepository) GetAll(ctx context.Context) ([]postgres_entity.OAuthTokenEntity, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "OAuthTokenRepository.GetAll")
	defer spans.Finish()

	var entities []postgres_entity.OAuthTokenEntity

	err := repo.db.Where("needs_manual_refresh = ?", false).Find(&entities).Error

	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return entities, nil
}

func (repo oAuthTokenRepository) GetByTenant(ctx context.Context, tenant string) ([]postgres_entity.OAuthTokenEntity, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "OAuthTokenRepository.GetByTenant")
	defer spans.Finish()

	var entities []postgres_entity.OAuthTokenEntity

	err := repo.db.
		Where("tenant_name = ?", tenant).
		Find(&entities).Error

	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.count", len(entities))
	return entities, nil
}

func (repo oAuthTokenRepository) GetByEmailAndProvider(ctx context.Context, tenant, provider, email string) (*postgres_entity.OAuthTokenEntity, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "OAuthTokenRepository.GetByEmailAndProvider")
	defer spans.Finish()

	spans.LogKV("tenant", tenant, "provider", provider, "email", email)

	var oAuthTokenEntity postgres_entity.OAuthTokenEntity

	err := repo.db.
		Where("LOWER(provider) = LOWER(?)", provider).
		Where("tenant_name = ?", tenant).
		Where("LOWER(email_address) = LOWER(?)", email).
		First(&oAuthTokenEntity).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		spans.TraceError(err)
		return nil, err
	}

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		spans.LogKV("result.found", false)
		return nil, nil
	}

	spans.LogKV("result.found", true)

	return &oAuthTokenEntity, nil
}

func (repo oAuthTokenRepository) GetByEmail(ctx context.Context, tenant, email string) (*postgres_entity.OAuthTokenEntity, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "OAuthTokenRepository.GetByEmail")
	defer spans.Finish()
	spans.LogKV("email", email)

	var oAuthTokenEntity postgres_entity.OAuthTokenEntity

	err := repo.db.
		Where("tenant_name = ?", tenant).
		Where("LOWER(email_address) = LOWER(?)", email).
		First(&oAuthTokenEntity).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		spans.TraceError(err)
		return nil, err
	}

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		spans.LogKV("result.found", false)
		return nil, nil
	}

	spans.LogKV("result.found", true)

	return &oAuthTokenEntity, nil
}

func (repo oAuthTokenRepository) GetByPlayerId(ctx context.Context, tenant, provider, playerId string) (*postgres_entity.OAuthTokenEntity, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "OAuthTokenRepository.GetByPlayerId")
	defer spans.Finish()

	spans.LogKV("playerId", playerId, "provider", provider)

	var oAuthTokenEntity postgres_entity.OAuthTokenEntity

	err := repo.db.
		Where("tenant_name = ?", tenant).
		Where("player_identity_id = ?", playerId).
		Where("provider = ?", provider).
		First(&oAuthTokenEntity).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		spans.TraceError(err)
		return nil, err
	}

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		spans.LogKV("result.found", false)
		return nil, nil
	}

	spans.LogKV("result.found", true)

	return &oAuthTokenEntity, nil
}

func (repo oAuthTokenRepository) Save(ctx context.Context, oAuthToken postgres_entity.OAuthTokenEntity) (*postgres_entity.OAuthTokenEntity, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "OAuthTokenRepository.Save")
	defer spans.Finish()

	spans.LogObjectAsJson("oAuthToken", oAuthToken)

	result := repo.db.Save(&oAuthToken)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, fmt.Errorf("saving oauth token failed: %w", result.Error)
	}
	return &oAuthToken, nil
}

func (repo oAuthTokenRepository) Update(ctx context.Context, tenant, playerId, provider, accessToken, refreshToken string, expiresAt time.Time) (*postgres_entity.OAuthTokenEntity, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "OAuthTokenRepository.Update")
	defer spans.Finish()

	spans.LogKV("playerId", playerId, "provider", provider, "expiresAt", expiresAt.String())

	existing, err := repo.GetByPlayerId(ctx, tenant, provider, playerId)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	if existing == nil {
		err := fmt.Errorf("oauth token not found")
		spans.TraceError(err)
		return nil, err
	}

	existing.AccessToken = accessToken
	existing.RefreshToken = refreshToken
	existing.ExpiresAt = expiresAt

	result := repo.db.Save(&existing)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, fmt.Errorf("updating oauth token failed: %w", result.Error)
	}

	return existing, nil
}

func (repo oAuthTokenRepository) MarkForManualRefresh(ctx context.Context, tenant, playerId, provider string) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "OAuthTokenRepository.MarkForManualRefresh")
	defer spans.Finish()

	spans.LogKV("playerId", playerId, "provider", provider)

	existing, err := repo.GetByPlayerId(ctx, tenant, provider, playerId)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	if existing == nil {
		err := fmt.Errorf("oauth token not found")
		spans.TraceError(err)
		return err
	}

	existing.NeedsManualRefresh = true

	result := repo.db.Save(&existing)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return fmt.Errorf("updating oauth token failed: %w", result.Error)
	}

	return nil
}

func (repo oAuthTokenRepository) DeleteByEmail(ctx context.Context, tenant, provider, email string) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "OAuthTokenRepository.DeleteByEmail")
	defer spans.Finish()

	spans.LogKV("provider", provider, "email", email)

	existing, err := repo.GetByEmailAndProvider(ctx, tenant, provider, email)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	if existing == nil {
		return nil
	}

	err = repo.db.Delete(&existing).Error
	if err != nil {
		spans.TraceError(err)
		return fmt.Errorf("deleting oauth token failed: %w", err)
	}

	return nil
}
