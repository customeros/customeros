package graph

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoice"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/mocked_grpc"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
)

func TestInvoiceEventHandler_OnInvoiceCreateForContractV1(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContract(ctx, testDatabase.Driver, tenantName, organizationId, neo4jentity.ContractEntity{})

	eventHandler := &InvoiceEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	now := utils.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)
	invoiceId := uuid.New().String()

	aggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenantName, invoiceId)
	newEvent, err := invoice.NewInvoiceForContractCreateEvent(
		aggregate,
		commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		},
		contractId,
		"EUR",
		"INV-123",
		neo4jenum.BillingCycleMonthlyBilling.String(),
		true,
		now,
		yesterday,
		tomorrow,
	)
	require.Nil(t, err)

	// EXECUTE
	err = eventHandler.OnInvoiceCreateForContractV1(context.Background(), newEvent)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		neo4jutil.NodeLabelInvoice:                    1,
		neo4jutil.NodeLabelInvoice + "_" + tenantName: 1})

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jutil.NodeLabelInvoice, invoiceId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	// verify
	createdInvoice := neo4jmapper.MapDbNodeToInvoiceEntity(dbNode)

	require.Equal(t, invoiceId, createdInvoice.Id)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), createdInvoice.Source)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), createdInvoice.SourceOfTruth)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, createdInvoice.AppSource)
	require.Equal(t, now, createdInvoice.CreatedAt)
	require.Equal(t, now, createdInvoice.UpdatedAt)
	require.Equal(t, true, createdInvoice.DryRun)
	require.Equal(t, "INV-123", createdInvoice.Number)
	require.Equal(t, yesterday, createdInvoice.PeriodStartDate)
	require.Equal(t, tomorrow, createdInvoice.PeriodEndDate)
	require.Equal(t, neo4jenum.BillingCycleMonthlyBilling, createdInvoice.BillingCycle)
	require.Equal(t, float64(0), createdInvoice.Amount)
	require.Equal(t, float64(0), createdInvoice.Vat)
	require.Equal(t, float64(0), createdInvoice.Amount)
	require.Equal(t, neo4jenum.CurrencyEUR, createdInvoice.Currency)
	require.Equal(t, "", createdInvoice.RepositoryFileId)
}

