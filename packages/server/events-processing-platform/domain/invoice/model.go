package invoice

import (
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"time"
)

const (
	PARAM_INVOICE_NUMBER = "invoiceNumber"
)

type Invoice struct {
	ID                   string                  `json:"id"`
	Tenant               string                  `json:"tenant"`
	ContractId           string                  `json:"contractId"`
	CreatedAt            time.Time               `json:"createdAt"`
	UpdatedAt            time.Time               `json:"updatedAt"`
	SourceFields         events.Source           `json:"source"`
	DryRun               bool                    `json:"dryRun"`
	OffCycle             bool                    `json:"offCycle"`
	Preview              bool                    `json:"preview"`
	Postpaid             bool                    `json:"postpaid"`
	InvoiceNumber        string                  `json:"invoiceNumber"`
	Currency             string                  `json:"currency"`
	PeriodStartDate      time.Time               `json:"periodStartDate"`
	PeriodEndDate        time.Time               `json:"periodEndDate"`
	Amount               float64                 `json:"amount"`
	VAT                  float64                 `json:"vat"`
	TotalAmount          float64                 `json:"totalAmount"`
	InvoiceLines         []InvoiceLine           `json:"invoiceLines"`
	RepositoryFileId     string                  `json:"repositoryFileId"`
	DryRunLines          []DryRunServiceLineItem `json:"dryRunLines"`
	Status               string                  `json:"status"`
	Note                 string                  `json:"note"`
	PaymentLink          string                  `json:"paymentLink"`
	BillingCycleInMonths int64                   `json:"billingCycleInMonths"`
}

type DryRunServiceLineItem struct {
	ServiceLineItemId string  `json:"serviceLineItemId"`
	Name              string  `json:"name"`
	Billed            string  `json:"billed"`
	Price             float64 `json:"price"`
	Quantity          int64   `json:"quantity"`
}

type InvoiceLine struct {
	ID                      string        `json:"id"`
	CreatedAt               time.Time     `json:"createdAt"`
	UpdatedAt               time.Time     `json:"updatedAt"`
	SourceFields            events.Source `json:"source"`
	Name                    string        `json:"name"`
	Price                   float64       `json:"price"`
	Quantity                int64         `json:"quantity"`
	Amount                  float64       `json:"amount"`
	VAT                     float64       `json:"vat"`
	TotalAmount             float64       `json:"totalAmount"`
	ServiceLineItemId       string        `json:"serviceLineItemId"`
	ServiceLineItemParentId string        `json:"serviceLineItemParentId"`
	BilledType              string        `json:"billedType"`
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
	InitializedInvoiceStatus
	DueInvoiceStatus
	PaidInvoiceStatus
	VoidInvoiceStatus
	ScheduledInvoiceStatus
	OverdueInvoiceStatus
	OnHoldInvoiceStatus
	EmptyInvoiceStatus
)

// This function provides a string representation of the BillingCyckle enum.
func (is InvoiceStatus) String() string {
	switch is {
	case NoneInvoiceStatus:
		return ""
	case InitializedInvoiceStatus:
		return string(neo4jenum.InvoiceStatusInitialized)
	case DueInvoiceStatus:
		return string(neo4jenum.InvoiceStatusDue)
	case PaidInvoiceStatus:
		return string(neo4jenum.InvoiceStatusPaid)
	case VoidInvoiceStatus:
		return string(neo4jenum.InvoiceStatusVoid)
	case ScheduledInvoiceStatus:
		return string(neo4jenum.InvoiceStatusScheduled)
	case OverdueInvoiceStatus:
		return string(neo4jenum.InvoiceStatusOverdue)
	case OnHoldInvoiceStatus:
		return string(neo4jenum.InvoiceStatusOnHold)
	case EmptyInvoiceStatus:
		return string(neo4jenum.InvoiceStatusEmpty)
	default:
		return ""
	}
}
