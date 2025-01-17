package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type industryMappingRepository struct {
	gormDb *gorm.DB
}

type IndustryMappingRepository interface {
	GetAllIndustryMappingsAsMap(ctx context.Context) (map[string]string, error)
}

func NewIndustryMappingRepository(gormDb *gorm.DB) IndustryMappingRepository {
	return &industryMappingRepository{gormDb: gormDb}
}

func (r industryMappingRepository) GetAllIndustryMappingsAsMap(ctx context.Context) (map[string]string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "IndustryMappingRepository.GetAllIndustryMappingsAsMap")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var industryMappings []entity.IndustryMapping
	err := r.gormDb.Find(&industryMappings).Error
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get industry mappings"))
		return nil, err
	}

	industryMappingsMap := make(map[string]string)
	for _, industryMapping := range industryMappings {
		industryMappingsMap[industryMapping.InputIndustry] = industryMapping.OutputIndustry
	}

	return industryMappingsMap, nil
}
