package entity

import "time"

type InvoiceNumberEntity struct {
	InvoiceNumber string    `gorm:"column:invoice_number;size:16;primaryKey"`
	Tenant        string    `gorm:"column:tenant;size:50"`
	CreatedDate   time.Time `gorm:"default:current_timestamp"`
	Attempts      int       `gorm:"column:attempts"`
}

func (InvoiceNumberEntity) TableName() string {
	return "invoice_numbers"
}
