package servicet

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoice"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInvoiceService_NewInvoice(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// setup test environment
	tenant := "ziggy"
	now := utils.Now()

	aggregateStore := eventstoret.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err, "Failed to get grpc connection")
	invoiceServiceClient := invoicepb.NewInvoiceGrpcServiceClient(grpcConnection)

	response, err := invoiceServiceClient.NewInvoice(ctx, &invoicepb.NewInvoiceRequest{
		Tenant:         tenant,
		OrganizationId: "1",
		SourceFields: &commonpb.SourceFields{
			AppSource: "app",
			Source:    "source",
		},
		CreatedAt: utils.ConvertTimeToTimestampPtr(&now),
		Date:      utils.ConvertTimeToTimestampPtr(&now),
	})
	require.Nil(t, err)
	require.NotNil(t, response)

	invoiceId := response.Id
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))

	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenant, response.Id)
	eventList := eventsMap[invoiceAggregate.ID]
	require.Equal(t, 1, len(eventList))
	require.Equal(t, invoice.InvoiceNewV1, eventList[0].GetEventType())
	require.Equal(t, string(invoice.InvoiceAggregateType)+"-"+tenant+"-"+invoiceId, eventList[0].GetAggregateID())

	var eventData invoice.InvoiceNewEvent
	err = eventList[0].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	// Assertions to validate the contract create event data
	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, "1", eventData.OrganizationId)
	require.Equal(t, now, eventData.CreatedAt)
	require.Equal(t, false, eventData.DryRun)
	require.Equal(t, now, eventData.Date)
	require.Equal(t, now, eventData.DueDate)
	require.Equal(t, "app", eventData.SourceFields.AppSource)
	require.Equal(t, "source", eventData.SourceFields.Source)
	require.Equal(t, "source", eventData.SourceFields.SourceOfTruth)
}

func TestInvoiceService_FillInvoice(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// setup test environment
	tenant := "ziggy"
	invoiceId := "invoice-id"
	now := utils.Now()

	aggregateStore := eventstoret.NewTestAggregateStore()
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenant, invoiceId)

	newEvent, _ := invoice.NewInvoiceNewEvent(invoiceAggregate, "1", false, now, now, now, commonmodel.Source{})
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
		Amount:         1,
		Vat:            2,
		Total:          3,
		Lines: []*invoicepb.InvoiceLine{
			{
				Index:    1,
				Name:     "name",
				Price:    2,
				Quantity: 3,
				Amount:   4,
				Vat:      5,
				Total:    6,
			},
		},
		UpdatedAt: utils.ConvertTimeToTimestampPtr(&now),
		SourceFields: &commonpb.SourceFields{
			AppSource: "app",
		},
	})
	require.Nil(t, err)
	require.NotNil(t, response)

	// verify
	require.Equal(t, invoiceId, response.Id)

	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))

	eventList := eventsMap[invoiceAggregate.ID]
	require.Equal(t, 2, len(eventList))

	require.Equal(t, invoice.InvoiceNewV1, eventList[0].GetEventType())
	require.Equal(t, invoice.InvoiceFillV1, eventList[1].GetEventType())
	require.Equal(t, string(invoice.InvoiceAggregateType)+"-"+tenant+"-"+invoiceId, eventList[1].GetAggregateID())

	var eventData invoice.InvoiceFillEvent
	err = eventList[1].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	require.Equal(t, float64(1), eventData.Amount)
	require.Equal(t, float64(2), eventData.VAT)
	require.Equal(t, float64(3), eventData.Total)
	require.Equal(t, 1, len(eventData.Lines))

	require.Equal(t, int64(1), eventData.Lines[0].Index)
	require.Equal(t, "name", eventData.Lines[0].Name)
	require.Equal(t, float64(2), eventData.Lines[0].Price)
	require.Equal(t, int64(3), eventData.Lines[0].Quantity)
	require.Equal(t, float64(4), eventData.Lines[0].Amount)
	require.Equal(t, float64(5), eventData.Lines[0].VAT)
	require.Equal(t, float64(6), eventData.Lines[0].Total)
}

