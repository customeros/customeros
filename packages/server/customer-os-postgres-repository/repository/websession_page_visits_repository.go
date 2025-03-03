package repository

import postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

type WebSessionPageVisitsRepository interface {
	Create(ctx context.Context, pageVisit postgres_entity.WebSessionPageVisit) (*postgres_entity.WebSessionPageVisit, error)
}

type webSessionPageVisitRepository struct {
	gormDb *gorm.DB
}

func NewWebSessionPageVisitRepository(gormDb *gorm.DB) WebSessionPageVisitRepository {
	return &webSessionPageVisitRepository{gormDb: gormDb}
}
