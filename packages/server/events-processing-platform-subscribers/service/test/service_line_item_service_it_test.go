package servicet

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	contractaggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	servicelineitempb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/service_line_item"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func TestServiceLineItemService_CreateServiceLineItem(t *testing.T) {
	ctx := context.Background()

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
	timeNow := utils.Now()
	createRequest := &servicelineitempb.CreateServiceLineItemGrpcRequest{
		Tenant:         tenant,
		LoggedInUserId: "User123",
		ContractId:     contractId,
		Billed:         commonpb.BilledType_MONTHLY_BILLED,
		Quantity:       10,
		Price:          -1.0123456789,
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
	require.Equal(t, -1.0123456789, eventData.Price)
	require.Equal(t, "Test service line item", eventData.Name)
	require.Equal(t, contractId, eventData.ContractId)
}

func TestServiceLineItemService_UpdateServiceLineItem(t *testing.T) {
	ctx := context.Background()
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
		Tenant:                  tenant,
		LoggedInUserId:          "User456",
		Id:                      serviceLineItemId,
		Name:                    "Updated Service Line Item",
		Quantity:                10,
		Price:                   150.0004,
		Comments:                "Some comments",
		IsRetroactiveCorrection: true,
		UpdatedAt:               timestamppb.New(updatedAt),
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
	require.Equal(t, model.NoneBilled.String(), eventData.Billed)
	require.Equal(t, int64(10), eventData.Quantity)
	require.Equal(t, 150.0004, eventData.Price)
	require.Equal(t, "Some comments", eventData.Comments)
	require.Equal(t, "Updated Service Line Item", eventData.Name)
	require.Equal(t, tenant, eventData.Tenant)
}

func TestServiceLineItemService_DeleteServiceLineItem(t *testing.T) {
	ctx := context.Background()
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

func TestServiceLineItemService_UpdateServiceLineItemCreateNewVersion(t *testing.T) {
	//test if the update service line item creates a new version of the service line item
	//and if we provide comments in the update request, it should be added to the new version of the service line item
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Setup test environment
	tenant := "ziggy"
	contractId := uuid.New().String()
	serviceLineItemId := "SLI123"

	// Create and save the initial Service Line Item aggregate

	aggregateStore := eventstoret.NewTestAggregateStore()
	serviceLineItemAggregate := aggregate.NewServiceLineItemAggregateWithTenantAndID(tenant, serviceLineItemId)
	aggregateStore.Save(ctx, serviceLineItemAggregate)
	contractAggregate := contractaggregate.NewContractAggregateWithTenantAndID(tenant, contractId)
	aggregateStore.Save(ctx, contractAggregate)

	// Prepare the gRPC connection and client
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err, "Failed to connect to processing platform")
	serviceLineItemClient := servicelineitempb.NewServiceLineItemGrpcServiceClient(grpcConnection)

	// Create the update request
	updatedAt := utils.Now()
	updateRequest := &servicelineitempb.UpdateServiceLineItemGrpcRequest{
		Tenant:                  tenant,
		LoggedInUserId:          "User456",
		Id:                      serviceLineItemId,
		ContractId:              contractId,
		Name:                    "Updated Service Line Item",
		Billed:                  commonpb.BilledType_MONTHLY_BILLED,
		Quantity:                10,
		Price:                   150.0004,
		Comments:                "Some comments",
		IsRetroactiveCorrection: false,
		UpdatedAt:               timestamppb.New(updatedAt),
		ParentId:                "SLI121",
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
	require.Equal(t, 3, len(eventsMap))
	var eventList []eventstore.Event
	//pick the latest event which is creating a new version of the service line item
	for _, value := range eventsMap {
		eventList = append(eventList, value...)
	}
	createEvent := eventList[1]
	var eventData event.ServiceLineItemCreateEvent
	err = createEvent.GetJsonData(&eventData)

	if createEvent.EventType == event.ServiceLineItemCreateV1 {
		require.Nil(t, err, "Failed to unmarshal event data")
		require.Equal(t, int64(10), eventData.Quantity)
		require.Equal(t, 150.0004, eventData.Price)
		require.Equal(t, "Some comments", eventData.Comments)
		require.Equal(t, "Updated Service Line Item", eventData.Name)
		require.Equal(t, tenant, eventData.Tenant)
		require.Equal(t, contractId, eventData.ContractId)
	}
}
