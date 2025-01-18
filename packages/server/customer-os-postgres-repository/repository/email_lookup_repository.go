package postgres_repository

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailLookupRepository.GetById")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	span.LogFields(tracingLog.String("id", id))

	var result postgres_entity.EmailLookup
	err := e.gormDb.
		Where("id = ?", id).
		First(&result).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			span.LogFields(tracingLog.Bool("result.found", false))
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "failed to get email lookup by ID")
	}

	span.LogFields(tracingLog.Bool("result.found", true))
	return &result, nil
}

func (e emailLookupRepository) Create(ctx context.Context, emailLookup postgres_entity.EmailLookup) (*postgres_entity.EmailLookup, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailLookupRepository.Create")
	defer span.Finish()
	tracing.TagTenant(span, emailLookup.Tenant)
	tracing.TagComponentPostgresRepository(span)

	emailLookup.ID = utils.GenerateRandomString(64)

	err := e.gormDb.Create(&emailLookup).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "failed to store email lookup")
	}

	return &emailLookup, nil
}
