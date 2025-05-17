package postgres_repository

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
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
	spans, ctx := telemetry.StartPostgresSpan(ctx, "WebSessionPageVisitRepository.Create")
	defer spans.Finish()

	spans.LogObjectAsJson("pageVisit", pageVisit)

	var created postgres_entity.WebSessionPageVisit
	err := r.gormDb.Create(&pageVisit).Scan(&created).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return &created, nil
}
