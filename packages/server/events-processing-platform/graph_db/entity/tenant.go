package entity

import (
	"fmt"
	"time"
)

type TenantEntity struct {
	Id        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Source    DataSource
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
