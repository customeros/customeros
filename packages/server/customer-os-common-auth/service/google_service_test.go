package service

import (
	"context"
	postgrest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/test/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestServiceAccountCredentialsExistsForTenant(t *testing.T) {
	tenant := "test_tenant"
	err := postgrest.InsertTenantAPIKey(postgresGormDB, tenant, "GSUITE_SERVICE_PRIVATE_KEY", "test1")
	require.NoError(t, err)
	err = postgrest.InsertTenantAPIKey(postgresGormDB, tenant, "GSUITE_SERVICE_EMAIL_ADDRESS", "test2")
	require.NoError(t, err)

	// positive case
	exists, err := serviceContainer.GoogleService.ServiceAccountCredentialsExistsForTenant(context.Background(), tenant)
	require.NoError(t, err)
	assert.True(t, exists)
}
