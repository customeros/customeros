package postgres_repository

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type EmailLookupRepository interface {
	GetById(ctx context.Context, id string) (*postgres_entity.EmailLookup, error)
	Create(ctx context.Context, emailLookup postgres_entity.EmailLookup) (*postgres_entity.EmailLookup, error)
}

type emailLookupRepository struct {
	gormDb *gorm.DB
}

func NewEmailLookupRepository(gormDb *gorm.DB) EmailLookupRepository {
	return &emailLookupRepository{gormDb: gormDb}
}

func (e emailLookupRepository) GetById(ctx context.Context, id string) (*postgres_entity.EmailLookup, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "EmailLookupRepository.GetById")
	defer spans.Finish()

	spans.LogKV("id", id)

	var result postgres_entity.EmailLookup
	err := e.gormDb.
		Where("id = ?", id).
		First(&result).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			spans.LogKV("result.found", false)
			return nil, nil
		}
		spans.TraceError(err)
		return nil, errors.Wrap(err, "failed to get email lookup by ID")
	}

	spans.LogKV("result.found", true)
	return &result, nil
}

func (e emailLookupRepository) Create(ctx context.Context, emailLookup postgres_entity.EmailLookup) (*postgres_entity.EmailLookup, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "EmailLookupRepository.Create")
	defer spans.Finish()

	emailLookup.ID = utils.GenerateRandomString(64)

	err := e.gormDb.Create(&emailLookup).Error

	if err != nil {
		spans.TraceError(err)
		return nil, errors.Wrap(err, "failed to store email lookup")
	}

	return &emailLookup, nil
}
