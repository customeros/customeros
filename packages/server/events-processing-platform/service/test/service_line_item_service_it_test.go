package servicet

import (
	"context"
	"github.com/google/uuid"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/common"
	service_line_item_pb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/service_line_item"
	contractaggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
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

	aggregateStore := eventstore.NewTestAggregateStore()
	contractAggregate := contractaggregate.NewContractAggregateWithTenantAndID(tenant, contractId)
	aggregateStore.Save(ctx, contractAggregate)

	// Mock GRPC Connection and ServiceLineItem Client
	grpcConnection, _ := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	serviceLineItemClient := service_line_item_pb.NewServiceLineItemGrpcServiceClient(grpcConnection)

	// Prepare request
	timeNow := time.Now()
	createRequest := &service_line_item_pb.CreateServiceLineItemGrpcRequest{
		Tenant:         tenant,
		LoggedInUserId: "User123",
		ContractId:     contractId,
		Billed:         service_line_item_pb.BilledType_MONTHLY_BILLED,
		Licenses:       10,
		Price:          100.50,
		Description:    "Test service line item",
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
	require.Equal(t, int32(10), eventData.Licenses)
	require.Equal(t, float32(100.50), eventData.Price)
	require.Equal(t, "Test service line item", eventData.Description)
	require.Equal(t, contractId, eventData.ContractId)
}
