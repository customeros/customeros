package servicet

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	contractaggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoice"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func TestInvoiceService_NewOnCycleInvoiceForContract(t *testing.T) {
	ctx := context.Background()

	// Setup test environment
	tenant := "tenant1"
	contractId := uuid.New().String()

	aggregateStore := eventstoret.NewTestAggregateStore()
	contractAggregate := contractaggregate.NewContractAggregateWithTenantAndID(tenant, contractId)
	aggregateStore.Save(ctx, contractAggregate)

	// Mock GRPC Connection and Invoice Client
	grpcConnection, _ := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	invoiceClient := invoicepb.NewInvoiceGrpcServiceClient(grpcConnection)

	// Prepare request
	timeNow := utils.Now()
	yesterday := timeNow.AddDate(0, 0, -1)
	nextWeek := timeNow.AddDate(0, 0, 7)
	createRequest := &invoicepb.NewInvoiceForContractRequest{
		Tenant:         tenant,
		LoggedInUserId: "User123",
		ContractId:     contractId,
		SourceFields: &commonpb.SourceFields{
			AppSource: "unit-test",
		},
		CreatedAt:          timestamppb.New(timeNow),
		Currency:           "USD",
		InvoicePeriodStart: utils.ConvertTimeToTimestampPtr(&yesterday),
		InvoicePeriodEnd:   utils.ConvertTimeToTimestampPtr(&nextWeek),
		DryRun:             false,
	}

	// Execute service method
	response, err := invoiceClient.NewInvoiceForContract(ctx, createRequest)

	// Validate response and error
	require.Nil(t, err, "Failed to create invoice")
	require.NotNil(t, response)
	invoiceId := response.Id
	require.NotEmpty(t, invoiceId)

	// Validate created invoice in event store
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 2, len(eventsMap))

	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenant, invoiceId)
	eventList := eventsMap[invoiceAggregate.GetID()]
	require.Equal(t, 1, len(eventList))

	// Validate event data
	createEvent := eventList[0]
	require.Equal(t, invoice.InvoiceCreateForContractV1, createEvent.GetEventType())
	require.Equal(t, string(invoice.InvoiceAggregateType)+"-"+tenant+"-"+invoiceId, eventList[0].GetAggregateID())

	var eventData invoice.InvoiceForContractCreateEvent
	err = createEvent.GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")
	require.Equal(t, contractId, eventData.ContractId)
	require.Equal(t, timeNow, eventData.CreatedAt)
	require.Equal(t, "unit-test", eventData.SourceFields.AppSource)
	require.Equal(t, false, eventData.DryRun)
	require.NotEmpty(t, eventData.InvoiceNumber)
	require.Equal(t, "USD", eventData.Currency)
	require.Equal(t, yesterday, eventData.PeriodStartDate)
	require.Equal(t, nextWeek, eventData.PeriodEndDate)
}

