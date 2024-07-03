package caches

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

type Cache struct {
	organizationsByTenantCache map[string][]*model.Organization
}

func NewCache() *Cache {
	cache := Cache{
		organizationsByTenantCache: make(map[string][]*model.Organization),
	}
	return &cache
}

func (c *Cache) SetOrganizations(tenant string, organizations []*model.Organization) {
	if c.organizationsByTenantCache[tenant] != nil {
		c.organizationsByTenantCache[tenant] = nil
	}

	c.organizationsByTenantCache[tenant] = organizations
}

func (c *Cache) GetOrganizations(tenant string) []*model.Organization {
	return c.organizationsByTenantCache[tenant]
}
