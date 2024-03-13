package order

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"time"
)

type Order struct {
	ID             string             `json:"id"`
	Tenant         string             `json:"tenant"`
	OrganizationId string             `json:"organizationId"`
	CreatedAt      time.Time          `json:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt"`
	SourceFields   commonmodel.Source `json:"source"`

	ConfirmedAt time.Time `json:"confirmedAt"`
	PaidAt      time.Time `json:"paidAt"`
	FulfilledAt time.Time `json:"fulfilledAt"`
	CanceledAt  time.Time `json:"canceledAt"`
}