func TestInvoiceService_FillInvoice(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// setup test environment
	tenant := "ziggy"
	invoiceId := "invoice-id"
	now := utils.Now()

	aggregateStore := eventstoret.NewTestAggregateStore()
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenant, invoiceId)

	newEvent, _ := invoice.NewInvoiceForContractCreateEvent(invoiceAggregate, commonmodel.Source{}, "contract-1", "USD", "1", "MONTHLY_BILLED", "test note", false, now, now, now)
	invoiceAggregate.UncommittedEvents = append(invoiceAggregate.UncommittedEvents, newEvent)
	aggregateStore.Save(ctx, invoiceAggregate)

	// prepare connection to grpc server
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	invoiceClient := invoicepb.NewInvoiceGrpcServiceClient(grpcConnection)

	// Execute the command
	response, err := invoiceClient.FillInvoice(ctx, &invoicepb.FillInvoiceRequest{
		Tenant:         tenant,
		LoggedInUserId: "user-id",
		InvoiceId:      invoiceId,
		Amount:         1.01,
		Vat:            2.02,
		Total:          3.03,
		InvoiceLines: []*invoicepb.InvoiceLine{
			{
				Name:                    "name",
				Price:                   2.02,
				Quantity:                3,
				Amount:                  4.04,
				Vat:                     5.05,
				Total:                   6.06,
				ServiceLineItemId:       "service-line-item-id",
				ServiceLineItemParentId: "service-line-item-parent-id",
				BilledType:              commonpb.BilledType_MONTHLY_BILLED,
			},
		},
		UpdatedAt: utils.ConvertTimeToTimestampPtr(&now),
		AppSource: "test",
	})
	require.Nil(t, err)
	require.NotNil(t, response)

	// verify
	require.Equal(t, invoiceId, response.Id)

	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))

	eventList := eventsMap[invoiceAggregate.ID]
	require.Equal(t, 2, len(eventList))

	require.Equal(t, invoice.InvoiceCreateForContractV1, eventList[0].GetEventType())
	require.Equal(t, invoice.InvoiceFillV1, eventList[1].GetEventType())
	require.Equal(t, string(invoice.InvoiceAggregateType)+"-"+tenant+"-"+invoiceId, eventList[1].GetAggregateID())

	var eventData invoice.InvoiceFillEvent
	err = eventList[1].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	require.Equal(t, float64(1.01), eventData.Amount)
	require.Equal(t, float64(2.02), eventData.VAT)
	require.Equal(t, float64(3.03), eventData.TotalAmount)
	require.Equal(t, "contract-1", eventData.ContractId)
	require.Equal(t, "USD", eventData.Currency)
	require.Equal(t, 1, len(eventData.InvoiceLines))
	require.Equal(t, "name", eventData.InvoiceLines[0].Name)
	require.Equal(t, float64(2.02), eventData.InvoiceLines[0].Price)
	require.Equal(t, int64(3), eventData.InvoiceLines[0].Quantity)
	require.Equal(t, float64(4.04), eventData.InvoiceLines[0].Amount)
	require.Equal(t, float64(5.05), eventData.InvoiceLines[0].VAT)
	require.Equal(t, float64(6.06), eventData.InvoiceLines[0].TotalAmount)
	require.Equal(t, "service-line-item-id", eventData.InvoiceLines[0].ServiceLineItemId)
	require.Equal(t, "service-line-item-parent-id", eventData.InvoiceLines[0].ServiceLineItemParentId)
	require.Equal(t, neo4jenum.BilledTypeMonthly.String(), eventData.InvoiceLines[0].BilledType)
}

func TestInvoiceService_GenerateInvoicePdf(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// setup test environment
	tenant := "ziggy"
	invoiceId := "invoice-id"
	now := utils.Now()

	aggregateStore := eventstoret.NewTestAggregateStore()
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenant, invoiceId)

	newEvent, _ := invoice.NewInvoiceForContractCreateEvent(invoiceAggregate, commonmodel.Source{}, "contract-1", "USD", "1", "MONTHLY_BILLED", "test note", false, now, now, now)
	invoiceAggregate.UncommittedEvents = append(invoiceAggregate.UncommittedEvents, newEvent)
	aggregateStore.Save(ctx, invoiceAggregate)

	// prepare connection to grpc server
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	invoiceClient := invoicepb.NewInvoiceGrpcServiceClient(grpcConnection)

	// Execute the command
	response, err := invoiceClient.GenerateInvoicePdf(ctx, &invoicepb.GenerateInvoicePdfRequest{
		Tenant:         tenant,
		LoggedInUserId: "user-id",
		InvoiceId:      invoiceId,
		AppSource:      "test",
	})
	require.Nil(t, err)
	require.NotNil(t, response)

	// verify
	require.Equal(t, invoiceId, response.Id)

	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))

	eventList := eventsMap[invoiceAggregate.ID]
	require.Equal(t, 2, len(eventList))

	require.Equal(t, invoice.InvoiceCreateForContractV1, eventList[0].GetEventType())
	require.Equal(t, invoice.InvoicePdfRequestedV1, eventList[1].GetEventType())
	require.Equal(t, string(invoice.InvoiceAggregateType)+"-"+tenant+"-"+invoiceId, eventList[1].GetAggregateID())

	var eventData invoice.InvoicePdfRequestedEvent
	err = eventList[1].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	require.Equal(t, tenant, eventData.Tenant)
}

