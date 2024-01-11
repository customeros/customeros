package invoice

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"time"
)

type Invoice struct {
	ID           string             `json:"id"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
	Date         time.Time          `json:"date"`
	DueDate      time.Time          `json:"dueDate"`
	Lines        []InvoiceLine      `json:"invoiceLines"`
	TotalAmount  float64            `json:"amount"`
	SourceFields commonmodel.Source `json:"source"`
}

type InvoiceLine struct {
	ID           string             `json:"id"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
	Name         string             `json:"name"`
	Price        float64            `json:"price"`
	Quantity     int                `json:"quantity"`
	Total        float64            `json:"amount"`
	SourceFields commonmodel.Source `json:"source"`
}