func TestInvoiceService_PdfGeneratedInvoice(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// setup test environment
	tenant := "ziggy"
	invoiceId := "invoice-id"
	now := utils.Now()

	aggregateStore := eventstoret.NewTestAggregateStore()
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenant, invoiceId)

	newEvent, _ := invoice.NewInvoiceNewEvent(invoiceAggregate, "1", false, now, now, now, commonmodel.Source{})
	fillEvent, _ := invoice.NewInvoiceFillEvent(invoiceAggregate, &now, commonmodel.Source{}, &invoicepb.FillInvoiceRequest{
		Amount: 1,
		Vat:    2,
		Total:  3,
		Lines:  []*invoicepb.InvoiceLine{},
	})
	invoiceAggregate.UncommittedEvents = append(invoiceAggregate.UncommittedEvents, newEvent)
	invoiceAggregate.UncommittedEvents = append(invoiceAggregate.UncommittedEvents, fillEvent)
	aggregateStore.Save(ctx, invoiceAggregate)

	// prepare connection to grpc server
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	invoiceClient := invoicepb.NewInvoiceGrpcServiceClient(grpcConnection)

	// Execute the command
	response, err := invoiceClient.PdfGeneratedInvoice(ctx, &invoicepb.PdfGeneratedInvoiceRequest{
		Tenant:           tenant,
		LoggedInUserId:   "user-id",
		InvoiceId:        invoiceId,
		RepositoryFileId: "repository-file-id",
		UpdatedAt:        utils.ConvertTimeToTimestampPtr(&now),
		SourceFields: &commonpb.SourceFields{
			AppSource: "app",
		},
	})
	require.Nil(t, err)
	require.NotNil(t, response)

	// verify
	require.Equal(t, invoiceId, response.Id)

	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))

	eventList := eventsMap[invoiceAggregate.ID]
	require.Equal(t, 3, len(eventList))

	require.Equal(t, invoice.InvoiceNewV1, eventList[0].GetEventType())
	require.Equal(t, invoice.InvoiceFillV1, eventList[1].GetEventType())
	require.Equal(t, invoice.InvoicePdfGeneratedV1, eventList[2].GetEventType())
	require.Equal(t, string(invoice.InvoiceAggregateType)+"-"+tenant+"-"+invoiceId, eventList[1].GetAggregateID())

	var eventData invoice.InvoicePdfGeneratedEvent
	err = eventList[2].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	require.Equal(t, "repository-file-id", eventData.RepositoryFileId)
}

func TestInvoiceService_PayInvoice(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// setup test environment
	tenant := "ziggy"
	invoiceId := "invoice-id"
	now := utils.Now()

	aggregateStore := eventstoret.NewTestAggregateStore()
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenant, invoiceId)

	newEvent, _ := invoice.NewInvoiceNewEvent(invoiceAggregate, "1", false, now, now, now, commonmodel.Source{})
	fillEvent, _ := invoice.NewInvoiceFillEvent(invoiceAggregate, &now, commonmodel.Source{}, &invoicepb.FillInvoiceRequest{
		Amount: 1,
		Vat:    2,
		Total:  3,
		Lines:  []*invoicepb.InvoiceLine{},
	})
	pdfGeneratedEvent, _ := invoice.NewInvoicePdfGeneratedEvent(invoiceAggregate, &now, commonmodel.Source{}, &invoicepb.PdfGeneratedInvoiceRequest{
		InvoiceId:        invoiceId,
		RepositoryFileId: "repository-file-id",
	})
	invoiceAggregate.UncommittedEvents = append(invoiceAggregate.UncommittedEvents, newEvent)
	invoiceAggregate.UncommittedEvents = append(invoiceAggregate.UncommittedEvents, fillEvent)
	invoiceAggregate.UncommittedEvents = append(invoiceAggregate.UncommittedEvents, pdfGeneratedEvent)
	aggregateStore.Save(ctx, invoiceAggregate)

	// prepare connection to grpc server
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	invoiceClient := invoicepb.NewInvoiceGrpcServiceClient(grpcConnection)

	// Execute the command
	response, err := invoiceClient.PayInvoice(ctx, &invoicepb.PayInvoiceRequest{
		Tenant:         tenant,
		LoggedInUserId: "user-id",
		InvoiceId:      invoiceId,
		UpdatedAt:      utils.ConvertTimeToTimestampPtr(&now),
		SourceFields: &commonpb.SourceFields{
			AppSource: "app",
		},
	})
	require.Nil(t, err)
	require.NotNil(t, response)

	// verify
	require.Equal(t, invoiceId, response.Id)

	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))

	eventList := eventsMap[invoiceAggregate.ID]
	require.Equal(t, 4, len(eventList))

	require.Equal(t, invoice.InvoiceNewV1, eventList[0].GetEventType())
	require.Equal(t, invoice.InvoiceFillV1, eventList[1].GetEventType())
	require.Equal(t, invoice.InvoicePdfGeneratedV1, eventList[2].GetEventType())
	require.Equal(t, invoice.InvoicePayV1, eventList[3].GetEventType())
	require.Equal(t, string(invoice.InvoiceAggregateType)+"-"+tenant+"-"+invoiceId, eventList[0].GetAggregateID())

	var eventData invoice.InvoicePayEvent
	err = eventList[3].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	test.AssertRecentTime(t, eventData.UpdatedAt)
}
