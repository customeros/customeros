package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"time"
)

type ServiceLineItemUpdateEvent struct {
	Tenant    string        `json:"tenant" validate:"required"`
	Name      string        `json:"name"`
	Quantity  int64         `json:"quantity,omitempty" validate:"min=0"`
	Price     float64       `json:"price,omitempty"`
	UpdatedAt time.Time     `json:"updatedAt"`
	Billed    string        `json:"billed"`
	Source    common.Source `json:"source"`
	Comments  string        `json:"comments"`
	VatRate   float64       `json:"vatRate"`
	StartedAt *time.Time    `json:"startedAt,omitempty"`
}
