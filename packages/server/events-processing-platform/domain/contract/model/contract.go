package model

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
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
	Status                 string                       `json:"status"`
	Source                 events.Source                `json:"source"`
	ExternalSystems        []commonmodel.ExternalSystem `json:"externalSystems"`
	Currency               string                       `json:"currency"`
	BillingCycleInMonths   int64                        `json:"billingCycleInMonths"`
	InvoicingStartDate     *time.Time                   `json:"invoicingStartDate,omitempty"`
	AddressLine1           string                       `json:"addressLine1"`
	AddressLine2           string                       `json:"addressLine2"`
	Locality               string                       `json:"locality"`
	Country                string                       `json:"country"`
	Region                 string                       `json:"region"`
	Zip                    string                       `json:"zip"`
	OrganizationLegalName  string                       `json:"organizationLegalName"`
	InvoiceEmail           string                       `json:"invoiceEmail"`
	InvoiceEmailCC         []string                     `json:"invoiceEmailCC"`
	InvoiceEmailBCC        []string                     `json:"invoiceEmailBCC"`
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
	Check                  bool                         `json:"check"`
	DueDays                int64                        `json:"dueDays"`
	LengthInMonths         int64                        `json:"lengthInMonths"`
	Approved               bool                         `json:"approved"`
}

type ContractDataFields struct {
	OrganizationId         string
	Name                   string
	ContractUrl            string
	CreatedByUserId        string
	ServiceStartedAt       *time.Time
	SignedAt               *time.Time
	EndedAt                *time.Time
	BillingCycleInMonths   int64
	Currency               string
	InvoicingStartDate     *time.Time
	NextInvoiceDate        *time.Time
	AddressLine1           string   `json:"addressLine1"`
	AddressLine2           string   `json:"addressLine2"`
	Locality               string   `json:"locality"`
	Country                string   `json:"country"`
	Region                 string   `json:"region"`
	Zip                    string   `json:"zip"`
	OrganizationLegalName  string   `json:"organizationLegalName"`
	InvoiceEmail           string   `json:"invoiceEmail"`
	InvoiceEmailCC         []string `json:"invoiceEmailCC"`
	InvoiceEmailBCC        []string `json:"invoiceEmailBCC"`
	InvoiceNote            string   `json:"invoiceNote"`
	CanPayWithCard         bool     `json:"canPayWithCard"`
	CanPayWithDirectDebit  bool     `json:"canPayWithDirectDebit"`
	CanPayWithBankTransfer bool     `json:"canPayWithBankTransfer"`
	InvoicingEnabled       bool     `json:"invoicingEnabled"`
	PayOnline              bool     `json:"payOnline"`
	PayAutomatically       bool     `json:"payAutomatically"`
	AutoRenew              bool     `json:"autoRenew"`
	Check                  bool     `json:"check"`
	DueDays                int64    `json:"dueDays"`
	LengthInMonths         int64    `json:"lengthInMonths"`
	Approved               bool     `json:"approved"`
}
