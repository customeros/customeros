package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type EmailLookupRepository interface {
	GetById(ctx context.Context, id string) (*entity.EmailLookup, error)
	Create(ctx context.Context, emailLookup entity.EmailLookup) (*entity.EmailLookup, error)
}

type emailLookupRepository struct {
	gormDb *gorm.DB
}

func NewEmailLookupRepository(gormDb *gorm.DB) EmailLookupRepository {
	return &emailLookupRepository{gormDb: gormDb}
}

func (e emailLookupRepository) GetById(ctx context.Context, id string) (*entity.EmailLookup, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TrackingRepository.GetById")
	defer span.Finish()

	span.LogFields(tracingLog.String("id", id))

	var result entity.EmailLookup
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

func (e emailLookupRepository) Create(ctx context.Context, emailLookup entity.EmailLookup) (*entity.EmailLookup, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailLookupRepository.Create")
	defer span.Finish()
	emailLookup.ID = utils.GenerateRandomString(64)

	err := e.gormDb.Create(&emailLookup).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "failed to store email lookup")
	}

	return &emailLookup, nil
}