func TestInvoiceService_UpdateInvoice(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// setup test environment
	tenant := "ziggy"
	invoiceId := "invoice-id"
	now := utils.Now()

	aggregateStore := eventstoret.NewTestAggregateStore()
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenant, invoiceId)

	newEvent, _ := invoice.NewInvoiceForContractCreateEvent(invoiceAggregate, commonmodel.Source{}, "contract-1", "USD", "1", "MONTHLY_BILLED", "test note", false, now, now, now)
	invoiceAggregate.UncommittedEvents = append(invoiceAggregate.UncommittedEvents, newEvent)
	aggregateStore.Save(ctx, invoiceAggregate)

	// prepare connection to grpc server
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	invoiceClient := invoicepb.NewInvoiceGrpcServiceClient(grpcConnection)

	// Execute the command
	response, err := invoiceClient.UpdateInvoice(ctx, &invoicepb.UpdateInvoiceRequest{
		Tenant:         tenant,
		LoggedInUserId: "user-id",
		InvoiceId:      invoiceId,
		AppSource:      "test",
		Status:         invoicepb.InvoiceStatus_INVOICE_STATUS_PAID,
		FieldsMask:     []invoicepb.InvoiceFieldMask{invoicepb.InvoiceFieldMask_INVOICE_FIELD_STATUS},
	})
	require.Nil(t, err)
	require.NotNil(t, response)

	// verify
	require.Equal(t, invoiceId, response.Id)

	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))

	eventList := eventsMap[invoiceAggregate.ID]
	require.Equal(t, 2, len(eventList))

	require.Equal(t, invoice.InvoiceCreateForContractV1, eventList[0].GetEventType())
	require.Equal(t, invoice.InvoiceUpdateV1, eventList[1].GetEventType())
	require.Equal(t, string(invoice.InvoiceAggregateType)+"-"+tenant+"-"+invoiceId, eventList[1].GetAggregateID())

	var eventData invoice.InvoiceUpdateEvent
	err = eventList[1].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, neo4jenum.InvoiceStatusPaid.String(), eventData.Status)
	require.Equal(t, 1, len(eventData.FieldsMask))
	require.Equal(t, []string{invoice.FieldMaskStatus}, eventData.FieldsMask)
}

