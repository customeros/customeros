package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	tenant "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/tenant"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/tenant/event"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTenantEventHandler_OnUpdateBillingProfileV1(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	profileId := neo4jtest.CreateTenantBillingProfile(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantBillingProfileEntity{})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		neo4jutil.NodeLabelTenant:               1,
		neo4jutil.NodeLabelTenantBillingProfile: 1,
	})

	// Prepare the event handler
	eventHandler := &TenantEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	timeNow := utils.Now()

	aggregate := tenant.NewTenantAggregate(tenantName)
	updateEvent, err := event.NewTenantBillingProfileUpdateEvent(
		aggregate,
		profileId,
		&tenantpb.UpdateBillingProfileRequest{
			Phone:                         "phone",
			AddressLine1:                  "addressLine1",
			AddressLine2:                  "addressLine2",
			AddressLine3:                  "addressLine3",
			Locality:                      "locality",
			Country:                       "country",
			Zip:                           "zip",
			LegalName:                     "legalName",
			DomesticPaymentsBankInfo:      "domesticPaymentsBankInfo",
			InternationalPaymentsBankInfo: "internationalPaymentsBankInfo",
			VatNumber:                     "vatNumber",
			SendInvoicesFrom:              "sendInvoicesFrom",
			SendInvoicesBcc:               "sendInvoicesBcc",
			CanPayWithCard:                true,
			CanPayWithDirectDebitSEPA:     true,
			CanPayWithDirectDebitACH:      true,
			CanPayWithDirectDebitBacs:     true,
			CanPayWithPigeon:              true,
			CanPayWithBankTransfer:        true,
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
		neo4jutil.NodeLabelTenant:               1,
		neo4jutil.NodeLabelTenantBillingProfile: 1,
	})

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jutil.NodeLabelTenantBillingProfile, profileId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	// verify
	tenantBillingProfileEntity := neo4jmapper.MapDbNodeToTenantBillingProfileEntity(dbNode)
	require.Equal(t, profileId, tenantBillingProfileEntity.Id)
	require.Equal(t, timeNow, tenantBillingProfileEntity.UpdatedAt)
	require.Equal(t, "phone", tenantBillingProfileEntity.Phone)
	require.Equal(t, "addressLine1", tenantBillingProfileEntity.AddressLine1)
	require.Equal(t, "addressLine2", tenantBillingProfileEntity.AddressLine2)
	require.Equal(t, "addressLine3", tenantBillingProfileEntity.AddressLine3)
	require.Equal(t, "locality", tenantBillingProfileEntity.Locality)
	require.Equal(t, "country", tenantBillingProfileEntity.Country)
	require.Equal(t, "zip", tenantBillingProfileEntity.Zip)
	require.Equal(t, "legalName", tenantBillingProfileEntity.LegalName)
	require.Equal(t, "domesticPaymentsBankInfo", tenantBillingProfileEntity.DomesticPaymentsBankInfo)
	require.Equal(t, "internationalPaymentsBankInfo", tenantBillingProfileEntity.InternationalPaymentsBankInfo)
	require.Equal(t, "vatNumber", tenantBillingProfileEntity.VatNumber)
	require.Equal(t, "sendInvoicesFrom", tenantBillingProfileEntity.SendInvoicesFrom)
	require.Equal(t, "sendInvoicesBcc", tenantBillingProfileEntity.SendInvoicesBcc)
	require.Equal(t, true, tenantBillingProfileEntity.CanPayWithCard)
	require.Equal(t, true, tenantBillingProfileEntity.CanPayWithDirectDebitSEPA)
	require.Equal(t, true, tenantBillingProfileEntity.CanPayWithDirectDebitACH)
	require.Equal(t, true, tenantBillingProfileEntity.CanPayWithDirectDebitBacs)
	require.Equal(t, true, tenantBillingProfileEntity.CanPayWithPigeon)
	require.Equal(t, true, tenantBillingProfileEntity.CanPayWithBankTransfer)
}

func TestTenantEventHandler_OnUpdateTenantSettingsV1(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	settingsId := neo4jtest.CreateTenantSettings(ctx, testDatabase.Driver, tenantName, neo4jentity.TenantSettingsEntity{})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		neo4jutil.NodeLabelTenant:         1,
		neo4jutil.NodeLabelTenantSettings: 1,
	})

	// Prepare the event handler
	eventHandler := &TenantEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	timeNow := utils.Now()

	aggregate := tenant.NewTenantAggregate(tenantName)
	updateEvent, err := event.NewTenantSettingsUpdateEvent(
		aggregate,
		&tenantpb.UpdateTenantSettingsRequest{
			LogoRepositoryFileId: "logoRepositoryFileId",
			BaseCurrency:         neo4jenum.CurrencyAUD.String(),
			InvoicingEnabled:     true,
			InvoicingPostpaid:    true,
		},
		timeNow,
		[]string{},
	)
	require.Nil(t, err)

	// EXECUTE
	err = eventHandler.OnUpdateTenantSettingsV1(context.Background(), updateEvent)
	require.Nil(t, err)

	// check still same nodes available
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		neo4jutil.NodeLabelTenant:         1,
		neo4jutil.NodeLabelTenantSettings: 1,
	})

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jutil.NodeLabelTenantSettings, settingsId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	// verify
	tenantSettingsEntity := neo4jmapper.MapDbNodeToTenantSettingsEntity(dbNode)
	require.Equal(t, settingsId, tenantSettingsEntity.Id)
	require.Equal(t, timeNow, tenantSettingsEntity.UpdatedAt)
	require.Equal(t, "logoRepositoryFileId", tenantSettingsEntity.LogoRepositoryFileId)
	require.Equal(t, true, tenantSettingsEntity.InvoicingEnabled)
	require.Equal(t, true, tenantSettingsEntity.InvoicingPostpaid)
	require.Equal(t, neo4jenum.CurrencyAUD, tenantSettingsEntity.BaseCurrency)
}
