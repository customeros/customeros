package data_fields

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type ContractSaveFields struct {
	CreatedAt              *time.Time          `json:"createdAt,omitempty"`
	AppSource              *string             `json:"appSource,omitempty"`
	Source                 *string             `json:"source,omitempty"`
	OrganizationId         *string             `json:"organizationId,omitempty"`
	CreatedByUserId        *string             `json:"createdByUserId,omitempty"`
	Name                   *string             `json:"name,omitempty"`
	ContractUrl            *string             `json:"contractUrl,omitempty"`
	ServiceStartedAt       *time.Time          `json:"serviceStartedAt,omitempty"`
	SignedAt               *time.Time          `json:"signedAt,omitempty"`
	LengthInMonths         *int64              `json:"lengthInMonths,omitempty"`
	Status                 *string             `json:"status,omitempty"`
	BillingCycleInMonths   *int64              `json:"billingCycleInMonths,omitempty"`
	Currency               *neo4jenum.Currency `json:"currency,omitempty"`
	InvoicingStartDate     *time.Time          `json:"invoicingStartDate,omitempty"`
	InvoicingEnabled       *bool               `json:"invoicingEnabled,omitempty"`
	PayOnline              *bool               `json:"payOnline,omitempty"`
	PayAutomatically       *bool               `json:"payAutomatically,omitempty"`
	CanPayWithCard         *bool               `json:"canPayWithCard,omitempty"`
	CanPayWithDirectDebit  *bool               `json:"canPayWithDirectDebit,omitempty"`
	CanPayWithBankTransfer *bool               `json:"canPayWithBankTransfer,omitempty"`
	AutoRenew              *bool               `json:"autoRenew,omitempty"`
	Check                  *bool               `json:"check,omitempty"`
	DueDays                *int64              `json:"dueDays,omitempty"`
	Country                *string             `json:"country,omitempty"`
	Approved               *bool               `json:"approved,omitempty"`
}

func (c ContractSaveFields) GetCreatedByUserId() string {
	return utils.IfNotNilString(c.CreatedByUserId)
}

func (c ContractSaveFields) GetOrganizationId() string {
	return utils.IfNotNilString(c.OrganizationId)
}
