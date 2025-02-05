package postgres_repository

import (
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"gorm.io/gorm"
)

type InvoiceRepository interface {
	Reserve(invoiceNumber postgres_entity.InvoiceNumberEntity) error
}

type invoiceRepository struct {
	gormDb *gorm.DB
}

func NewInvoiceRepository(gormDb *gorm.DB) InvoiceRepository {
	repo := invoiceRepository{gormDb: gormDb}
	return &repo
}

func (r *invoiceRepository) Reserve(invoiceNumber postgres_entity.InvoiceNumberEntity) error {
	err := r.gormDb.Save(&invoiceNumber).Error
	return err
}
