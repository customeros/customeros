package postgres_repository

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type WebSessionPageVisitRepository interface {
	Create(ctx context.Context, pageVisit postgres_entity.WebSessionPageVisit) (*postgres_entity.WebSessionPageVisit, error)
}

type webSessionPageVisitRepository struct {
	gormDb *gorm.DB
}

func NewWebSessionPageVisitRepository(gormDb *gorm.DB) WebSessionPageVisitRepository {
	return &webSessionPageVisitRepository{gormDb: gormDb}
}

func (r *webSessionPageVisitRepository) Create(ctx context.Context, pageVisit postgres_entity.WebSessionPageVisit) (*postgres_entity.WebSessionPageVisit, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionPageVisitRepository.Create")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "pageVisit", pageVisit)

	var created postgres_entity.WebSessionPageVisit
	err := r.gormDb.Create(&pageVisit).Scan(&created).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &created, nil
}
