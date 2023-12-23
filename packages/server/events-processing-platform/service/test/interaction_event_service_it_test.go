package servicet

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	iepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/interaction_event"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInteractionEventService_RequestSummary(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	if err != nil {
		t.Fatalf("Failed to connect to emailEvents processing platform: %v", err)
	}
	client := iepb.NewInteractionEventGrpcServiceClient(grpcConnection)
	interactionEventId := uuid.New().String()
	tenant := "ziggy"
	response, err := client.RequestGenerateSummary(ctx, &iepb.RequestGenerateSummaryGrpcRequest{
		Tenant:             tenant,
		InteractionEventId: interactionEventId,
	})
	if err != nil {
		t.Errorf("Failed to request generate summary: %v", err)
	}
	require.NotNil(t, response)
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	aggr := aggregate.NewInteractionEventAggregateWithTenantAndID(tenant, response.Id)
	eventList := eventsMap[aggr.ID]
	require.Equal(t, 1, len(eventList))
	require.Equal(t, event.InteractionEventRequestSummaryV1, eventList[0].GetEventType())
	require.Equal(t, string(aggregate.InteractionEventAggregateType)+"-"+tenant+"-"+interactionEventId, eventList[0].GetAggregateID())
	var eventData event.InteractionEventRequestSummaryEvent
	if err := eventList[0].GetJsonData(&eventData); err != nil {
		t.Errorf("Failed to unmarshal event data: %v", err)
	}
	require.Equal(t, tenant, eventData.Tenant)
	require.NotNil(t, eventData.RequestedAt)
}
