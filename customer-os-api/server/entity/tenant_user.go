package entity

import (
	"fmt"
	"time"
)

type TenantUserEntity struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
	CreatedAt time.Time
}

func (tenantUser TenantUserEntity) ToString() string {
	return fmt.Sprintf("id: %s\nfirstName: %s\nlastName: %s", tenantUser.Id, tenantUser.FirstName, tenantUser.LastName)
}

type TenantUserEntities []TenantUserEntity
