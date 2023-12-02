package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

type ContractEntity struct {
	Id               string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Source           DataSource
	SourceOfTruth    DataSource
	AppSource        string
	Name             string
	ContractUrl      string
	Status           string
	RenewalCycle     string
	RenewalPeriods   *int64
	SignedAt         *time.Time
	ServiceStartedAt *time.Time
	EndedAt          *time.Time
}

func (c ContractEntity) IsEnded() bool {
	return c.EndedAt != nil && c.EndedAt.Before(utils.Now())
}