//
//func TestInvoiceService_PdfGeneratedInvoice(t *testing.T) {
//	ctx := context.TODO()
//	defer tearDownTestCase(ctx, testDatabase)(t)
//
//	// setup test environment
//	tenant := "ziggy"
//	invoiceId := "invoice-id"
//	now := utils.Now()
//
//	aggregateStore := eventstoret.NewTestAggregateStore()
//	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenant, invoiceId)
//
//	newEvent, _ := invoice.NewInvoiceCreateEvent(invoiceAggregate, commonmodel.Source{}, &invoicepb.NewInvoiceForContractRequest{
//		ContractId:         "1",
//		CreatedAt:          utils.ConvertTimeToTimestampPtr(&now),
//		InvoicePeriodStart: utils.ConvertTimeToTimestampPtr(&now),
//	})
//	fillEvent, _ := invoice.NewInvoiceFillEvent(invoiceAggregate, &now, commonmodel.Source{}, &invoicepb.FillInvoiceRequest{
//		Amount: 1,
//		Vat:    2,
//		Total:  3,
//		Lines:  []*invoicepb.InvoiceLine{},
//	})
//	invoiceAggregate.UncommittedEvents = append(invoiceAggregate.UncommittedEvents, newEvent)
//	invoiceAggregate.UncommittedEvents = append(invoiceAggregate.UncommittedEvents, fillEvent)
//	aggregateStore.Save(ctx, invoiceAggregate)
//
//	// prepare connection to grpc server
//	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
//	require.Nil(t, err)
//	invoiceClient := invoicepb.NewInvoiceGrpcServiceClient(grpcConnection)
//
//	// Execute the command
//	response, err := invoiceClient.PdfGeneratedInvoice(ctx, &invoicepb.PdfGeneratedInvoiceRequest{
//		Tenant:           tenant,
//		LoggedInUserId:   "user-id",
//		InvoiceId:        invoiceId,
//		RepositoryFileId: "repository-file-id",
//		UpdatedAt:        utils.ConvertTimeToTimestampPtr(&now),
//		AppSource:        "app",
//	})
//	require.Nil(t, err)
//	require.NotNil(t, response)
//
//	// verify
//	require.Equal(t, invoiceId, response.Id)
//
//	eventsMap := aggregateStore.GetEventMap()
//	require.Equal(t, 1, len(eventsMap))
//
//	eventList := eventsMap[invoiceAggregate.ID]
//	require.Equal(t, 3, len(eventList))
//
//	require.Equal(t, invoice.InvoiceCreateForContractV1, eventList[0].GetEventType())
//	require.Equal(t, invoice.InvoiceFillV1, eventList[1].GetEventType())
//	require.Equal(t, invoice.InvoicePdfGeneratedV1, eventList[2].GetEventType())
//	require.Equal(t, string(invoice.InvoiceAggregateType)+"-"+tenant+"-"+invoiceId, eventList[1].GetAggregateID())
//
//	var eventData invoice.InvoicePdfGeneratedEvent
//	err = eventList[2].GetJsonData(&eventData)
//	require.Nil(t, err, "Failed to unmarshal event data")
//
//	require.Equal(t, "repository-file-id", eventData.RepositoryFileId)
//}
//
//func TestInvoiceService_PayInvoice(t *testing.T) {
//	ctx := context.TODO()
//	defer tearDownTestCase(ctx, testDatabase)(t)
//
//	// setup test environment
//	tenant := "ziggy"
//	invoiceId := "invoice-id"
//	now := utils.Now()
//
//	aggregateStore := eventstoret.NewTestAggregateStore()
//	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenant, invoiceId)
//
//	newEvent, _ := invoice.NewInvoiceCreateEvent(invoiceAggregate, commonmodel.Source{}, &invoicepb.NewInvoiceForContractRequest{
//		ContractId:         "1",
//		CreatedAt:          utils.ConvertTimeToTimestampPtr(&now),
//		InvoicePeriodStart: utils.ConvertTimeToTimestampPtr(&now),
//	})
//	fillEvent, _ := invoice.NewInvoiceFillEvent(invoiceAggregate, &now, commonmodel.Source{}, &invoicepb.FillInvoiceRequest{
//		Amount: 1,
//		Vat:    2,
//		Total:  3,
//		Lines:  []*invoicepb.InvoiceLine{},
//	})
//	pdfGeneratedEvent, _ := invoice.NewInvoicePdfGeneratedEvent(invoiceAggregate, now, "repository-file-id")
//	invoiceAggregate.UncommittedEvents = append(invoiceAggregate.UncommittedEvents, newEvent)
//	invoiceAggregate.UncommittedEvents = append(invoiceAggregate.UncommittedEvents, fillEvent)
//	invoiceAggregate.UncommittedEvents = append(invoiceAggregate.UncommittedEvents, pdfGeneratedEvent)
//	aggregateStore.Save(ctx, invoiceAggregate)
//
//	// prepare connection to grpc server
//	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
//	require.Nil(t, err)
//	invoiceClient := invoicepb.NewInvoiceGrpcServiceClient(grpcConnection)
//
//	// Execute the command
//	response, err := invoiceClient.PayInvoice(ctx, &invoicepb.PayInvoiceRequest{
//		Tenant:         tenant,
//		LoggedInUserId: "user-id",
//		InvoiceId:      invoiceId,
//		UpdatedAt:      utils.ConvertTimeToTimestampPtr(&now),
//		SourceFields: &commonpb.SourceFields{
//			AppSource: "app",
//		},
//	})
//	require.Nil(t, err)
//	require.NotNil(t, response)
//
//	// verify
//	require.Equal(t, invoiceId, response.Id)
//
//	eventsMap := aggregateStore.GetEventMap()
//	require.Equal(t, 1, len(eventsMap))
//
//	eventList := eventsMap[invoiceAggregate.ID]
//	require.Equal(t, 4, len(eventList))
//
//	require.Equal(t, invoice.InvoiceCreateForContractV1, eventList[0].GetEventType())
//	require.Equal(t, invoice.InvoiceFillV1, eventList[1].GetEventType())
//	require.Equal(t, invoice.InvoicePdfGeneratedV1, eventList[2].GetEventType())
//	require.Equal(t, invoice.InvoicePayV1, eventList[3].GetEventType())
//	require.Equal(t, string(invoice.InvoiceAggregateType)+"-"+tenant+"-"+invoiceId, eventList[0].GetAggregateID())
//
//	var eventData invoice.InvoicePayEvent
//	err = eventList[3].GetJsonData(&eventData)
//	require.Nil(t, err, "Failed to unmarshal event data")
//
//	test.AssertRecentTime(t, eventData.UpdatedAt)
//}
//
//func TestInvoiceService_SimulateInvoice(t *testing.T) {
//	ctx := context.TODO()
//	defer tearDownTestCase(ctx, testDatabase)(t)
//
//	// setup test environment
//	tenant := "ziggy"
//	now := utils.Now()
//
//	aggregateStore := eventstoret.NewTestAggregateStore()
//	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
//	require.Nil(t, err, "Failed to get grpc connection")
//	invoiceServiceClient := invoicepb.NewInvoiceGrpcServiceClient(grpcConnection)
//
//	response, err := invoiceServiceClient.SimulateInvoice(ctx, &invoicepb.SimulateInvoiceRequest{
//		Tenant:     tenant,
//		ContractId: "1",
//		SourceFields: &commonpb.SourceFields{
//			AppSource: "app",
//			Source:    "source",
//		},
//		CreatedAt: utils.ConvertTimeToTimestampPtr(&now),
//		Date:      utils.ConvertTimeToTimestampPtr(&now),
//		DryRunInvoices: []*invoicepb.DryRunInvoice{
//			{
//				InvoiceId: "1",
//				Name:              "name 1",
//				Billed:            commonpb.BilledType_MONTHLY_BILLED,
//				Price:             1,
//				Quantity:          2,
//			},
//			{
//				Name:     "name 2",
//				Price:    1,
//				Quantity: 2,
//			},
//		},
//	})
//	require.Nil(t, err)
//	require.NotNil(t, response)
//
//	invoiceId := response.Id
//	eventsMap := aggregateStore.GetEventMap()
//	require.Equal(t, 1, len(eventsMap))
//
//	invoiceTempAggregate := invoice.NewInvoiceTempAggregateWithTenantAndID(tenant, response.Id)
//	eventList := eventsMap[invoiceTempAggregate.ID]
//
//	require.Equal(t, 1, len(eventList))
//
//	require.Equal(t, invoice.InvoiceCreateForContractV1, eventList[0].GetEventType())
//	require.Equal(t, string(invoice.InvoiceAggregateType)+"-temp-"+tenant+"-"+invoiceId, eventList[0].GetAggregateID())
//
//	var eventData invoice.InvoiceCreateEvent
//	err = eventList[0].GetJsonData(&eventData)
//	require.Nil(t, err, "Failed to unmarshal event data")
//
//	// Assertions to validate the contract create event data
//	require.Equal(t, tenant, eventData.Tenant)
//	require.Equal(t, "1", eventData.ContractId)
//	require.Equal(t, now, eventData.CreatedAt)
//	require.Equal(t, "app", eventData.SourceFields.AppSource)
//	require.Equal(t, "source", eventData.SourceFields.Source)
//	require.Equal(t, "source", eventData.SourceFields.SourceOfTruth)
//
//	require.Equal(t, true, eventData.DryRun)
//	require.Equal(t, 36, len(eventData.Number))
//	require.Equal(t, now, eventData.Date)
//	require.Equal(t, now, eventData.DueDate)
//	require.Equal(t, 2, len(eventData.DryRunLines))
//
//	require.Equal(t, "1", eventData.DryRunLines[0].InvoiceId)
//	require.Equal(t, "name 1", eventData.DryRunLines[0].Name)
//	require.Equal(t, commonpb.BilledType_MONTHLY_BILLED.String(), eventData.DryRunLines[0].Billed)
//	require.Equal(t, float64(1), eventData.DryRunLines[0].Price)
//	require.Equal(t, int64(2), eventData.DryRunLines[0].Quantity)
//
//	require.Equal(t, "", eventData.DryRunLines[1].InvoiceId)
//	require.Equal(t, "name 2", eventData.DryRunLines[1].Name)
//	require.Equal(t, commonpb.BilledType_NONE_BILLED.String(), eventData.DryRunLines[1].Billed)
//	require.Equal(t, float64(1), eventData.DryRunLines[1].Price)
//	require.Equal(t, int64(2), eventData.DryRunLines[1].Quantity)
//}
