package servicet

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	contractaggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoice"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func TestInvoiceService_NextPreviewInvoiceForContract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Setup test environment
	aprilFirst := utils.FirstTimeOfMonth(2024, 4)
	expectedPeriodStart := utils.FirstTimeOfMonth(2024, 4)
	expectedPeriodEnd := utils.LastTimeOfMonth(2024, 4)
	tenant := "tenant1"

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenant)
	organizationId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenant, neo4jentity.OrganizationEntity{})
	contractId := neo4jtest.CreateContractForOrganization(ctx, testDatabase.Driver, tenant, organizationId, neo4jentity.ContractEntity{
		BillingCycleInMonths: 1,
		NextInvoiceDate:      utils.ToDatePtr(&aprilFirst),
		Currency:             neo4jenum.CurrencyUSD,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{
		"Tenant":       1,
		"Organization": 1,
		"Contract":     1,
	})

	aggregateStore := eventstoret.NewTestAggregateStore()
	contractAggregate := contractaggregate.NewContractAggregateWithTenantAndID(tenant, contractId)
	aggregateStore.Save(ctx, contractAggregate)

	// Mock GRPC Connection and Invoice Client
	grpcConnection, _ := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	invoiceClient := invoicepb.NewInvoiceGrpcServiceClient(grpcConnection)

	// Prepare request
	request := &invoicepb.NextPreviewInvoiceForContractRequest{
		Tenant:     tenant,
		ContractId: contractId,
		AppSource:  "event-processing-platform",
	}

	// Execute service method
	response, err := invoiceClient.NextPreviewInvoiceForContract(ctx, request)

	// Validate response and error
	require.Nil(t, err)
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
	require.Equal(t, true, eventData.DryRun)
	require.Equal(t, true, eventData.Preview)
	require.Equal(t, false, eventData.OffCycle)
	require.Equal(t, "USD", eventData.Currency)
	require.Equal(t, expectedPeriodStart.Format("dd.MM.yyyy"), eventData.PeriodStartDate.Format("dd.MM.yyyy"))
	require.Equal(t, expectedPeriodEnd.Format("dd.MM.yyyy"), eventData.PeriodEndDate.Format("dd.MM.yyyy"))
}

func TestInvoiceService_NewOnCycleInvoiceForContract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

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
			AppSource: "event-processing-platform",
		},
		CreatedAt:          timestamppb.New(timeNow),
		Currency:           "USD",
		InvoicePeriodStart: utils.ConvertTimeToTimestampPtr(&yesterday),
		InvoicePeriodEnd:   utils.ConvertTimeToTimestampPtr(&nextWeek),
		DryRun:             false,
		OffCycle:           true,
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
	require.Equal(t, "event-processing-platform", eventData.SourceFields.AppSource)
	require.Equal(t, false, eventData.DryRun)
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

	newEvent, _ := invoice.NewInvoiceForContractCreateEvent(invoiceAggregate, common.Source{}, "contract-1", "USD", "test note", 1, false, false, false, false, now, now, now)
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
		Note:           "abc",
		Provider: &invoicepb.FillInvoiceProvider{
			LogoRepositoryFileId: "c",
			Name:                 "d",
			AddressLine1:         "e",
			AddressLine2:         "f",
			Zip:                  "g",
			Locality:             "h",
			Country:              "i",
			Region:               "ii",
		},
		Customer: &invoicepb.FillInvoiceCustomer{
			Name:         "j",
			Email:        "k",
			AddressLine1: "l",
			AddressLine2: "m",
			Zip:          "n",
			Locality:     "o",
			Country:      "p",
			Region:       "pp",
		},
		Amount: 1.01,
		Vat:    2.02,
		Total:  3.03,
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
	require.Equal(t, "abc", eventData.Note)
	require.Equal(t, 9, len(eventData.InvoiceNumber))
	require.Equal(t, "c", eventData.Provider.LogoRepositoryFileId)
	require.Equal(t, "d", eventData.Provider.Name)
	require.Equal(t, "e", eventData.Provider.AddressLine1)
	require.Equal(t, "f", eventData.Provider.AddressLine2)
	require.Equal(t, "g", eventData.Provider.Zip)
	require.Equal(t, "h", eventData.Provider.Locality)
	require.Equal(t, "i", eventData.Provider.Country)
	require.Equal(t, "ii", eventData.Provider.Region)
	require.Equal(t, "j", eventData.Customer.Name)
	require.Equal(t, "k", eventData.Customer.Email)
	require.Equal(t, "l", eventData.Customer.AddressLine1)
	require.Equal(t, "m", eventData.Customer.AddressLine2)
	require.Equal(t, "n", eventData.Customer.Zip)
	require.Equal(t, "o", eventData.Customer.Locality)
	require.Equal(t, "p", eventData.Customer.Country)
	require.Equal(t, "pp", eventData.Customer.Region)
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

	newEvent, _ := invoice.NewInvoiceForContractCreateEvent(invoiceAggregate, common.Source{}, "contract-1", "USD", "test note", 1, false, false, false, false, now, now, now)
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

