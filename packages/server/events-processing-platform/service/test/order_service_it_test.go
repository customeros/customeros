package servicet

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/order"
	organizationaggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	orderpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/order"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func TestOrderService_UpsertOrder_MissingOrganization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Setup test environment
	tenant := "tenant1"
	organizationId := uuid.New().String()

	aggregateStore := eventstoret.NewTestAggregateStore()

	// Mock GRPC Connection and Order Client
	grpcConnection, _ := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	orderClient := orderpb.NewOrderGrpcServiceClient(grpcConnection)

	// Prepare request
	timeNow := utils.Now()

	upsertRequest := &orderpb.UpsertOrderGrpcRequest{
		Tenant:         tenant,
		LoggedInUserId: "User123",
		OrganizationId: organizationId,
		SourceFields: &commonpb.SourceFields{
			AppSource: "unit-test",
		},
		CreatedAt:   timestamppb.New(timeNow),
		ConfirmedAt: timestamppb.New(timeNow),
		PaidAt:      timestamppb.New(timeNow),
		FulfilledAt: timestamppb.New(timeNow),
		CanceledAt:  timestamppb.New(timeNow),
	}

	// Execute service method
	response, err := orderClient.UpsertOrder(ctx, upsertRequest)

	// Validate response and error
	require.NotNil(t, err, "Failed to create order")
	require.Nil(t, response)
}

func TestOrderService_UpsertOrder_AllProperties(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Setup test environment
	tenant := "tenant1"
	organizationId := uuid.New().String()

	aggregateStore := eventstoret.NewTestAggregateStore()
	organizationAggregate := organizationaggregate.NewOrganizationAggregateWithTenantAndID(tenant, organizationId)
	aggregateStore.Save(ctx, organizationAggregate)

	// Mock GRPC Connection and Order Client
	grpcConnection, _ := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	orderClient := orderpb.NewOrderGrpcServiceClient(grpcConnection)

	// Prepare request
	timeNow := utils.Now()

	upsertRequest := &orderpb.UpsertOrderGrpcRequest{
		Tenant:         tenant,
		LoggedInUserId: "User123",
		OrganizationId: organizationId,
		SourceFields: &commonpb.SourceFields{
			AppSource: "unit-test",
		},
		CreatedAt:   timestamppb.New(timeNow),
		ConfirmedAt: timestamppb.New(timeNow),
		PaidAt:      timestamppb.New(timeNow),
		FulfilledAt: timestamppb.New(timeNow),
		CanceledAt:  timestamppb.New(timeNow),
	}

	// Execute service method
	response, err := orderClient.UpsertOrder(ctx, upsertRequest)

	// Validate response and error
	require.Nil(t, err)
	require.NotNil(t, response)

	orderId := response.Id
	require.NotEmpty(t, orderId)

	// Validate created order in event store
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 2, len(eventsMap))

	orderAggregate := order.NewOrderAggregateWithTenantAndID(tenant, orderId)
	eventList := eventsMap[orderAggregate.GetID()]
	require.Equal(t, 1, len(eventList))

	// Validate event data
	createEvent := eventList[0]
	require.Equal(t, order.OrderUpsertV1, createEvent.GetEventType())
	require.Equal(t, string(order.OrderAggregateType)+"-"+tenant+"-"+orderId, eventList[0].GetAggregateID())

	var eventData order.OrderUpsertEvent
	err = createEvent.GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")
	require.Equal(t, organizationId, eventData.OrganizationId)
	require.Equal(t, timeNow, eventData.CreatedAt)
	require.Equal(t, "unit-test", eventData.SourceFields.AppSource)
	require.Equal(t, timeNow, *eventData.ConfirmedAt)
	require.Equal(t, timeNow, *eventData.PaidAt)
	require.Equal(t, timeNow, *eventData.FulfilledAt)
	require.Equal(t, timeNow, *eventData.CanceledAt)
}

func TestOrderService_UpsertOrder_PartialProperties(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Setup test environment
	tenant := "tenant1"
	organizationId := uuid.New().String()

	aggregateStore := eventstoret.NewTestAggregateStore()
	organizationAggregate := organizationaggregate.NewOrganizationAggregateWithTenantAndID(tenant, organizationId)
	aggregateStore.Save(ctx, organizationAggregate)

	// Mock GRPC Connection and Order Client
	grpcConnection, _ := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	orderClient := orderpb.NewOrderGrpcServiceClient(grpcConnection)

	// Prepare request
	timeNow := utils.Now()

	upsertRequest := &orderpb.UpsertOrderGrpcRequest{
		Tenant:         tenant,
		LoggedInUserId: "User123",
		OrganizationId: organizationId,
		SourceFields: &commonpb.SourceFields{
			AppSource: "unit-test",
		},
		PaidAt: timestamppb.New(timeNow),
	}

	// Execute service method
	response, err := orderClient.UpsertOrder(ctx, upsertRequest)

	// Validate response and error
	require.Nil(t, err)
	require.NotNil(t, response)

	orderId := response.Id
	require.NotEmpty(t, orderId)

	// Validate created order in event store
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 2, len(eventsMap))

	orderAggregate := order.NewOrderAggregateWithTenantAndID(tenant, orderId)
	eventList := eventsMap[orderAggregate.GetID()]
	require.Equal(t, 1, len(eventList))

	// Validate event data
	createEvent := eventList[0]
	require.Equal(t, order.OrderUpsertV1, createEvent.GetEventType())
	require.Equal(t, string(order.OrderAggregateType)+"-"+tenant+"-"+orderId, eventList[0].GetAggregateID())

	var eventData order.OrderUpsertEvent
	err = createEvent.GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")
	require.Equal(t, organizationId, eventData.OrganizationId)
	require.Equal(t, "unit-test", eventData.SourceFields.AppSource)
	require.Nil(t, eventData.ConfirmedAt)
	require.Equal(t, timeNow, *eventData.PaidAt)
	require.Nil(t, eventData.FulfilledAt)
	require.Nil(t, eventData.CanceledAt)
}
