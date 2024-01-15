package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type ContractEntity struct {
	Id               string
	Name             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ServiceStartedAt *time.Time
	SignedAt         *time.Time
	EndedAt          *time.Time
	RenewalCycle     enum.RenewalCycle
	RenewalPeriods   *int64
	ContractStatus   enum.ContractStatus
	Source           DataSource
	SourceOfTruth    DataSource
	AppSource        string
	ContractUrl      string
}
