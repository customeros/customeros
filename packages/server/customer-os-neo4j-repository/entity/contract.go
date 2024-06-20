package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type ContractInternalFields struct {
	StatusRenewalRequestedAt      *time.Time
	RolloutRenewalRequestedAt     *time.Time
	NextPreviewInvoiceRequestedAt *time.Time
}

type ContractEntity struct {
	DataLoaderKey
	EventStoreAggregate
	Id                              string
	Name                            string `neo4jDb:"property:name;lookupName:NAME;supportCaseSensitive:true"`
	CreatedAt                       time.Time
	UpdatedAt                       time.Time
	ServiceStartedAt                *time.Time // DateTime
	SignedAt                        *time.Time // DateTime
	EndedAt                         *time.Time `neo4jDb:"property:endedAt;lookupName:ENDED_AT;supportCaseSensitive:false"` // DateTime
	ContractStatus                  enum.ContractStatus
	Source                          DataSource
	SourceOfTruth                   DataSource
	AppSource                       string
	ContractUrl                     string
	InvoicingStartDate              *time.Time // Date only
	NextInvoiceDate                 *time.Time // Date only
	BillingCycleInMonths            int64
	Currency                        enum.Currency
	TriggeredOnboardingStatusChange bool
	AddressLine1                    string
	AddressLine2                    string
	Zip                             string
	Locality                        string
	Country                         string
	Region                          string
	OrganizationLegalName           string
	InvoiceEmail                    string
	InvoiceEmailCC                  []string
	InvoiceEmailBCC                 []string
	InvoiceNote                     string
	CanPayWithCard                  bool
	CanPayWithDirectDebit           bool
	CanPayWithBankTransfer          bool
	InvoicingEnabled                bool
	PayOnline                       bool
	PayAutomatically                bool
	AutoRenew                       bool
	Check                           bool
	DueDays                         int64
	ContractInternalFields          ContractInternalFields
	LengthInMonths                  int64
	Approved                        bool
	Ltv                             float64
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

func (c ContractEntity) IsDraft() bool {
	return c.ContractStatus == enum.ContractStatusDraft
}
