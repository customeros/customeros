package repository

import (
	"gorm.io/gorm"
)

type InvoiceRepository interface {
	GetNextInvoiceNumberSequenceValue() int
}

type invoiceRepository struct {
	gormDb *gorm.DB
}

func NewInvoiceRepository(gormDb *gorm.DB) InvoiceRepository {
	repo := invoiceRepository{gormDb: gormDb}
	repo.createInvoiceNumberSequence()
	return &repo
}

func (r *invoiceRepository) createInvoiceNumberSequence() {
	// Create the sequence if it doesn't exist
	createSequenceSQL := `
        CREATE SEQUENCE IF NOT EXISTS invoice_number_sequence
        START 1
        INCREMENT 1
        MAXVALUE 999999
        CYCLE
    `
	result := r.gormDb.Exec(createSequenceSQL)
	if result.Error != nil {
		panic(result.Error)
	}
}

func (r *invoiceRepository) GetNextInvoiceNumberSequenceValue() int {
	var nextValue int
	query := "SELECT nextval('invoice_number_sequence')"
	r.gormDb.Raw(query).Scan(&nextValue)
	return nextValue
}
