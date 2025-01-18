package postgres_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type aiLocationMappingRepository struct {
	gormDb *gorm.DB
}

type AiLocationMappingRepository interface {
	AddLocationMapping(ctx context.Context, aiLocationMapping postgres_entity.AiLocationMapping) error
	GetLatestLocationMappingByInput(ctx context.Context, input string) (*postgres_entity.AiLocationMapping, error)
}

func NewAiLocationMappingRepository(gormDb *gorm.DB) AiLocationMappingRepository {
	return &aiLocationMappingRepository{gormDb: gormDb}
}

func (r aiLocationMappingRepository) AddLocationMapping(ctx context.Context, aiLocationMapping postgres_entity.AiLocationMapping) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "AiLocationMappingRepository.AddLocationMapping")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := r.gormDb.Create(&aiLocationMapping).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r aiLocationMappingRepository) GetLatestLocationMappingByInput(ctx context.Context, input string) (*postgres_entity.AiLocationMapping, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "AiLocationMappingRepository.GetLatestLocationMappingByInput")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var aiLocationMapping postgres_entity.AiLocationMapping
	err := r.gormDb.Where("input = ?", input).Order("created_at desc").First(&aiLocationMapping).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &aiLocationMapping, nil
}
