package model

import (
	"time"

	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
)

// Contract represents the state of a contract aggregate.
type Contract struct {
	ID               string                       `json:"id"`
	Tenant           string                       `json:"tenant"`
	OrganizationId   string                       `json:"organizationId"`
	Name             string                       `json:"name"`
	CreatedByUserId  string                       `json:"createdByUserId"`
	CreatedAt        time.Time                    `json:"createdAt"`
	UpdatedAt        time.Time                    `json:"updatedAt"`
	ServiceStartedAt *time.Time                   `json:"serviceStartedAt,omitempty"`
	SignedAt         *time.Time                   `json:"signedAt,omitempty"`
	EndedAt          *time.Time                   `json:"endedAt,omitempty"`
	RenewalCycle     string                       `json:"renewalCycle"`
	Status           string                       `json:"status"`
	Source           commonmodel.Source           `json:"source"`
	ExternalSystems  []commonmodel.ExternalSystem `json:"externalSystems"`
}

type ContractDataFields struct {
	OrganizationId       string
	Name                 string
	CreatedByUserId      string
	ServiceStartedAt     *time.Time
	SignedAt             *time.Time
	EndedAt              *time.Time
	RenewalCycle         RenewalCycle
	Status               ContractStatus
	Source               commonmodel.Source
	ExternalSystemFields commonmodel.ExternalSystem
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// ContractStatus represents the status of a contract.
type ContractStatus int32

const (
	Draft ContractStatus = iota
	Live
	Ended
)

// RenewalCycle represents the renewal cycle of a contract.
type RenewalCycle int32

const (
	None RenewalCycle = iota
	MonthlyRenewal
	AnnuallyRenewal
)

// This function provides a string representation of the ContractStatus enum.
func (cs ContractStatus) String() string {
	switch cs {
	case Draft:
		return "DRAFT"
	case Live:
		return "LIVE"
	case Ended:
		return "ENDED"
	default:
		return ""
	}
}

// This function provides a string representation of the RenewalCycle enum.
func (rc RenewalCycle) String() string {
	switch rc {
	case None:
		return "NONE"
	case MonthlyRenewal:
		return "MONTHLY_RENEWAL"
	case AnnuallyRenewal:
		return "ANNUALLY_RENEWAL"
	default:
		return ""
	}
}
