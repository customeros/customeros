package postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
	"time"
)

type OAuthTokenRepository interface {
	GetAll(ctx context.Context) ([]postgres_entity.OAuthTokenEntity, error)
	GetByTenant(ctx context.Context, tenant string) ([]postgres_entity.OAuthTokenEntity, error)
	GetByEmail(ctx context.Context, tenant, provider, email string) (*postgres_entity.OAuthTokenEntity, error)
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "OAuthTokenRepository.GetAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var entities []postgres_entity.OAuthTokenEntity

	err := repo.db.Where("needs_manual_refresh = ?", false).Find(&entities).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return entities, nil
}

func (repo oAuthTokenRepository) GetByTenant(ctx context.Context, tenant string) ([]postgres_entity.OAuthTokenEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OAuthTokenRepository.GetByTenant")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	var entities []postgres_entity.OAuthTokenEntity

	err := repo.db.
		Where("tenant_name = ?", tenant).
		Find(&entities).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(log.Int("result.count", len(entities)))
	return entities, nil
}

func (repo oAuthTokenRepository) GetByEmail(ctx context.Context, tenant, provider, email string) (*postgres_entity.OAuthTokenEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OAuthTokenRepository.GetByEmail")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	span.LogFields(log.String("tenant", tenant), log.String("provider", provider), log.String("email", email))

	var oAuthTokenEntity postgres_entity.OAuthTokenEntity

	err := repo.db.
		Where("LOWER(provider) = LOWER(?)", provider).
		Where("tenant_name = ?", tenant).
		Where("email_address = ?", email).
		First(&oAuthTokenEntity).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}

	span.LogFields(log.Bool("result.found", true))

	return &oAuthTokenEntity, nil
}

func (repo oAuthTokenRepository) GetByPlayerId(ctx context.Context, tenant, provider, playerId string) (*postgres_entity.OAuthTokenEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OAuthTokenRepository.GetByPlayerId")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("playerId", playerId), log.String("provider", provider))

	var oAuthTokenEntity postgres_entity.OAuthTokenEntity

	err := repo.db.
		Where("tenant_name = ?", tenant).
		Where("player_identity_id = ?", playerId).
		Where("provider = ?", provider).
		First(&oAuthTokenEntity).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}

	span.LogFields(log.Bool("result.found", true))

	return &oAuthTokenEntity, nil
}

func (repo oAuthTokenRepository) Save(ctx context.Context, oAuthToken postgres_entity.OAuthTokenEntity) (*postgres_entity.OAuthTokenEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OAuthTokenRepository.Save")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "oAuthToken", oAuthToken)

	result := repo.db.Save(&oAuthToken)
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return nil, fmt.Errorf("saving oauth token failed: %w", result.Error)
	}
	return &oAuthToken, nil
}

func (repo oAuthTokenRepository) Update(ctx context.Context, tenant, playerId, provider, accessToken, refreshToken string, expiresAt time.Time) (*postgres_entity.OAuthTokenEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OAuthTokenRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("playerId", playerId), log.String("provider", provider), log.String("expiresAt", expiresAt.String()))

	existing, err := repo.GetByPlayerId(ctx, tenant, provider, playerId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if existing == nil {
		err := fmt.Errorf("oauth token not found")
		tracing.TraceErr(span, err)
		return nil, err
	}

	existing.AccessToken = accessToken
	existing.RefreshToken = refreshToken
	existing.ExpiresAt = expiresAt

	result := repo.db.Save(&existing)
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return nil, fmt.Errorf("updating oauth token failed: %w", result.Error)
	}

	return existing, nil
}

func (repo oAuthTokenRepository) MarkForManualRefresh(ctx context.Context, tenant, playerId, provider string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OAuthTokenRepository.MarkForManualRefresh")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("playerId", playerId), log.String("provider", provider))

	existing, err := repo.GetByPlayerId(ctx, tenant, provider, playerId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if existing == nil {
		err := fmt.Errorf("oauth token not found")
		tracing.TraceErr(span, err)
		return err
	}

	existing.NeedsManualRefresh = true

	result := repo.db.Save(&existing)
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return fmt.Errorf("updating oauth token failed: %w", result.Error)
	}

	return nil
}

func (repo oAuthTokenRepository) DeleteByEmail(ctx context.Context, tenant, provider, email string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OAuthTokenRepository.DeleteByEmail")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("provider", provider), log.String("email", email))

	existing, err := repo.GetByEmail(ctx, tenant, provider, email)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if existing == nil {
		return nil
	}

	err = repo.db.Delete(&existing).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("deleting oauth token failed: %w", err)
	}

	return nil
}
