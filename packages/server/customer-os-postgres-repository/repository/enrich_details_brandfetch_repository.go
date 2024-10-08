package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type enrichDetailsBrandfetchRepository struct {
	db *gorm.DB
}

type EnrichDetailsBrandfetchRepository interface {
	Create(ctx context.Context, data entity.EnrichDetailsBrandfetch) (*entity.EnrichDetailsBrandfetch, error)
	GetAllSuccessByDomain(ctx context.Context, domain string) ([]entity.EnrichDetailsBrandfetch, error)
	GetLatestByDomain(ctx context.Context, domain string) (*entity.EnrichDetailsBrandfetch, error)
}

func NewEnrichDetailsBrandfetchRepository(gormDb *gorm.DB) EnrichDetailsBrandfetchRepository {
	return &enrichDetailsBrandfetchRepository{db: gormDb}
}

func (r enrichDetailsBrandfetchRepository) GetAllSuccessByDomain(ctx context.Context, domain string) ([]entity.EnrichDetailsBrandfetch, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "enrichDetailsBrandfetchRepository.GetAllSuccessByDomain")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogKV("domain", domain)

	var data []entity.EnrichDetailsBrandfetch
	err := r.db.Where("domain = ? AND success = ?", domain, true).Order("created_at desc").Find(&data).Error
	if err != nil {
		return nil, err
	}

	span.LogFields(tracingLog.Int("result.count", len(data)))

	return data, nil
}

func (r enrichDetailsBrandfetchRepository) GetLatestByDomain(ctx context.Context, domain string) (*entity.EnrichDetailsBrandfetch, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "enrichDetailsBrandfetchRepository.GetLatestByDomain")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogKV("domain", domain)

	var data entity.EnrichDetailsBrandfetch
	err := r.db.Where("domain = ?", domain).Order("created_at desc").First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.LogFields(tracingLog.Bool("result.found", false))
			return nil, nil // Return nil if no record found
		}
		return nil, err // Return other errors as usual
	}

	span.LogFields(tracingLog.Bool("result.found", true))
	return &data, nil
}

func (r enrichDetailsBrandfetchRepository) Create(ctx context.Context, data entity.EnrichDetailsBrandfetch) (*entity.EnrichDetailsBrandfetch, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "enrichDetailsBrandfetchRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.LogObjectAsJson(span, "data", data)

	data.CreatedAt = utils.Now()
	data.UpdatedAt = utils.Now()
	if err := r.db.WithContext(ctx).Create(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}
