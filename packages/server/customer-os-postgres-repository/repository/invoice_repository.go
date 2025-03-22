package postgres_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type InvoiceRepository interface {
	Reserve(ctx context.Context, invoiceNumber postgres_entity.InvoiceNumberEntity) error
}

type invoiceRepository struct {
	gormDb *gorm.DB
}

func NewInvoiceRepository(gormDb *gorm.DB) InvoiceRepository {
	repo := invoiceRepository{gormDb: gormDb}
	return &repo
}

func (r *invoiceRepository) Reserve(ctx context.Context, invoiceNumber postgres_entity.InvoiceNumberEntity) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InvoiceRepository.Reserve")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	err := r.gormDb.Create(&invoiceNumber).Error
	return err
}
