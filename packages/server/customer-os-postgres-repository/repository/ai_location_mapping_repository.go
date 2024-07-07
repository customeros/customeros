package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type aiLocationMappingRepository struct {
	gormDb *gorm.DB
}

type AiLocationMappingRepository interface {
	AddLocationMapping(ctx context.Context, aiLocationMapping entity.AiLocationMapping) error
}

func NewAiLocationMappingRepository(gormDb *gorm.DB) AiLocationMappingRepository {
	return &aiLocationMappingRepository{gormDb: gormDb}
}

func (r aiLocationMappingRepository) AddLocationMapping(ctx context.Context, aiLocationMapping entity.AiLocationMapping) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "AiLocationMappingRepository.AddLocationMapping")
	defer span.Finish()

	err := r.gormDb.Create(&aiLocationMapping).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
