package invoice

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"time"
)

type Invoice struct {
	ID             string             `json:"id"`
	Tenant         string             `json:"tenant"`
	OrganizationId string             `json:"organizationId"`
	CreatedAt      time.Time          `json:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt"`
	SourceFields   commonmodel.Source `json:"source"`

	DryRun           bool          `json:"dryRun"`
	Date             time.Time     `json:"date"`
	DueDate          time.Time     `json:"dueDate"`
	Amount           float64       `json:"amount"`
	VAT              float64       `json:"vat"`
	Total            float64       `json:"total"`
	Lines            []InvoiceLine `json:"invoiceLines"`
	RepositoryFileId string        `json:"repositoryFileId"`
}

type InvoiceLine struct {
	ID           string             `json:"id"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
	SourceFields commonmodel.Source `json:"source"`

	Index    int64   `json:"index"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int64   `json:"quantity"`
	Amount   float64 `json:"amount"`
	VAT      float64 `json:"vat"`
	Total    float64 `json:"total"`
}
