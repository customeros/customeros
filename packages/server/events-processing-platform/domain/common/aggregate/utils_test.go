package aggregate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetAggregateTenant(t *testing.T) {
	invoiceAggregate := NewCommonAggregateWithTenantAndId("invoice", "tenantName", "invoiceId")
	organizationAggregate := NewCommonAggregateWithTenantAndId("organization", "tenantName", "invoiceId")

	invoiceTenant := GetTenantFromAggregate(invoiceAggregate.GetID(), "invoice")
	organizationTenant := GetTenantFromAggregate(organizationAggregate.GetID(), "organization")

	assert.Equal(t, invoiceTenant, "tenantName")
	assert.Equal(t, organizationTenant, "tenantName")
}
