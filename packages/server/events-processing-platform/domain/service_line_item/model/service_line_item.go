package model

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"time"
)

// ServiceLineItem represents the state of a service line item aggregate.
type ServiceLineItem struct {
	ID          string    `json:"id"`
	ContractId  string    `json:"contractId"`
	Billed      string    `json:"billed"`
	Licenses    int32     `json:"licenses"` // Relevant only for Subscription type
	Price       float32   `json:"price"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Source      commonmodel.Source
}

// ServiceLineItemDataFields contains all the fields that may be used to create or update a service line item.
type ServiceLineItemDataFields struct {
	Billed      BilledType `json:"billed"`
	Licenses    int32      `json:"licenses"` // Relevant only for Subscription type
	Price       float32    `json:"price"`
	Description string     `json:"description"`
	ContractId  string     `json:"contractId"`
}

// BilledType enum represents the billing type for a service line item.
type BilledType int32

const (
	MonthlyBilled BilledType = iota
	AnnuallyBilled
	OnceBilled // For One-Time
)

func (bt BilledType) String() string {
	return [...]string{"MONTHLY", "ANNUALLY", "ONCE"}[bt]
}

func (bt BilledType) IsOneTime() bool {
	return bt == OnceBilled
}
