package postgres_repository

import (
	"context"
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"gorm.io/gorm"
)

type NylasGrantRepository interface {
	GetByTenantAndEmail(ctx context.Context, tenant, email string) (*postgres_entity.NylasGrant, error)
	GetByTenantAndUserId(ctx context.Context, tenant, userId string) (*postgres_entity.NylasGrant, error)
	Save(ctx context.Context, grant *postgres_entity.NylasGrant) (*postgres_entity.NylasGrant, error)
	DeleteByGrantId(ctx context.Context, grantId string) error
}

type nylasGrantRepository struct {
	db *gorm.DB
}

func NewNylasGrantRepository(db *gorm.DB) NylasGrantRepository {
	return &nylasGrantRepository{
		db: db,
	}
}

func (r *nylasGrantRepository) GetByTenantAndEmail(ctx context.Context, tenant, email string) (*postgres_entity.NylasGrant, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "NylasGrantRepository.GetByTenantAndEmail")
	defer spans.Finish()
	spans.LogKV("tenant", tenant, "email", email)

	var grant postgres_entity.NylasGrant
	err := r.db.
		Where("tenant = ? AND email = ?", tenant, email).
		First(&grant).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			spans.LogKV("result.found", false)
			return nil, nil
		}
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.found", true)
	return &grant, nil
}

func (r *nylasGrantRepository) GetByTenantAndUserId(ctx context.Context, tenant, userId string) (*postgres_entity.NylasGrant, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "NylasGrantRepository.GetByTenantAndUserId")
	defer spans.Finish()
	spans.LogKV("tenant", tenant, "userId", userId)

	var grant postgres_entity.NylasGrant
	err := r.db.
		Where("tenant = ? AND user_id = ?", tenant, userId).
		First(&grant).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			spans.LogKV("result.found", false)
			return nil, nil
		}
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.found", true)
	return &grant, nil
}

func (r *nylasGrantRepository) Save(ctx context.Context, grant *postgres_entity.NylasGrant) (*postgres_entity.NylasGrant, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "NylasGrantRepository.Save")
	defer spans.Finish()
	spans.LogObjectAsJson("grant", grant)

	err := r.db.Save(grant).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return grant, nil
}

func (r *nylasGrantRepository) DeleteByGrantId(ctx context.Context, grantId string) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "NylasGrantRepository.DeleteByGrantId")
	defer spans.Finish()
	spans.TagEntity(grantId)

	err := r.db.Where("nylas_grant_id = ?", grantId).Delete(&postgres_entity.NylasGrant{}).Error
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}
