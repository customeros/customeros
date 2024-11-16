package graph

import (
	"context"
	"testing"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test"
	tenant "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/tenant"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/tenant/event"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	"github.com/stretchr/testify/require"
)

func TestTenantEventHandler_OnUpdateBillingProfileV1(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	profileId := neo4jtest.CreateTenantBillingProfile(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantBillingProfileEntity{})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model.NodeLabelTenant:               1,
		model.NodeLabelTenantBillingProfile: 1,
	})

	// Prepare the event handler
	eventHandler := &TenantEventHandler{
		log:      testLogger,
		services: testDatabase.Services,
	}

	timeNow := utils.Now()

	aggregate := tenant.NewTenantAggregate(tenantName)
	updateEvent, err := event.NewTenantBillingProfileUpdateEvent(
		aggregate,
		profileId,
		&tenantpb.UpdateBillingProfileRequest{
			Phone:                  "phone",
			AddressLine1:           "addressLine1",
			AddressLine2:           "addressLine2",
			AddressLine3:           "addressLine3",
			Locality:               "locality",
			Country:                "country",
			Region:                 "region",
			Zip:                    "zip",
			LegalName:              "legalName",
			VatNumber:              "vatNumber",
			SendInvoicesFrom:       "sendInvoicesFrom",
			SendInvoicesBcc:        "sendInvoicesBcc",
			CanPayWithPigeon:       true,
			CanPayWithBankTransfer: true,
			Check:                  true,
		},
		timeNow,
		[]string{},
	)
	require.Nil(t, err)

	// EXECUTE
	err = eventHandler.OnUpdateBillingProfileV1(context.Background(), updateEvent)
	require.Nil(t, err)

	// check still same nodes available
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		model.NodeLabelTenant:               1,
		model.NodeLabelTenantBillingProfile: 1,
	})

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, model.NodeLabelTenantBillingProfile, profileId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	// verify
	tenantBillingProfileEntity := neo4jmapper.MapDbNodeToTenantBillingProfileEntity(dbNode)
	require.Equal(t, profileId, tenantBillingProfileEntity.Id)
	test.AssertRecentTime(t, tenantBillingProfileEntity.UpdatedAt)
	require.Equal(t, "phone", tenantBillingProfileEntity.Phone)
	require.Equal(t, "addressLine1", tenantBillingProfileEntity.AddressLine1)
	require.Equal(t, "addressLine2", tenantBillingProfileEntity.AddressLine2)
	require.Equal(t, "addressLine3", tenantBillingProfileEntity.AddressLine3)
	require.Equal(t, "locality", tenantBillingProfileEntity.Locality)
	require.Equal(t, "country", tenantBillingProfileEntity.Country)
	require.Equal(t, "region", tenantBillingProfileEntity.Region)
	require.Equal(t, "zip", tenantBillingProfileEntity.Zip)
	require.Equal(t, "legalName", tenantBillingProfileEntity.LegalName)
	require.Equal(t, "vatNumber", tenantBillingProfileEntity.VatNumber)
	require.Equal(t, "sendInvoicesFrom", tenantBillingProfileEntity.SendInvoicesFrom)
	require.Equal(t, "sendInvoicesBcc", tenantBillingProfileEntity.SendInvoicesBcc)
	require.Equal(t, true, tenantBillingProfileEntity.CanPayWithPigeon)
	require.Equal(t, true, tenantBillingProfileEntity.CanPayWithBankTransfer)
	require.Equal(t, true, tenantBillingProfileEntity.Check)
}
