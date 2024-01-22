package invoice

import (
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"time"
)

type Invoice struct {
	ID               string                  `json:"id"`
	Tenant           string                  `json:"tenant"`
	ContractId       string                  `json:"contractId"`
	CreatedAt        time.Time               `json:"createdAt"`
	UpdatedAt        time.Time               `json:"updatedAt"`
	SourceFields     commonmodel.Source      `json:"source"`
	DryRun           bool                    `json:"dryRun"`
	InvoiceNumber    string                  `json:"invoiceNumber"`
	Currency         string                  `json:"currency"`
	PeriodStartDate  time.Time               `json:"periodStartDate"`
	PeriodEndDate    time.Time               `json:"periodEndDate"`
	DueDate          time.Time               `json:"dueDate"`
	Amount           float64                 `json:"amount"`
	VAT              float64                 `json:"vat"`
	TotalAmount      float64                 `json:"totalAmount"`
	InvoiceLines     []InvoiceLine           `json:"invoiceLines"`
	RepositoryFileId string                  `json:"repositoryFileId"`
	DryRunLines      []DryRunServiceLineItem `json:"dryRunLines"`
	BillingCycle     string                  `json:"billingCycle"`
	Status           string                  `json:"status"`
}

type DryRunServiceLineItem struct {
	ServiceLineItemId string  `json:"serviceLineItemId"`
	Name              string  `json:"name"`
	Billed            string  `json:"billed"`
	Price             float64 `json:"price"`
	Quantity          int64   `json:"quantity"`
}

type InvoiceLine struct {
	ID                      string             `json:"id"`
	CreatedAt               time.Time          `json:"createdAt"`
	UpdatedAt               time.Time          `json:"updatedAt"`
	SourceFields            commonmodel.Source `json:"source"`
	Name                    string             `json:"name"`
	Price                   float64            `json:"price"`
	Quantity                int64              `json:"quantity"`
	Amount                  float64            `json:"amount"`
	VAT                     float64            `json:"vat"`
	TotalAmount             float64            `json:"totalAmount"`
	ServiceLineItemId       string             `json:"serviceLineItemId"`
	ServiceLineItemParentId string             `json:"serviceLineItemParentId"`
	BilledType              string             `json:"billedType"`
}

// BillingCycle represents the billing cycle of a contract.
type BillingCycle int32

const (
	NoneBilling BillingCycle = iota
	MonthlyBilling
	QuarterlyBilling
	AnnuallyBilling
)

// This function provides a string representation of the BillingCyckle enum.
func (bc BillingCycle) String() string {
	switch bc {
	case NoneBilling:
		return ""
	case MonthlyBilling:
		return string(neo4jenum.BillingCycleMonthlyBilling)
	case QuarterlyBilling:
		return string(neo4jenum.BillingCycleQuarterlyBilling)
	case AnnuallyBilling:
		return string(neo4jenum.BillingCycleAnnuallyBilling)
	default:
		return ""
	}
}

// BilledType enum represents the billing type for a service line item.
type BilledType int32

const (
	NoneBilled BilledType = iota
	MonthlyBilled
	AnnuallyBilled
	OnceBilled  // For One-Time
	UsageBilled // For Usage-Based
	QuarterlyBilled
)

func (bt BilledType) String() string {
	switch bt {
	case NoneBilled:
		return ""
	case MonthlyBilled:
		return string(neo4jenum.BilledTypeMonthly)
	case QuarterlyBilled:
		return string(neo4jenum.BilledTypeQuarterly)
	case AnnuallyBilled:
		return string(neo4jenum.BilledTypeAnnually)
	case OnceBilled:
		return string(neo4jenum.BilledTypeOnce)
	case UsageBilled:
		return string(neo4jenum.BilledTypeUsage)
	default:
		return ""
	}
}

type InvoiceStatus int32

const (
	NoneInvoiceStatus InvoiceStatus = iota
	DraftInvoiceStatus
	DueInvoiceStatus
	PaidInvoiceStatus
)

// This function provides a string representation of the BillingCyckle enum.
func (is InvoiceStatus) String() string {
	switch is {
	case NoneInvoiceStatus:
		return ""
	case DraftInvoiceStatus:
		return string(neo4jenum.InvoiceStatusDraft)
	case DueInvoiceStatus:
		return string(neo4jenum.InvoiceStatusDue)
	case PaidInvoiceStatus:
		return string(neo4jenum.InvoiceStatusPaid)
	default:
		return ""
	}
}
