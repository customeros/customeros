package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository/postgres/entity"
	"gorm.io/gorm"
)

type InvoiceRepository interface {
	Reserve(invoiceNumber entity.InvoiceNumberEntity) error
}

type invoiceRepository struct {
	gormDb *gorm.DB
}

func NewInvoiceRepository(gormDb *gorm.DB) InvoiceRepository {
	repo := invoiceRepository{gormDb: gormDb}
	return &repo
}

func (r *invoiceRepository) Reserve(invoiceNumber entity.InvoiceNumberEntity) error {
	err := r.gormDb.Save(&invoiceNumber).Error
	return err
}
