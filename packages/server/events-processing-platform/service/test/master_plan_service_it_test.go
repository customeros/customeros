package servicet

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
	eventstoret "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	masterplanpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/master_plan"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMasterPlanService_CreateMasterPlan(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	// setup test environment
	tenant := "ziggy"

	aggregateStore := eventstoret.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err, "Failed to get grpc connection")
	masterPlanClient := masterplanpb.NewMasterPlanGrpcServiceClient(grpcConnection)

	response, err := masterPlanClient.CreateMasterPlan(ctx, &masterplanpb.CreateMasterPlanGrpcRequest{
		Tenant: tenant,
		Name:   "New Master Plan",
		SourceFields: &commonpb.SourceFields{
			AppSource: "app",
			Source:    "source",
		},
	})
	require.Nil(t, err)
	require.NotNil(t, response)

	masterPlanId := response.Id
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	masterPlanAggregate := aggregate.NewMasterPlanAggregateWithTenantAndID(tenant, response.Id)
	eventList := eventsMap[masterPlanAggregate.ID]
	require.Equal(t, 1, len(eventList))
	require.Equal(t, event.MasterPlanCreateV1, eventList[0].GetEventType())
	require.Equal(t, string(aggregate.MasterPlanAggregateType)+"-"+tenant+"-"+masterPlanId, eventList[0].GetAggregateID())

	var eventData event.MasterPlanCreateEvent
	err = eventList[0].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	// Assertions to validate the contract create event data
	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, "New Master Plan", eventData.Name)
	test.AssertRecentTime(t, eventData.CreatedAt)
	require.Equal(t, "app", eventData.SourceFields.AppSource)
	require.Equal(t, "source", eventData.SourceFields.Source)
	require.Equal(t, "source", eventData.SourceFields.SourceOfTruth)
}
