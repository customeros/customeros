package test

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	log_entry_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/log_entry"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
	teventstore "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
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

func TestLogEntryService_UpsertLogEntry_CreateLogEntry(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := teventstore.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err, "Failed to connect to processing platform")
	logEntryClient := log_entry_grpc_service.NewLogEntryGrpcServiceClient(grpcConnection)
	timeNow := utils.Now()
	startedAt := timeNow.Add(-1 * time.Hour)
	tenant := "ziggy"
	response, err := logEntryClient.UpsertLogEntry(ctx, &log_entry_grpc_service.UpsertLogEntryGrpcRequest{
		Tenant:               tenant,
		Content:              "This is a log entry",
		ContentType:          "text/plain",
		StartedAt:            timestamppb.New(startedAt),
		CreatedAt:            timestamppb.New(timeNow),
		Source:               "openline",
		SourceOfTruth:        "openline",
		AppSource:            "unit-test",
		AuthorUserId:         utils.StringPtr("123"),
		LoggedOrganizationId: utils.StringPtr("456"),
	})
	require.Nil(t, err, "Failed to create log entry")

	require.NotNil(t, response)
	logEntryId := response.Id
	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	logEntryAggregate := aggregate.NewLogEntryAggregateWithTenantAndID(tenant, response.Id)
	eventList := eventsMap[logEntryAggregate.ID]
	require.Equal(t, 1, len(eventList))

	require.Equal(t, events.LogEntryCreateV1, eventList[0].GetEventType())
	require.Equal(t, string(aggregate.LogEntryAggregateType)+"-"+tenant+"-"+logEntryId, eventList[0].GetAggregateID())

	var eventData events.LogEntryCreateEvent
	err = eventList[0].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, "unit-test", eventData.AppSource)
	require.Equal(t, "openline", eventData.Source)
	require.Equal(t, "openline", eventData.SourceOfTruth)
	require.Equal(t, timeNow, eventData.CreatedAt)
	require.Equal(t, timeNow, eventData.UpdatedAt)
	require.Equal(t, startedAt, eventData.StartedAt)
	require.Equal(t, "This is a log entry", eventData.Content)
	require.Equal(t, "text/plain", eventData.ContentType)
	require.Equal(t, "123", eventData.AuthorUserId)
	require.Equal(t, "456", eventData.LoggedOrganizationId)
}

func TestLogEntryService_UpsertLogEntry_UpdateLogEntry(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	aggregateStore := teventstore.NewTestAggregateStore()
	grpcConnection, err := dialFactory.GetEventsProcessingPlatformConn(testDatabase.Repositories, aggregateStore)
	require.Nil(t, err)

	logEntryClient := log_entry_grpc_service.NewLogEntryGrpcServiceClient(grpcConnection)

	timeNow := utils.Now()
	logEntryId := uuid.New().String()
	startedAt := timeNow.Add(-1 * time.Hour)
	tenant := "ziggy"

	// prepare aggregate
	logEntryAggregate := aggregate.NewLogEntryAggregateWithTenantAndID(tenant, logEntryId)
	event := eventstore.NewBaseEvent(logEntryAggregate, events.LogEntryCreateV1)
	preconfiguredEventData := events.LogEntryCreateEvent{
		StartedAt:     utils.Now(),
		SourceOfTruth: "openline",
	}
	err = event.SetJsonData(&preconfiguredEventData)
	require.Nil(t, err)
	logEntryAggregate.UncommittedEvents = []eventstore.Event{
		event,
	}
	err = aggregateStore.Save(ctx, logEntryAggregate)
	require.Nil(t, err)

	response, err := logEntryClient.UpsertLogEntry(ctx, &log_entry_grpc_service.UpsertLogEntryGrpcRequest{
		Tenant:               tenant,
		Id:                   logEntryId,
		Content:              "This is a log entry",
		ContentType:          "text/plain",
		StartedAt:            timestamppb.New(startedAt),
		UpdatedAt:            timestamppb.New(timeNow),
		AuthorUserId:         utils.StringPtr("123"),
		LoggedOrganizationId: utils.StringPtr("456"),
	})
	require.Nil(t, err, "Failed to create log entry")

	require.NotNil(t, response)
	require.Equal(t, logEntryId, response.Id)

	eventsMap := aggregateStore.GetEventMap()
	require.Equal(t, 1, len(eventsMap))
	eventList := eventsMap[aggregate.NewLogEntryAggregateWithTenantAndID(tenant, response.Id).ID]
	require.Equal(t, 2, len(eventList))

	require.Equal(t, events.LogEntryCreateV1, eventList[0].GetEventType())
	require.Equal(t, events.LogEntryUpdateV1, eventList[1].GetEventType())
	require.Equal(t, string(aggregate.LogEntryAggregateType)+"-"+tenant+"-"+logEntryId, eventList[1].GetAggregateID())

	var eventData events.LogEntryUpdateEvent
	err = eventList[1].GetJsonData(&eventData)
	require.Nil(t, err, "Failed to unmarshal event data")

	require.Equal(t, tenant, eventData.Tenant)
	require.Equal(t, "openline", eventData.SourceOfTruth)
	require.Equal(t, timeNow, eventData.UpdatedAt)
	require.Equal(t, startedAt, eventData.StartedAt)
	require.Equal(t, "This is a log entry", eventData.Content)
	require.Equal(t, "text/plain", eventData.ContentType)
}

func tearDownTestCase(ctx context.Context, database *test.TestDatabase) func(tb testing.TB) {
	return func(tb testing.TB) {
		tb.Logf("Teardown test %v, cleaning neo4j DB", tb.Name())
		neo4jt.CleanupAllData(ctx, database.Driver)
	}
}
