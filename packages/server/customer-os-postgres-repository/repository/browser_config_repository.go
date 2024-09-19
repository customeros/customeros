package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type BrowserConfigRepository interface {
	Get(ctx context.Context) ([]entity.BrowserConfig, error)
}

type browserConfigRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewBrowserConfigRepository(gormDb *gorm.DB) BrowserConfigRepository {
	return &browserConfigRepositoryImpl{gormDb: gormDb}
}

func (repo *browserConfigRepositoryImpl) Get(ctx context.Context) ([]entity.BrowserConfig, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "BrowserConfigRepository.Get")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var result []entity.BrowserConfig
	err := repo.gormDb.Where("session_status = 'VALID'").Find(&result).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return result, nil
}
