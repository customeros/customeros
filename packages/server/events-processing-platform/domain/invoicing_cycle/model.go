package invoicing_cycle

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"time"
)

type InvoicingCycleType int32

const (
	DATE InvoicingCycleType = iota
	ANNIVERSARY
)

type InvoicingCycleTypeString string

const (
	InvoicingCycleTypeDate        InvoicingCycleTypeString = "DATE"
	InvoicingCycleTypeAnniversary InvoicingCycleTypeString = "ANNIVERSARY"
)

func (t InvoicingCycleType) StringValue() InvoicingCycleTypeString {
	switch t {
	case DATE:
		return InvoicingCycleTypeDate
	case ANNIVERSARY:
		return InvoicingCycleTypeAnniversary
	default:
		return ""
	}
}

type InvoicingCycle struct {
	ID           string                   `json:"id"`
	Type         InvoicingCycleTypeString `json:"type"`
	CreatedAt    time.Time                `json:"createdAt"`
	UpdatedAt    time.Time                `json:"updatedAt"`
	SourceFields commonmodel.Source       `json:"source"`
}
