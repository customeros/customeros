package entity

import (
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
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
	RenewalCycle     RenewalCycle
	RenewalPeriods   *int64
	ContractStatus   ContractStatus
	Source           neo4jentity.DataSource
	SourceOfTruth    neo4jentity.DataSource
	AppSource        string
	ContractUrl      string

	DataloaderKey string
}

type ContractEntities []ContractEntity
