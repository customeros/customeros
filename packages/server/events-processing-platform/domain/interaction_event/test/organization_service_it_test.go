package test

import (
	"context"
	"github.com/google/uuid"
	interaction_event_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/interaction_event"
	interactionEventAggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/aggregate"
	interactionEventEvents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/grpc"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var testDatabase *neo4jt.TestDatabase
var dialFactory *grpc.TestDialFactoryImpl

func TestMain(m *testing.M) {
	myDatabase, shutdown := neo4jt.SetupTestDatabase()
	testDatabase = &myDatabase

	dialFactory = &grpc.TestDialFactoryImpl{}
	defer shutdown()

	os.Exit(m.Run())
}

func TestInteractionEventService_RequestSummary(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	if err != nil {
		t.Fatalf("Failed to connect to emailEvents processing platform: %v", err)
	}
	client := interaction_event_grpc_service.NewInteractionEventGrpcServiceClient(grpcConnection)
	interactionEventId := uuid.New().String()
	tenant := "ziggy"
	response, err := client.RequestGenerateSummary(ctx, &interaction_event_grpc_service.RequestGenerateSummaryGrpcRequest{
		Tenant:             tenant,
		InteractionEventId: interactionEventId,
	})
	if err != nil {
		t.Errorf("Failed to request generate summary: %v", err)
	}
	require.NotNil(t, response)
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	aggregate := interactionEventAggregate.NewInteractionEventAggregateWithTenantAndID(tenant, response.Id)
	eventList := eventsMap[aggregate.ID]
	require.Equal(t, 1, len(eventList))
	require.Equal(t, interactionEventEvents.InteractionEventRequestSummaryV1, eventList[0].GetEventType())
	require.Equal(t, string(interactionEventAggregate.InteractionEventAggregateType)+"-"+tenant+"-"+interactionEventId, eventList[0].GetAggregateID())
	var eventData interactionEventEvents.InteractionEventRequestSummaryEvent
	if err := eventList[0].GetJsonData(&eventData); err != nil {
		t.Errorf("Failed to unmarshal event data: %v", err)
	}
	require.Equal(t, tenant, eventData.Tenant)
	require.NotNil(t, eventData.RequestedAt)
}

func tearDownTestCase(ctx context.Context, database *neo4jt.TestDatabase) func(tb testing.TB) {
	return func(tb testing.TB) {
		tb.Logf("Teardown test %v, cleaning neo4j DB", tb.Name())
		neo4jt.CleanupAllData(ctx, database.Driver)
	}
}
