package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetAggregateTenant(t *testing.T) {
	invoiceAggregate := eventstore.NewCommonAggregateWithTenantAndId("invoice", "tenantName", "invoiceId")
	organizationAggregate := eventstore.NewCommonAggregateWithTenantAndId("organization", "tenantName", "invoiceId")

	invoiceTenant := GetTenantFromAggregate(invoiceAggregate.GetID(), "invoice")
	organizationTenant := GetTenantFromAggregate(organizationAggregate.GetID(), "organization")

	assert.Equal(t, invoiceTenant, "tenantName")
	assert.Equal(t, organizationTenant, "tenantName")
}
