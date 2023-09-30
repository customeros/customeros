package test

import (
	"context"
	"github.com/google/uuid"
	common_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/common"
	email_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/email"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/grpc"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"testing"
	"time"
)

var testDatabase *test.TestDatabase
var dialFactory *grpc.TestDialFactoryImpl

func TestMain(m *testing.M) {
	myDatabase, shutdown := test.SetupTestDatabase()
	testDatabase = &myDatabase

	dialFactory = &grpc.TestDialFactoryImpl{}
	defer shutdown()

	os.Exit(m.Run())
}

func TestEmailService_UpsertEmail(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := eventstore.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	if err != nil {
		t.Fatalf("Failed to connect to events processing platform: %v", err)
	}
	emailClient := email_grpc_service.NewEmailGrpcServiceClient(grpcConnection)
	timeNow := time.Now().UTC()
	emailId, _ := uuid.NewUUID()
	response, err := emailClient.UpsertEmail(ctx, &email_grpc_service.UpsertEmailGrpcRequest{
		Tenant:   "ziggy",
		RawEmail: "test@openline.ai",
		SourceFields: &common_grpc_service.SourceFields{
			AppSource:     "unit-test",
			Source:        "N/A",
			SourceOfTruth: "N/A",
		},
		CreatedAt: timestamppb.New(timeNow),
		UpdatedAt: timestamppb.New(timeNow),
		Id:        emailId.String(),
	})
	require.Nil(t, err)
	require.NotNil(t, response)
	require.Equal(t, emailId.String(), response.Id)
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList := eventsMap[aggregate.NewEmailAggregateWithTenantAndID("ziggy", emailId.String()).ID]
	require.Equal(t, 1, len(eventList))
	require.Equal(t, events.EmailCreateV1, eventList[0].GetEventType())
	var eventData events.EmailCreateEvent
	if err := eventList[0].GetJsonData(&eventData); err != nil {
		t.Errorf("Failed to unmarshal event data: %v", err)
	}
	require.Equal(t, "test@openline.ai", eventData.RawEmail)
	require.Equal(t, "unit-test", eventData.SourceFields.AppSource)
	require.Equal(t, "N/A", eventData.SourceFields.Source)
	require.Equal(t, "N/A", eventData.SourceFields.SourceOfTruth)
	require.Equal(t, timeNow, eventData.CreatedAt)
	require.Equal(t, timeNow, eventData.UpdatedAt)
	require.Equal(t, "ziggy", eventData.Tenant)

}

func tearDownTestCase(ctx context.Context, database *test.TestDatabase) func(tb testing.TB) {
	return func(tb testing.TB) {
		tb.Logf("Teardown test %v, cleaning neo4j DB", tb.Name())
		neo4jt.CleanupAllData(ctx, database.Driver)
	}
}
