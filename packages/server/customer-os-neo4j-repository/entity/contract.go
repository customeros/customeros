package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type ContractEntity struct {
	Id                 string
	Name               string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	ServiceStartedAt   *time.Time
	SignedAt           *time.Time
	EndedAt            *time.Time
	RenewalCycle       enum.RenewalCycle
	RenewalPeriods     *int64
	ContractStatus     enum.ContractStatus
	Source             DataSource
	SourceOfTruth      DataSource
	AppSource          string
	ContractUrl        string
	InvoicingStartDate *time.Time
	BillingCycle       enum.BillingCycle
	NextInvoiceDate    *time.Time
	Currency           enum.Currency

	TriggeredOnboardingStatusChange bool

	DataloaderKey string
}

type ContractEntities []ContractEntity

func (c ContractEntity) IsEnded() bool {
	return c.EndedAt != nil && c.EndedAt.Before(utils.Now())
}

func (c ContractEntity) IsSigned() bool {
	return c.SignedAt != nil && c.SignedAt.Before(utils.Now())
}

func (c ContractEntity) IsServiceStarted() bool {
	return c.ServiceStartedAt != nil && c.ServiceStartedAt.Before(utils.Now())
}

func (c ContractEntity) IsEligibleToStartOnboarding() bool {
	return !c.TriggeredOnboardingStatusChange && (c.IsSigned() || c.IsServiceStarted()) && !c.IsEnded()
}
