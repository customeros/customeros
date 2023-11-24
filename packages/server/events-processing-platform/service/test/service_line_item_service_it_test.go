package servicet

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/common"
	servicelineitempb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/service_line_item"
	contractaggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/model"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestServiceLineItemService_CreateServiceLineItem(t *testing.T) {
	ctx := context.TODO()

	// Setup test environment
	tenant := "tenant1"
	contractId := uuid.New().String()

	aggregateStore := eventstoret.NewTestAggregateStore()
	contractAggregate := contractaggregate.NewContractAggregateWithTenantAndID(tenant, contractId)
	aggregateStore.Save(ctx, contractAggregate)

	// Mock GRPC Connection and ServiceLineItem Client
	grpcConnection, _ := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	serviceLineItemClient := servicelineitempb.NewServiceLineItemGrpcServiceClient(grpcConnection)

	// Prepare request
	timeNow := time.Now()
	createRequest := &servicelineitempb.CreateServiceLineItemGrpcRequest{
		Tenant:         tenant,
		LoggedInUserId: "User123",
		ContractId:     contractId,
		Billed:         servicelineitempb.BilledType_MONTHLY_BILLED,
		Quantity:       10,
		Price:          100.50,
		Name:           "Test service line item",
		SourceFields: &commonpb.SourceFields{
			AppSource: "unit-test",
		},
		CreatedAt: timestamppb.New(timeNow),
		UpdatedAt: timestamppb.New(timeNow),
	}

	// Execute service method
	response, err := serviceLineItemClient.CreateServiceLineItem(ctx, createRequest)

	// Validate response and error
	require.Nil(t, err, "Failed to create service line item")
	require.NotNil(t, response)
	serviceLineItemId := response.Id
	require.NotEmpty(t, serviceLineItemId)

	// Validate created service line item in event store
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 2, len(eventsMap))

	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenant, serviceLineItemId)
	eventList := eventsMap[serviceLineItemAggregate.GetID()]
	require.Equal(t, 1, len(eventList))

	// Validate event data
	createEvent := eventList[0]
	require.Equal(t, event.ServiceLineItemCreateV1, createEvent.GetEventType())
	require.Equal(t, string(aggregate.ServiceLineItemAggregateType)+"-"+tenant+"-"+serviceLineItemId, eventList[0].GetAggregateID())

	var eventData event.ServiceLineItemCreateEvent
	err = createEvent.GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")
	require.Equal(t, model.MonthlyBilled.String(), eventData.Billed)
	require.Equal(t, int64(10), eventData.Quantity)
	require.Equal(t, float64(100.50), eventData.Price)
	require.Equal(t, "Test service line item", eventData.Name)
	require.Equal(t, contractId, eventData.ContractId)
}

func TestServiceLineItemService_UpdateServiceLineItem(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Setup test environment
	tenant := "ziggy"
	serviceLineItemId := "SLI123"

	// Create and save the initial Service Line Item aggregate
	aggregateStore := eventstoret.NewTestAggregateStore()
	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenant, serviceLineItemId)
	aggregateStore.Save(ctx, serviceLineItemAggregate)

	// Prepare the gRPC connection and client
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err, "Failed to connect to processing platform")
	serviceLineItemClient := servicelineitempb.NewServiceLineItemGrpcServiceClient(grpcConnection)

	// Create the update request
	updatedAt := utils.Now()
	updateRequest := &servicelineitempb.UpdateServiceLineItemGrpcRequest{
		Tenant:         tenant,
		LoggedInUserId: "User456",
		Id:             serviceLineItemId,
		Name:           "Updated Service Line Item",
		Quantity:       10,
		Price:          150.0,
		Comments:       "Some comments",
		UpdatedAt:      timestamppb.New(updatedAt),
		SourceFields: &commonpb.SourceFields{
			Source:    "openline",
			AppSource: "unit-test",
		},
	}

	// Call the update service
	response, err := serviceLineItemClient.UpdateServiceLineItem(ctx, updateRequest)
	require.Nil(t, err, "Failed to update service line item")

	// Assertions
	require.NotNil(t, response)
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList := eventsMap[serviceLineItemAggregate.ID]
	require.Equal(t, 1, len(eventList))

	createEvent := eventList[0]
	require.Equal(t, event.ServiceLineItemUpdateV1, createEvent.GetEventType())
	require.Equal(t, string(aggregate.ServiceLineItemAggregateType)+"-"+tenant+"-"+serviceLineItemId, eventList[0].GetAggregateID())

	var eventData event.ServiceLineItemUpdateEvent
	err = createEvent.GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")
	require.Equal(t, model.MonthlyBilled.String(), eventData.Billed)
	require.Equal(t, int64(10), eventData.Quantity)
	require.Equal(t, float64(150.0), eventData.Price)
	require.Equal(t, "Some comments", eventData.Comments)
	require.Equal(t, "Updated Service Line Item", eventData.Name)
	require.Equal(t, tenant, eventData.Tenant)
}

func TestServiceLineItemService_DeleteServiceLineItem(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Setup test environment
	tenant := "ziggy"
	serviceLineItemId := "SLI123"

	// Create and save the initial Service Line Item aggregate
	aggregateStore := eventstoret.NewTestAggregateStore()
	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenant, serviceLineItemId)
	aggregateStore.Save(ctx, serviceLineItemAggregate)

	// Prepare the gRPC connection and client
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	serviceLineItemClient := servicelineitempb.NewServiceLineItemGrpcServiceClient(grpcConnection)

	// Create the request
	request := &servicelineitempb.DeleteServiceLineItemGrpcRequest{
		Tenant:         tenant,
		LoggedInUserId: "User456",
		Id:             serviceLineItemId,
		AppSource:      "unit-test",
	}
	// Execute
	response, err := serviceLineItemClient.DeleteServiceLineItem(ctx, request)
	require.Nil(t, err)

	// Assertions
	require.NotNil(t, response)
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList := eventsMap[serviceLineItemAggregate.ID]
	require.Equal(t, 1, len(eventList))

	createEvent := eventList[0]
	require.Equal(t, event.ServiceLineItemDeleteV1, createEvent.GetEventType())
	require.Equal(t, string(aggregate.ServiceLineItemAggregateType)+"-"+tenant+"-"+serviceLineItemId, eventList[0].GetAggregateID())

	var eventData event.ServiceLineItemDeleteEvent
	err = createEvent.GetJsonData(&eventData)
	require.Nil(t, err)
	require.Equal(t, tenant, eventData.Tenant)
}
