package servicet

import (
	"context"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	invoicingcycle "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoicing_cycle"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	invoicingcyclepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoicing_cycle"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInvoicingCycleService_CreateInvoicingCycle(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// setup test environment
	tenant := "ziggy"

	aggregateStore := eventstoret.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err, "Failed to get grpc connection")
	invoicingCycleServiceClient := invoicingcyclepb.NewInvoicingCycleServiceClient(grpcConnection)

	response, err := invoicingCycleServiceClient.CreateInvoicingCycleType(ctx, &invoicingcyclepb.CreateInvoicingCycleTypeRequest{
		Tenant: tenant,
		Type:   invoicingcyclepb.InvoicingDateType_ANNIVERSARY,
		SourceFields: &commonpb.SourceFields{
			AppSource: "app",
			Source:    "source",
		},
	})
	require.Nil(t, err)
	require.NotNil(t, response)

	invoicingCycleId := response.Id
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))

	invoicingCycleAggregate := invoicingcycle.NewInvoicingCycleAggregateWithTenantAndID(tenant, response.Id)
	eventList := eventsMap[invoicingCycleAggregate.ID]
	require.Equal(t, 1, len(eventList))
	require.Equal(t, invoicingcycle.InvoicingCycleCreateV1, eventList[0].GetEventType())
	require.Equal(t, string(invoicingcycle.InvoicingCycleAggregateType)+"-"+tenant+"-"+invoicingCycleId, eventList[0].GetAggregateID())

	var eventData invoicingcycle.InvoicingCycleCreateEvent
	err = eventList[0].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	// Assertions to validate the contract create event data
	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, invoicingcycle.InvoicingCycleTypeAnniversary, eventData.Type)
	test.AssertRecentTime(t, eventData.CreatedAt)
	require.Equal(t, "app", eventData.SourceFields.AppSource)
	require.Equal(t, "source", eventData.SourceFields.Source)
	require.Equal(t, "source", eventData.SourceFields.SourceOfTruth)
}

func TestInvoicingCycleService_UpdateInvoicingCycle(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// setup test environment
	tenant := "ziggy"
	invoicingCycleId := "invoicing-cycle-id"

	aggregateStore := eventstoret.NewTestAggregateStore()
	invoicingCycleAggregate := invoicingcycle.NewInvoicingCycleAggregateWithTenantAndID(tenant, invoicingCycleId)
	createEvent, _ := invoicingcycle.NewInvoicingCycleCreateEvent(invoicingCycleAggregate, invoicingcycle.DATE, commonmodel.Source{})
	invoicingCycleAggregate.UncommittedEvents = append(invoicingCycleAggregate.UncommittedEvents, createEvent)
	aggregateStore.Save(ctx, invoicingCycleAggregate)

	// prepare connection to grpc server
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	invoicingCycleClient := invoicingcyclepb.NewInvoicingCycleServiceClient(grpcConnection)

	// Execute the command
	response, err := invoicingCycleClient.UpdateInvoicingCycleType(ctx, &invoicingcyclepb.UpdateInvoicingCycleTypeRequest{
		Tenant:               tenant,
		InvoicingCycleTypeId: invoicingCycleId,
		Type:                 invoicingcyclepb.InvoicingDateType_ANNIVERSARY,
		LoggedInUserId:       "user-id",
		SourceFields: &commonpb.SourceFields{
			AppSource: "app",
		},
	})
	require.Nil(t, err)
	require.NotNil(t, response)

	// verify
	require.Equal(t, invoicingCycleId, response.Id)

	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList := eventsMap[invoicingCycleAggregate.ID]
	require.Equal(t, 2, len(eventList))
	require.Equal(t, invoicingcycle.InvoicingCycleCreateV1, eventList[0].GetEventType())
	require.Equal(t, invoicingcycle.InvoicingCycleUpdateV1, eventList[1].GetEventType())
	require.Equal(t, string(invoicingcycle.InvoicingCycleAggregateType)+"-"+tenant+"-"+invoicingCycleId, eventList[1].GetAggregateID())

	var eventData invoicingcycle.InvoicingCycleUpdateEvent
	err = eventList[1].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	// Assertions to validate the contract create event data
	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, invoicingcycle.InvoicingCycleTypeAnniversary, eventData.Type)
	test.AssertRecentTime(t, eventData.UpdatedAt)
}
