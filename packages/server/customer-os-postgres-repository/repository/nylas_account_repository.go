package postgres_repository

import (
	"context"
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"gorm.io/gorm"
)

type NylasAccountRepository interface {
	GetByTenantAndEmail(ctx context.Context, tenant, email string) (*postgres_entity.NylasAccount, error)
	Save(ctx context.Context, account *postgres_entity.NylasAccount) (*postgres_entity.NylasAccount, error)
	Delete(ctx context.Context, account *postgres_entity.NylasAccount) error
}

type nylasAccountRepository struct {
	db *gorm.DB
}

func NewNylasAccountRepository(db *gorm.DB) NylasAccountRepository {
	return &nylasAccountRepository{
		db: db,
	}
}

func (r *nylasAccountRepository) GetByTenantAndEmail(ctx context.Context, tenant, email string) (*postgres_entity.NylasAccount, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "NylasAccountRepository.GetByTenantAndEmail")
	defer spans.Finish()
	spans.LogKV("tenant", tenant, "email", email)

	var account postgres_entity.NylasAccount
	err := r.db.
		Where("tenant = ? AND email = ?", tenant, email).
		First(&account).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			spans.LogKV("result.found", false)
			return nil, nil
		}
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.found", true)
	return &account, nil
}

func (r *nylasAccountRepository) Save(ctx context.Context, account *postgres_entity.NylasAccount) (*postgres_entity.NylasAccount, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "NylasAccountRepository.Save")
	defer spans.Finish()
	spans.LogObjectAsJson("account", account)

	err := r.db.Save(account).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return account, nil
}

func (r *nylasAccountRepository) Delete(ctx context.Context, account *postgres_entity.NylasAccount) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "NylasAccountRepository.Delete")
	defer spans.Finish()
	spans.LogObjectAsJson("account", account)

	return r.db.Delete(account).Error
}
