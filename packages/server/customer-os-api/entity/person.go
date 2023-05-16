package entity

import (
	"fmt"
	"time"
)

type PersonRelation string

const (
	IDENTIFIES PersonRelation = "IDENTIFIES"
)

type PersonEntity struct {
	Id            string
	IdentityId    *string
	Email         string
	Provider      string
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
	CreatedAt     time.Time
	UpdatedAt     time.Time

	DataloaderKey string
}

func (person PersonEntity) ToString() string {
	return fmt.Sprintf("id: %s\nemail: %s\nidentityId: %s", person.Id, person.Email, *person.IdentityId)
}

type PersonEntities []PersonEntity

func (PersonEntity) Labels(tenant string) []string {
	return []string{
		NodeLabel_Person,
	}
}
