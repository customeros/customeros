package servicet

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestContractService_SoftDeleteContract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// Setup test environment
	tenant := "ziggy"
	contractId := uuid.New().String()

	// Setup aggregate store and create initial contract
	aggregateStore := eventstoret.NewTestAggregateStore()
	contractAggregate := aggregate.NewContractAggregateWithTenantAndID(tenant, contractId)
	aggregateStore.Save(ctx, contractAggregate)

	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)
	contractClient := contractpb.NewContractGrpcServiceClient(grpcConnection)

	// Create update request
	deleteRequest := &contractpb.SoftDeleteContractGrpcRequest{
		Tenant: tenant,
		Id:     contractId,
	}

	// Execute update contract request
	response, err := contractClient.SoftDeleteContract(ctx, deleteRequest)
	require.Nil(t, err, "Failed to delete contract")

	// Assert response
	require.NotNil(t, response)

	// Retrieve and assert events
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	contractEvents := eventsMap[contractAggregate.ID]
	require.Equal(t, 1, len(contractEvents))

	require.Equal(t, event.ContractDeleteV1, contractEvents[0].GetEventType())

	var eventData event.ContractDeleteEvent
	err = contractEvents[0].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	// Assert event data
	require.NotNil(t, eventData.UpdatedAt)
	require.Equal(t, tenant, eventData.Tenant)
}
