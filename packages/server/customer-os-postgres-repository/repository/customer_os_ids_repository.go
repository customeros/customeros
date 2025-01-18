package postgres_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type CustomerOsIdsRepository interface {
	Reserve(ctx context.Context, customerOsIds postgres_entity.CustomerOsIds) error
}

type customerOsIdsRepository struct {
	gormDb *gorm.DB
}

func NewCustomerOsIdsRepository(gormDb *gorm.DB) CustomerOsIdsRepository {
	return &customerOsIdsRepository{gormDb: gormDb}
}

func (repo *customerOsIdsRepository) Reserve(ctx context.Context, customerOsIds postgres_entity.CustomerOsIds) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "CustomerOsIdsRepository.Reserve")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := repo.gormDb.Save(&customerOsIds).Error
	return err
}