func TestInvoiceService_PayInvoiceNotification(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// setup test environment
	tenant := "ziggy"
	invoiceId := "invoice-id"
	now := utils.Now()

	aggregateStore := eventstoret.NewTestAggregateStore()
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenant, invoiceId)

	newEvent, _ := invoice.NewInvoiceForContractCreateEvent(invoiceAggregate, common.Source{}, "contract-1", "USD", "test note", 1, false, false, false, false, now, now, now)
	invoiceAggregate.UncommittedEvents = append(invoiceAggregate.UncommittedEvents, newEvent)
	aggregateStore.Save(ctx, invoiceAggregate)

	// prepare connection to grpc server
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	invoiceClient := invoicepb.NewInvoiceGrpcServiceClient(grpcConnection)

	// Execute the command
	response, err := invoiceClient.PayInvoiceNotification(ctx, &invoicepb.PayInvoiceNotificationRequest{
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
	require.Equal(t, invoice.InvoicePayNotificationV1, eventList[1].GetEventType())
	require.Equal(t, string(invoice.InvoiceAggregateType)+"-"+tenant+"-"+invoiceId, eventList[1].GetAggregateID())

	var eventData invoice.InvoicePayNotificationEvent
	err = eventList[1].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	require.Equal(t, tenant, eventData.Tenant)
}

func TestInvoiceService_RequestFillInvoice(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// setup test environment
	tenant := "ziggy"
	invoiceId := "invoice-id"
	now := utils.Now()

	aggregateStore := eventstoret.NewTestAggregateStore()
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenant, invoiceId)

	newEvent, _ := invoice.NewInvoiceForContractCreateEvent(invoiceAggregate, common.Source{}, "contract-1", "USD", "test note", 1, false, false, false, false, now, now, now)
	invoiceAggregate.UncommittedEvents = append(invoiceAggregate.UncommittedEvents, newEvent)
	aggregateStore.Save(ctx, invoiceAggregate)

	// prepare connection to grpc server
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	invoiceClient := invoicepb.NewInvoiceGrpcServiceClient(grpcConnection)

	// Execute the command
	response, err := invoiceClient.RequestFillInvoice(ctx, &invoicepb.RequestFillInvoiceRequest{
		Tenant:         tenant,
		LoggedInUserId: "user-id",
		InvoiceId:      invoiceId,
		ContractId:     "contract-1",
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
	require.Equal(t, invoice.InvoiceFillRequestedV1, eventList[1].GetEventType())
	require.Equal(t, string(invoice.InvoiceAggregateType)+"-"+tenant+"-"+invoiceId, eventList[1].GetAggregateID())

	var eventData invoice.InvoiceFillRequestedEvent
	err = eventList[1].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, "contract-1", eventData.ContractId)
}

func TestInvoiceService_VoidInvoice(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// setup test environment
	tenant := "ziggy"
	invoiceId := "invoice-id"
	now := utils.Now()

	aggregateStore := eventstoret.NewTestAggregateStore()
	invoiceAggregate := invoice.NewInvoiceAggregateWithTenantAndID(tenant, invoiceId)

	newEvent, _ := invoice.NewInvoiceForContractCreateEvent(invoiceAggregate, common.Source{}, "contract-1", "USD", "test note", 1, false, false, false, false, now, now, now)
	invoiceAggregate.UncommittedEvents = append(invoiceAggregate.UncommittedEvents, newEvent)
	aggregateStore.Save(ctx, invoiceAggregate)

	// prepare connection to grpc server
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	invoiceClient := invoicepb.NewInvoiceGrpcServiceClient(grpcConnection)

	// Execute the command
	response, err := invoiceClient.VoidInvoice(ctx, &invoicepb.VoidInvoiceRequest{
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
	require.Equal(t, invoice.InvoiceVoidV1, eventList[1].GetEventType())
	require.Equal(t, string(invoice.InvoiceAggregateType)+"-"+tenant+"-"+invoiceId, eventList[1].GetAggregateID())

	var eventData invoice.InvoiceVoidEvent
	err = eventList[1].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	require.Equal(t, tenant, eventData.Tenant)
}
