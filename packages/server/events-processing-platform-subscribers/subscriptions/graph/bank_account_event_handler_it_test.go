package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	tenant "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/tenant"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/tenant/event"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBankAccountEventHandler_OnAddBankAccountV1(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		neo4jutil.NodeLabelTenant:      1,
		neo4jutil.NodeLabelBankAccount: 0,
	})

	// Prepare the event handler
	eventHandler := &BankAccountEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	now := utils.Now()
	bankAccountId := uuid.New().String()

	aggregate := tenant.NewTenantAggregate(tenantName)
	updateEvent, err := event.NewTenantBankAccountCreateEvent(
		aggregate,
		commonmodel.Source{
			Source:    "openline",
			AppSource: "appSource",
		},
		bankAccountId,
		&tenantpb.AddBankAccountGrpcRequest{
			BankName:            "bankName",
			BankTransferEnabled: true,
			AllowInternational:  true,
			Currency:            neo4jenum.CurrencyUSD.String(),
			AccountNumber:       "accountNumber",
			SortCode:            "sortCode",
			Iban:                "iban",
			Bic:                 "bic",
			RoutingNumber:       "routingNumber",
			OtherDetails:        "otherDetails",
		},
		now,
	)
	require.Nil(t, err)

	// EXECUTE
	err = eventHandler.OnAddBankAccountV1(context.Background(), updateEvent)
	require.Nil(t, err)

	// check still same nodes available
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		neo4jutil.NodeLabelTenant:      1,
		neo4jutil.NodeLabelBankAccount: 1,
	})

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jutil.NodeLabelBankAccount, bankAccountId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	// verify
	bankAccountEntity := neo4jmapper.MapDbNodeToBankAccountEntity(dbNode)
	require.Equal(t, bankAccountId, bankAccountEntity.Id)
	require.Equal(t, "bankName", bankAccountEntity.BankName)
	require.Equal(t, true, bankAccountEntity.BankTransferEnabled)
	require.Equal(t, true, bankAccountEntity.AllowInternational)
	require.Equal(t, neo4jenum.CurrencyUSD, bankAccountEntity.Currency)
	require.Equal(t, "accountNumber", bankAccountEntity.AccountNumber)
	require.Equal(t, "sortCode", bankAccountEntity.SortCode)
	require.Equal(t, "iban", bankAccountEntity.Iban)
	require.Equal(t, "bic", bankAccountEntity.Bic)
	require.Equal(t, "routingNumber", bankAccountEntity.RoutingNumber)
	require.Equal(t, "otherDetails", bankAccountEntity.OtherDetails)
	require.Equal(t, now, bankAccountEntity.CreatedAt)
	test.AssertRecentTime(t, bankAccountEntity.UpdatedAt)
	require.Equal(t, neo4jentity.DataSourceOpenline, bankAccountEntity.Source)
	require.Equal(t, neo4jentity.DataSourceOpenline, bankAccountEntity.SourceOfTruth)
	require.Equal(t, "appSource", bankAccountEntity.AppSource)
}

func TestBankAccountEventHandler_OnUpdateBankAccountV1(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	bankAccountId := neo4jtest.CreateBankAccount(ctx, testDatabase.Driver, tenantName, neo4jentity.BankAccountEntity{
		BankName:            "bankName",
		BankTransferEnabled: false,
		AllowInternational:  false,
		Currency:            neo4jenum.CurrencyEUR,
		AccountNumber:       "accountNumber",
		SortCode:            "sortCode",
		Iban:                "iban",
		Bic:                 "bic",
		RoutingNumber:       "routingNumber",
		OtherDetails:        "otherDetails",
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		neo4jutil.NodeLabelTenant:      1,
		neo4jutil.NodeLabelBankAccount: 1,
	})

	// Prepare the event handler
	eventHandler := &BankAccountEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	now := utils.Now()

	aggregate := tenant.NewTenantAggregate(tenantName)
	updateEvent, err := event.NewTenantBankAccountUpdateEvent(
		aggregate,
		bankAccountId,
		&tenantpb.UpdateBankAccountGrpcRequest{
			BankName:            "updatedName",
			BankTransferEnabled: true,
			AllowInternational:  true,
			Currency:            neo4jenum.CurrencyUSD.String(),
			AccountNumber:       "updatedAccountNumber",
			SortCode:            "updatedSortCode",
			Iban:                "updatedIban",
			Bic:                 "updatedBic",
			RoutingNumber:       "updatedRoutingNumber",
			OtherDetails:        "updatedOtherDetails",
		},
		now,
		[]string{
			event.FieldMaskBankAccountBankName,
			event.FieldMaskBankAccountBankTransferEnabled,
			event.FieldMaskBankAccountAllowInternational,
			event.FieldMaskBankAccountCurrency,
			event.FieldMaskBankAccountAccountNumber,
			event.FieldMaskBankAccountSortCode,
			event.FieldMaskBankAccountIban,
			event.FieldMaskBankAccountBic,
			event.FieldMaskBankAccountRoutingNumber,
			event.FieldMaskBankAccountOtherDetails,
		},
	)
	require.Nil(t, err)

	// EXECUTE
	err = eventHandler.OnUpdateBankAccountV1(context.Background(), updateEvent)
	require.Nil(t, err)

	// check still same nodes available
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		neo4jutil.NodeLabelTenant:      1,
		neo4jutil.NodeLabelBankAccount: 1,
	})

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jutil.NodeLabelBankAccount, bankAccountId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	// verify
	bankAccountEntity := neo4jmapper.MapDbNodeToBankAccountEntity(dbNode)
	require.Equal(t, bankAccountId, bankAccountEntity.Id)
	require.Equal(t, "updatedName", bankAccountEntity.BankName)
	require.Equal(t, true, bankAccountEntity.BankTransferEnabled)
	require.Equal(t, true, bankAccountEntity.AllowInternational)
	require.Equal(t, neo4jenum.CurrencyUSD, bankAccountEntity.Currency)
	require.Equal(t, "updatedAccountNumber", bankAccountEntity.AccountNumber)
	require.Equal(t, "updatedSortCode", bankAccountEntity.SortCode)
	require.Equal(t, "updatedIban", bankAccountEntity.Iban)
	require.Equal(t, "updatedBic", bankAccountEntity.Bic)
	require.Equal(t, "updatedRoutingNumber", bankAccountEntity.RoutingNumber)
	require.Equal(t, "updatedOtherDetails", bankAccountEntity.OtherDetails)
	test.AssertRecentTime(t, bankAccountEntity.UpdatedAt)
}

func TestBankAccountEventHandler_OnDeleteBankAccountV1(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	bankAccountId := neo4jtest.CreateBankAccount(ctx, testDatabase.Driver, tenantName, neo4jentity.BankAccountEntity{
		BankName: "bankName",
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		neo4jutil.NodeLabelTenant:      1,
		neo4jutil.NodeLabelBankAccount: 1,
	})

	// Prepare the event handler
	eventHandler := &BankAccountEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	now := utils.Now()

	aggregate := tenant.NewTenantAggregate(tenantName)
	deleteEvent, err := event.NewTenantBankAccountDeleteEvent(aggregate, bankAccountId, now)
	require.Nil(t, err)

	// EXECUTE
	err = eventHandler.OnDeleteBankAccountV1(context.Background(), deleteEvent)
	require.Nil(t, err)

	// check still same nodes available
	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		neo4jutil.NodeLabelTenant:      1,
		neo4jutil.NodeLabelBankAccount: 0,
	})
}
