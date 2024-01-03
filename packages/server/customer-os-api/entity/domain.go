package entity

import (
	"fmt"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

type DomainEntity struct {
	Id        string
	Domain    string `neo4jDb:"property:domain;lookupName:DOMAIN;supportCaseSensitive:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Source    neo4jentity.DataSource
	AppSource string

	DataloaderKey string
}

func (domain DomainEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", domain.Id, domain.Domain)
}

type DomainEntities []DomainEntity

func (domain DomainEntity) Labels(tenant string) []string {
	return []string{"Domain"}
}
