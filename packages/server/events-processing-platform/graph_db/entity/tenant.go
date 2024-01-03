package entity

import (
	"fmt"
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
	"time"
)

type TenantEntity struct {
	Id        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Source    neo4jentity.DataSource
	AppSource string

	DataloaderKey string
}

func (domain TenantEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", domain.Id, domain.Name)
}

type TenantEntities []TenantEntity

func (domain TenantEntity) Labels(tenant string) []string {
	return []string{"Tenant"}
}
