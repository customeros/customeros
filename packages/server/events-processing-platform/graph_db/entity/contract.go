package entity

import (
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
	Status           string
	RenewalCycle     string
	SignedAt         *time.Time
	ServiceStartedAt *time.Time
	EndedAt          *time.Time
}
