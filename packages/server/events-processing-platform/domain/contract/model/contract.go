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
	ContractUrl      string                       `json:"contractUrl"`
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
	OrganizationId   string
	Name             string
	ContractUrl      string
	CreatedByUserId  string
	ServiceStartedAt *time.Time
	SignedAt         *time.Time
	EndedAt          *time.Time
	RenewalCycle     RenewalCycle
	Status           ContractStatus
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
	ContractStatusStringDraft ContractStatusString = "DRAFT"
	ContractStatusStringLive  ContractStatusString = "LIVE"
	ContractStatusStringEnded ContractStatusString = "ENDED"
)

// RenewalCycle represents the renewal cycle of a contract.
type RenewalCycle int32

const (
	None RenewalCycle = iota
	MonthlyRenewal
	AnnuallyRenewal
	QuarterlyRenewal
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
		return ""
	case MonthlyRenewal:
		return string(MonthlyRenewalCycleString)
	case QuarterlyRenewal:
		return string(QuarterlyRenewalCycleString)
	case AnnuallyRenewal:
		return string(AnnuallyRenewalCycleString)
	default:
		return ""
	}
}

type RenewalCycleString string

const (
	MonthlyRenewalCycleString   RenewalCycleString = "MONTHLY"
	QuarterlyRenewalCycleString RenewalCycleString = "QUARTERLY"
	AnnuallyRenewalCycleString  RenewalCycleString = "ANNUALLY"
)

func RenewalCycleFromString(renewalCycle string) RenewalCycle {
	switch renewalCycle {
	case "MONTHLY":
		return MonthlyRenewal
	case "QUARTERLY":
		return QuarterlyRenewal
	case "ANNUALLY":
		return AnnuallyRenewal
	default:
		return None
	}
}

func IsFrequencyBasedRenewalCycle(renewalCycle string) bool {
	return renewalCycle == string(MonthlyRenewalCycleString) ||
		renewalCycle == string(AnnuallyRenewalCycleString) ||
		renewalCycle == string(QuarterlyRenewalCycleString)
}
