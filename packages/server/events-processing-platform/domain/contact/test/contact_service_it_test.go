package test

import (
	"context"
	contact_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contact"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/grpc"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"testing"
	"time"
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

func TestContactService_CreateContact(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	if err != nil {
		t.Fatalf("Failed to connect to events processing platform: %v", err)
	}
	contactClient := contact_grpc_service.NewContactGrpcServiceClient(grpcConnection)
	timeNow := time.Now().UTC()
	response, err := contactClient.CreateContact(ctx, &contact_grpc_service.CreateContactGrpcRequest{
		Tenant:        "openline",
		FirstName:     "Bob",
		LastName:      "Smith",
		Prefix:        "Mr.",
		Description:   "This is a contact description",
		AppSource:     "unit-test",
		Source:        "N/A",
		SourceOfTruth: "N/A",
		CreatedAt:     timestamppb.New(timeNow),
	})
	if err != nil {
		t.Errorf("Failed to create contact: %v", err)
	}
	require.NotNil(t, response)
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList := eventsMap[aggregate.NewContactAggregateWithTenantAndID("openline", response.Id).ID]
	require.Equal(t, 1, len(eventList))
	require.Equal(t, events.ContactCreatedV1, eventList[0].GetEventType())
	var eventData events.ContactCreatedEvent
	if err := eventList[0].GetJsonData(&eventData); err != nil {
		t.Errorf("Failed to unmarshal event data: %v", err)
	}
	require.Equal(t, "Bob", eventData.FirstName)
	require.Equal(t, "Smith", eventData.LastName)
	require.Equal(t, "Mr.", eventData.Prefix)
	require.Equal(t, "unit-test", eventData.AppSource)
	require.Equal(t, "N/A", eventData.Source)
	require.Equal(t, "N/A", eventData.SourceOfTruth)
	require.Equal(t, timeNow, eventData.CreatedAt)
	require.Equal(t, timeNow, eventData.UpdatedAt)
	require.Equal(t, "openline", eventData.Tenant)

}

func tearDownTestCase(ctx context.Context, database *neo4jt.TestDatabase) func(tb testing.TB) {
	return func(tb testing.TB) {
		tb.Logf("Teardown test %v, cleaning neo4j DB", tb.Name())
		neo4jt.CleanupAllData(ctx, database.Driver)
	}
}