func TestInvoiceEventHandler_OnInvoiceFillV1(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContract(ctx, testDatabase.Driver, tenantName, organizationId, neo4jentity.ContractEntity{})
	invoiceId := neo4jtest.CreateInvoiceForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.InvoiceEntity{
		Currency: neo4jenum.CurrencyEUR,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		neo4jutil.NodeLabelInvoice:      1,
		neo4jutil.NodeLabelOrganization: 1,
		neo4jutil.NodeLabelContract:     1,
		neo4jutil.NodeLabelInvoiceLine:  0,
	})

	// prepare grpc mock
	calledGenerateInvoicePdfGrpcRequest := false
	invoiceGrpcServiceCallbacks := mocked_grpc.MockInvoiceServiceCallbacks{
		GenerateInvoicePdf: func(context context.Context, inv *invoicepb.GenerateInvoicePdfRequest) (*invoicepb.InvoiceIdResponse, error) {
			require.Equal(t, tenantName, inv.Tenant)
			require.Equal(t, invoiceId, inv.InvoiceId)
			require.Equal(t, "", inv.LoggedInUserId)
			require.Equal(t, constants.AppSourceEventProcessingPlatform, inv.AppSource)
			calledGenerateInvoicePdfGrpcRequest = true
			return &invoicepb.InvoiceIdResponse{
				Id: invoiceId,
			}, nil
		},
	}
	mocked_grpc.SetInvoiceCallbacks(&invoiceGrpcServiceCallbacks)

	// Prepare the event handler
	eventHandler := &InvoiceEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
		grpcClients:  testMockedGrpcClient,
	}

	timeNow := utils.Now()

	aggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenantName, invoiceId)
	fillEvent, err := invoice.NewInvoiceFillEvent(
		aggregate,
		timeNow,
		invoice.Invoice{
			ContractId: contractId,
		},
		100,
		20,
		120,
		[]invoice.InvoiceLineEvent{
			{
				Id:        "invoice-line-id-1",
				CreatedAt: timeNow,
				SourceFields: commonmodel.Source{
					Source:    constants.SourceOpenline,
					AppSource: constants.AppSourceEventProcessingPlatform,
				},
				Name:                    "test",
				Price:                   50,
				Quantity:                2,
				Amount:                  100,
				VAT:                     20,
				TotalAmount:             120,
				ServiceLineItemId:       "service-line-item-id",
				ServiceLineItemParentId: "service-line-item-parent-id",
				BilledType:              neo4jenum.BilledTypeMonthly.String(),
			},
			{
				Id:        "invoice-line-id-2",
				CreatedAt: timeNow,
				SourceFields: commonmodel.Source{
					Source:    constants.SourceOpenline,
					AppSource: constants.AppSourceEventProcessingPlatform,
				},
				Name:                    "test02",
				Price:                   50.2,
				Quantity:                22,
				Amount:                  100.2,
				VAT:                     20.2,
				TotalAmount:             120.2,
				ServiceLineItemId:       "service-line-item-id-2",
				ServiceLineItemParentId: "service-line-item-parent-id-2",
				BilledType:              neo4jenum.BilledTypeAnnually.String(),
			},
		},
	)
	require.Nil(t, err)

	// EXECUTE
	err = eventHandler.OnInvoiceFillV1(context.Background(), fillEvent)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		neo4jutil.NodeLabelInvoice:                        1,
		neo4jutil.NodeLabelInvoice + "_" + tenantName:     1,
		neo4jutil.NodeLabelInvoiceLine:                    2,
		neo4jutil.NodeLabelInvoiceLine + "_" + tenantName: 2,
	})

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jutil.NodeLabelInvoice, invoiceId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	// verify
	invoiceEntity := neo4jmapper.MapDbNodeToInvoiceEntity(dbNode)
	require.Equal(t, invoiceId, invoiceEntity.Id)
	require.Equal(t, timeNow, invoiceEntity.UpdatedAt)
	require.Equal(t, float64(100), invoiceEntity.Amount)
	require.Equal(t, float64(20), invoiceEntity.Vat)
	require.Equal(t, float64(120), invoiceEntity.TotalAmount)

	// verify invoice lines
	dbNodes, err := neo4jtest.GetAllNodesByLabel(ctx, testDatabase.Driver, neo4jutil.NodeLabelInvoiceLine)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	invoiceLineEntity1 := neo4jmapper.MapDbNodeToInvoiceLineEntity(dbNodes[0])
	invoiceLineEntity2 := neo4jmapper.MapDbNodeToInvoiceLineEntity(dbNodes[1])
	var firstInvoiceLine, secondInvoiceLine neo4jentity.InvoiceLineEntity
	if invoiceLineEntity1.Id == "invoice-line-id-1" {
		firstInvoiceLine = *invoiceLineEntity1
		secondInvoiceLine = *invoiceLineEntity2
	} else {
		firstInvoiceLine = *invoiceLineEntity2
		secondInvoiceLine = *invoiceLineEntity1
	}

	require.Equal(t, "invoice-line-id-1", firstInvoiceLine.Id)
	require.Equal(t, timeNow, firstInvoiceLine.CreatedAt)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), firstInvoiceLine.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, firstInvoiceLine.AppSource)
	require.Equal(t, "test", firstInvoiceLine.Name)
	require.Equal(t, float64(50), firstInvoiceLine.Price)
	require.Equal(t, int64(2), firstInvoiceLine.Quantity)
	require.Equal(t, float64(100), firstInvoiceLine.Amount)
	require.Equal(t, float64(20), firstInvoiceLine.Vat)
	require.Equal(t, float64(120), firstInvoiceLine.TotalAmount)
	require.Equal(t, "service-line-item-id", firstInvoiceLine.ServiceLineItemId)
	require.Equal(t, "service-line-item-parent-id", firstInvoiceLine.ServiceLineItemParentId)
	require.Equal(t, neo4jenum.BilledTypeMonthly, firstInvoiceLine.BilledType)

	require.Equal(t, "invoice-line-id-2", secondInvoiceLine.Id)
	require.Equal(t, timeNow, secondInvoiceLine.CreatedAt)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), secondInvoiceLine.Source)
	require.Equal(t, constants.AppSourceEventProcessingPlatform, secondInvoiceLine.AppSource)
	require.Equal(t, "test02", secondInvoiceLine.Name)
	require.Equal(t, float64(50.2), secondInvoiceLine.Price)
	require.Equal(t, int64(22), secondInvoiceLine.Quantity)
	require.Equal(t, float64(100.2), secondInvoiceLine.Amount)
	require.Equal(t, float64(20.2), secondInvoiceLine.Vat)
	require.Equal(t, float64(120.2), secondInvoiceLine.TotalAmount)
	require.Equal(t, "service-line-item-id-2", secondInvoiceLine.ServiceLineItemId)
	require.Equal(t, "service-line-item-parent-id-2", secondInvoiceLine.ServiceLineItemParentId)
	require.Equal(t, neo4jenum.BilledTypeAnnually, secondInvoiceLine.BilledType)

	require.True(t, calledGenerateInvoicePdfGrpcRequest)
}

func TestInvoiceEventHandler_OnInvoicePdfGenerated(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// prepare neo4j data
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	organizationId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContract(ctx, testDatabase.Driver, tenantName, organizationId, neo4jentity.ContractEntity{})
	id := neo4jtest.CreateInvoiceForContract(ctx, testDatabase.Driver, tenantName, contractId, neo4jentity.InvoiceEntity{})

	// Prepare the event handler
	eventHandler := &InvoiceEventHandler{
		log:          testLogger,
		repositories: testDatabase.Repositories,
	}

	timeNow := utils.Now()

	aggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenantName, id)
	pdfGeneratedEvent, err := invoice.NewInvoicePdfGeneratedEvent(
		aggregate,
		timeNow,
		"test",
	)
	require.Nil(t, err)

	// EXECUTE
	err = eventHandler.OnInvoicePdfGenerated(context.Background(), pdfGeneratedEvent)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		neo4jutil.NodeLabelInvoice:                    1,
		neo4jutil.NodeLabelInvoice + "_" + tenantName: 1,
	})

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, neo4jutil.NodeLabelInvoice, id)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	// verify
	invoice := neo4jmapper.MapDbNodeToInvoiceEntity(dbNode)
	require.Equal(t, id, invoice.Id)
	require.Equal(t, timeNow, invoice.UpdatedAt)

	require.Equal(t, "test", invoice.RepositoryFileId)
}
