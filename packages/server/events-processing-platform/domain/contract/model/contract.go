package model

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"

	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
)

// Contract represents the state of a contract aggregate.
type Contract struct {
	ID                     string                       `json:"id"`
	Tenant                 string                       `json:"tenant"`
	OrganizationId         string                       `json:"organizationId"`
	Name                   string                       `json:"name"`
	ContractUrl            string                       `json:"contractUrl"`
	CreatedByUserId        string                       `json:"createdByUserId"`
	CreatedAt              time.Time                    `json:"createdAt"`
	UpdatedAt              time.Time                    `json:"updatedAt"`
	ServiceStartedAt       *time.Time                   `json:"serviceStartedAt,omitempty"`
	SignedAt               *time.Time                   `json:"signedAt,omitempty"`
	EndedAt                *time.Time                   `json:"endedAt,omitempty"`
	RenewalCycle           string                       `json:"renewalCycle"`
	RenewalPeriods         *int64                       `json:"renewalPeriods"`
	Status                 string                       `json:"status"`
	Source                 commonmodel.Source           `json:"source"`
	ExternalSystems        []commonmodel.ExternalSystem `json:"externalSystems"`
	Currency               string                       `json:"currency"`
	BillingCycle           string                       `json:"billingCycle"`
	InvoicingStartDate     *time.Time                   `json:"invoicingStartDate,omitempty"`
	AddressLine1           string                       `json:"addressLine1"`
	AddressLine2           string                       `json:"addressLine2"`
	Locality               string                       `json:"locality"`
	Country                string                       `json:"country"`
	Zip                    string                       `json:"zip"`
	OrganizationLegalName  string                       `json:"organizationLegalName"`
	InvoiceEmail           string                       `json:"invoiceEmail"`
	InvoiceNote            string                       `json:"invoiceNote"`
	NextInvoiceDate        *time.Time                   `json:"nextInvoiceDate,omitempty"`
	CanPayWithCard         bool                         `json:"canPayWithCard"`
	CanPayWithDirectDebit  bool                         `json:"canPayWithDirectDebit"`
	CanPayWithBankTransfer bool                         `json:"canPayWithBankTransfer"`
	InvoicingEnabled       bool                         `json:"invoicingEnabled"`
	Removed                bool                         `json:"removed"`
	PayOnline              bool                         `json:"payOnline"`
	PayAutomatically       bool                         `json:"payAutomatically"`
	AutoRenew              bool                         `json:"autoRenew"`
}

type ContractDataFields struct {
	OrganizationId         string
	Name                   string
	ContractUrl            string
	CreatedByUserId        string
	ServiceStartedAt       *time.Time
	SignedAt               *time.Time
	EndedAt                *time.Time
	RenewalCycle           string
	RenewalPeriods         *int64
	Status                 ContractStatus
	BillingCycle           string
	Currency               string
	InvoicingStartDate     *time.Time
	NextInvoiceDate        *time.Time
	AddressLine1           string `json:"addressLine1"`
	AddressLine2           string `json:"addressLine2"`
	Locality               string `json:"locality"`
	Country                string `json:"country"`
	Zip                    string `json:"zip"`
	OrganizationLegalName  string `json:"organizationLegalName"`
	InvoiceEmail           string `json:"invoiceEmail"`
	InvoiceNote            string `json:"invoiceNote"`
	CanPayWithCard         bool   `json:"canPayWithCard"`
	CanPayWithDirectDebit  bool   `json:"canPayWithDirectDebit"`
	CanPayWithBankTransfer bool   `json:"canPayWithBankTransfer"`
	InvoicingEnabled       bool   `json:"invoicingEnabled"`
	PayOnline              bool   `json:"payOnline"`
	PayAutomatically       bool   `json:"payAutomatically"`
	AutoRenew              bool   `json:"autoRenew"`
}

// ContractStatus represents the status of a contract.
type ContractStatus int32

const (
	Draft ContractStatus = iota
	Live
	Ended
)

type ContractStatusString string

const (
	ContractStatusStringLive  ContractStatusString = "LIVE"
	ContractStatusStringEnded ContractStatusString = "ENDED"
)

// RenewalCycle represents the renewal cycle of a contract.
type RenewalCycle int32

const (
	NoneRenewal RenewalCycle = iota
	MonthlyRenewal
	AnnuallyRenewal
	QuarterlyRenewal
)

// This function provides a string representation of the RenewalCycle enum.
func (rc RenewalCycle) String() string {
	switch rc {
	case NoneRenewal:
		return ""
	case MonthlyRenewal:
		return enum.RenewalCycleMonthlyRenewal.String()
	case QuarterlyRenewal:
		return enum.RenewalCycleQuarterlyRenewal.String()
	case AnnuallyRenewal:
		return enum.RenewalCycleAnnualRenewal.String()
	default:
		return ""
	}
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
		return string(enum.BillingCycleMonthlyBilling)
	case QuarterlyBilling:
		return string(enum.BillingCycleQuarterlyBilling)
	case AnnuallyBilling:
		return string(enum.BillingCycleAnnuallyBilling)
	default:
		return ""
	}
}

// This function provides a string representation of the ContractStatus enum.
func (cs ContractStatus) String() string {
	switch cs {
	case Draft:
		return string(enum.ContractStatusDraft)
	case Live:
		return string(enum.ContractStatusLive)
	case Ended:
		return string(enum.ContractStatusEnded)
	default:
		return ""
	}
}
